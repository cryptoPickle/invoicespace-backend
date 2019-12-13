package auth

import "github.com/dgrijalva/jwt-go"

var signKey = []byte("dummy")

func JWTDecode(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token)(interface{}, error){
		return signKey, nil
	})
}

func JwtCrate(userId string, expiredAt int64) string {
	claims := UserClaims{
		UserId:         userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt,
			Issuer: "dummy",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, _ := token.SignedString(signKey)

	return ss
}