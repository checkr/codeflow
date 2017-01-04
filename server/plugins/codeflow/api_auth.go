package codeflow

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/ant0ine/go-json-rest/rest"
	jwt_m "github.com/cheungpat/go-json-rest-middleware-jwt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lestrrat/go-jwx/jwk"
	"github.com/maxwellhealth/bongo"
	"github.com/spf13/viper"
)

type Auth struct {
	Path string
}

func (a *Auth) Register(api *rest.Api) []*rest.Route {
	jwt_middleware := &jwt_m.JWTMiddleware{
		Key:        []byte("secret key"),
		Realm:      "jwt auth",
		Timeout:    time.Hour * 24,
		MaxRefresh: time.Hour * 48,
		Authenticator: func(userId string, password string) bool {
			// Disabled because we dont support username/password login
			return false
		},
	}

	okta := &Okta{
		JWTMiddleware: jwt_middleware,
	}

	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			switch request.URL.Path {
			case a.Path + "/callback/okta":
				return false
			case "/ws":
				return false
			default:
				return true
			}
		},
		IfTrue: jwt_middleware,
	})

	var routes []*rest.Route
	routes = append(routes,
		rest.Post(a.Path+"/callback/okta", okta.oktaCallbackEventHandler),
		rest.Get(a.Path+"/refresh_token", jwt_middleware.RefreshHandler),
	)
	log.Printf("Started the codeflow auth handler on %s\n", a.Path)
	return routes
}

func (a *Auth) handle_auth(w rest.ResponseWriter, r *rest.Request) {

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

// oktaCallbackEventHandler can be used by clients to get a jwt token.
func (o *Okta) oktaCallbackEventHandler(writer rest.ResponseWriter, request *rest.Request) {
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

		user := User{
			Name:     name,
			Username: email,
			Email:    email,
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
