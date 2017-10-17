package resolvers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/codeamp/circuit/plugins/codeamp/models"
	"github.com/codeamp/circuit/plugins/codeamp/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	graphql "github.com/neelance/graphql-go"
	"github.com/spf13/viper"
)

type UserInput struct {
	Email    string
	Password string
}

func (r *Resolver) CreateUser(args *struct{ User *UserInput }) *UserResolver {
	passwordHash, _ := utils.HashPassword(args.User.Password)

	user := models.User{
		Email:    args.User.Email,
		Password: passwordHash,
	}

	r.db.Create(&user)

	return &UserResolver{db: r.db, User: user}
}

func (r *Resolver) UserToken(args *struct {
	Email    string
	Password string
}) (*UserResolver, error) {
	var user models.User

	if r.db.Where("email = ?", args.Email).First(&user).RecordNotFound() {
		return nil, errors.New("Authentication failed")
	}

	if !utils.CheckPasswordHash(args.Password, user.Password) {
		return nil, errors.New("Authentication failed")
	}

	return &UserResolver{db: r.db, User: user}, nil
}

func (r *Resolver) User(ctx context.Context, args *struct{ ID *graphql.ID }) (*UserResolver, error) {
	var err error
	var userId string
	var user models.User

	if userId, err = utils.CheckAuth(ctx, []string{"admin", fmt.Sprintf("user:%s", args.ID)}); err != nil {
		return nil, err
	}

	if args.ID != nil {
		userId = string(*args.ID)
	}

	if err := r.db.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}

	return &UserResolver{db: r.db, User: user}, nil
}

type UserResolver struct {
	db   *gorm.DB
	User models.User
}

func (r *UserResolver) ID() graphql.ID {
	return graphql.ID(r.User.Model.ID.String())
}

func (r *UserResolver) Email(ctx context.Context) (string, error) {
	return r.User.Email, nil
}

func (r *UserResolver) Permissions() []string {
	var permissions []string

	r.db.Model(r.User).Association("Permissions").Find(&r.User.Permissions)

	permissions = append(permissions, fmt.Sprintf("user:%s", r.User.Model.ID))

	for _, permission := range r.User.Permissions {
		permissions = append(permissions, permission.Value)
	}

	return permissions
}

func (r *UserResolver) Token() (string, error) {
	claims := utils.JWTClaims{
		UserId:      r.User.Model.ID.String(),
		Permissions: r.Permissions(),
		StandardClaims: jwt.StandardClaims{
			Issuer:    viper.GetString("plugins.codeamp.jwt_issuer"),
			IssuedAt:  time.Now().UTC().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	key := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := key.SignedString([]byte(viper.GetString("plugins.codeamp.jwt_secret")))
	if err != nil {
		return "", err
	}

	return token, nil
}
