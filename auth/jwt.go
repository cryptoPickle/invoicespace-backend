package auth

import (
	"github.com/cryptopickle/invoicespace/app/api"
	"github.com/cryptopickle/invoicespace/graphqlServer/models"
	"github.com/dgrijalva/jwt-go"
)

var signKey = []byte("dummy")

func JWTDecode(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token)(interface{}, error){
		return signKey, nil
	})
}

func JwtCrate(user *models.User, expiredAt int64) string {
	claims := UserClaims{
		User: &api.User{
			ID:             user.ID,
			Role:           *user.Role,

		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt,
			Issuer: "dummy",
		},
	}

	if user.OrganisationID != nil {
		claims.User.OrganisationId = *user.OrganisationID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, _ := token.SignedString(signKey)

	return ss
}