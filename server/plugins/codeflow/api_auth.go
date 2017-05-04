package codeflow

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/checkr/codeflow/server/agent"
	jwt_m "github.com/cheungpat/go-json-rest-middleware-jwt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lestrrat/go-jwx/jwk"
	"github.com/maxwellhealth/bongo"
	"github.com/spf13/viper"
)

type Auth struct {
	Path          string
	JWTMiddleware *jwt_m.JWTMiddleware
}

func (a *Auth) Register(api *rest.Api) []*rest.Route {
	if viper.GetString("plugins.codeflow.jwt_secret_key") == "" {
		panic("plugins.codeflow.jwt_secret_key is empty")
	}

	jwt_middleware := &jwt_m.JWTMiddleware{
		Key:        []byte(viper.GetString("plugins.codeflow.jwt_secret_key")),
		Realm:      "jwt auth",
		Timeout:    time.Hour * 24,
		MaxRefresh: time.Hour * 48,
		Authenticator: func(userId string, password string) bool {
			// Disabled because we dont support username/password login
			return false
		},
	}

	a.JWTMiddleware = jwt_middleware

	okta := &Okta{
		JWTMiddleware: jwt_middleware,
	}

	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			var mockEvents = regexp.MustCompile(`/mockEvents/*`)
			switch {
			case (request.URL.Path == (a.Path + "/handler")):
				return false
			case (request.URL.Path == (a.Path + "/callback/okta")):
				return false
			case (request.URL.Path == (a.Path + "/callback/demo")):
				if viper.GetString("environment") == "development" {
					return false
				} else {
					return true
				}
			case request.URL.Path == "/ws":
				return false
			case mockEvents.MatchString(request.URL.Path):
				return false
			default:
				return true
			}
		},
		IfTrue: jwt_middleware,
	})

	var routes []*rest.Route
	routes = append(routes,
		rest.Get(a.Path+"/handler", a.handler),
		rest.Get(a.Path+"/refresh_token", jwt_middleware.RefreshHandler),
		rest.Post(a.Path+"/callback/okta", okta.oktaCallbackHandler),
		rest.Post(a.Path+"/callback/demo", a.demoCallbackHandler),
	)
	log.Printf("Started the codeflow auth handler on %s\n", a.Path)
	return routes
}

func (a *Auth) handler(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(resultHandler{Handler: viper.GetString("plugins.codeflow.auth.handler")})
}

func (a *Auth) demoCallbackHandler(writer rest.ResponseWriter, r *rest.Request) {
	user := User{
		Name:       "Demo",
		Username:   "demo@development.com",
		Email:      "demo@development.com",
		IsAdmin:    true,
		IsEngineer: true,
	}

	usersCol := db.Collection("users")
	if err := usersCol.FindOne(bson.M{"email": user.Email}, &user); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Users::FindOne::DocumentNotFound: email: `%v`", user.Email)
			log.Printf("Creating new user: email: `%v`", user.Email)
			if err := usersCol.Save(&user); err != nil {
				log.Printf("Users::Save::Error: %v", err.Error())
				rest.Error(writer, err.Error(), http.StatusInternalServerError)
			}
		} else {
			log.Printf("Users::FindOne::Error: %v", err.Error())
			rest.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if err := usersCol.Save(&user); err != nil {
			log.Printf("Users::Save::Error: %v", err.Error())
			rest.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	token := jwt.New(jwt.GetSigningMethod(a.JWTMiddleware.SigningAlgorithm))
	claims := make(jwt.MapClaims)

	if a.JWTMiddleware.PayloadFunc != nil {
		for key, value := range a.JWTMiddleware.PayloadFunc(user.Username) {
			claims[key] = value
		}
	}

	claims["sub"] = user.Username
	claims["exp"] = time.Now().Add(a.JWTMiddleware.Timeout).Unix()

	if a.JWTMiddleware.MaxRefresh != 0 {
		claims["orig_iat"] = time.Now().Unix()
	}

	token.Claims = claims
	tokenString, err := token.SignedString(a.JWTMiddleware.Key)

	if err != nil {
		a.unauthorized(writer)
		return
	}

	writer.WriteJson(resultToken{Token: tokenString})
	return
}

type resultHandler struct {
	Handler string `json:"handler"`
}

type Okta struct {
	JWTMiddleware *jwt_m.JWTMiddleware
}

type resultToken struct {
	Token string `json:"token"`
}

type idToken struct {
	Token string `json:"idToken"`
}

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (o *Okta) unauthorized(writer rest.ResponseWriter) {
	writer.Header().Set("WWW-Authenticate", "JWT realm="+o.JWTMiddleware.Realm)
	rest.Error(writer, "Not Authorized", http.StatusUnauthorized)
}

func (a *Auth) unauthorized(writer rest.ResponseWriter) {
	writer.Header().Set("WWW-Authenticate", "JWT realm="+a.JWTMiddleware.Realm)
	rest.Error(writer, "Not Authorized", http.StatusUnauthorized)
}

func (o *Okta) getKey(token *jwt.Token) (interface{}, error) {
	jwksURL := fmt.Sprintf("https://%v.okta.com/oauth2/v1/keys", viper.GetString("plugins.codeflow.auth.okta_org"))
	set, err := jwk.FetchHTTP(jwksURL)
	if err != nil {
		return nil, err
	}

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have string kid")
	}

	if key := set.LookupKeyID(keyID); len(key) == 1 {
		return key[0].Materialize()
	}

	return nil, errors.New("unable to find key")
}

func updateUserPermissions(user *User, groups []interface{}) {
	// Set user permissions
	admin_groups := viper.GetStringSlice("plugins.codeflow.auth.admin_groups")
	engineer_groups := viper.GetStringSlice("plugins.codeflow.auth.engineer_groups")

	for _, g := range groups {
		if user.IsAdmin == false {
			user.IsAdmin = agent.SliceContains(g.(string), admin_groups)
		}
		if user.IsEngineer == false {
			user.IsEngineer = agent.SliceContains(g.(string), engineer_groups)
		}
	}
}

// oktaCallbackEventHandler can be used by clients to get a jwt token.
func (o *Okta) oktaCallbackHandler(writer rest.ResponseWriter, request *rest.Request) {
	var idToken map[string]string
	if err := request.DecodeJsonPayload(&idToken); err != nil {
		rest.Error(writer, fmt.Sprintf("No id token received: %#v", idToken), http.StatusBadRequest)
		return
	}
	token, err := jwt.Parse(idToken["idToken"], o.getKey)

	// JWT Successfull
	if token.Valid {
		oktaClaims := token.Claims.(jwt.MapClaims)
		email := oktaClaims["email"].(string)
		name := oktaClaims["name"].(string)
		groups := oktaClaims["groups"].([]interface{})

		user := User{
			Name:       name,
			Username:   email,
			Email:      email,
			IsAdmin:    false,
			IsEngineer: false,
		}

		usersCol := db.Collection("users")
		if err := usersCol.FindOne(bson.M{"email": user.Email}, &user); err != nil {
			if _, ok := err.(*bongo.DocumentNotFoundError); ok {
				log.Printf("Users::FindOne::DocumentNotFound: email: `%v`", user.Email)
				log.Printf("Creating new user: email: `%v`", user.Email)
				updateUserPermissions(&user, groups)
				if err := usersCol.Save(&user); err != nil {
					log.Printf("Users::Save::Error: %v", err.Error())
					rest.Error(writer, err.Error(), http.StatusInternalServerError)
				}
			} else {
				log.Printf("Users::FindOne::Error: %v", err.Error())
				rest.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			updateUserPermissions(&user, groups)
			if err := usersCol.Save(&user); err != nil {
				log.Printf("Users::Save::Error: %v", err.Error())
				rest.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		token := jwt.New(jwt.GetSigningMethod(o.JWTMiddleware.SigningAlgorithm))
		claims := make(jwt.MapClaims)

		if o.JWTMiddleware.PayloadFunc != nil {
			for key, value := range o.JWTMiddleware.PayloadFunc(user.Username) {
				claims[key] = value
			}
		}

		claims["sub"] = user.Username
		claims["exp"] = time.Now().Add(o.JWTMiddleware.Timeout).Unix()

		if o.JWTMiddleware.MaxRefresh != 0 {
			claims["orig_iat"] = time.Now().Unix()
		}

		token.Claims = claims
		tokenString, err := token.SignedString(o.JWTMiddleware.Key)

		if err != nil {
			o.unauthorized(writer)
			return
		}

		writer.WriteJson(resultToken{Token: tokenString})
		return
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			rest.Error(writer, fmt.Sprintf("ID Token is malformed"), http.StatusBadRequest)
			return
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			rest.Error(writer, fmt.Sprintf("ID Token is expired or not active yet: %s", err), http.StatusBadRequest)
			return
		}
		rest.Error(writer, fmt.Sprintf("Could not handle ID Token: %s", err), http.StatusBadRequest)
		return
	}
}
