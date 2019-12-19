package auth

import (
	"context"
	"github.com/cryptopickle/invoicespace/app/api"
	"github.com/cryptopickle/invoicespace/db/cache"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net"
	"net/http"
	"strings"
)


var userCtxKey = &ContextKey{"user"}

type ContextKey struct {
	name string
}

type UserClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

type AuthMiddleware struct {
	cache *cache.Client
}

func NewAuthMiddleware(cache *cache.Client) *AuthMiddleware {
	return &AuthMiddleware{cache}
}

func(m *AuthMiddleware) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h := r.Header.Get("Authorization");
		if h == "" {
			next.ServeHTTP(w,r);
			return
		}
		token := GetTokenFromHeader(h)
		userId := UserIdFromToken(token)
		isMatch := m.cache.IsTokenMatches(userId, token)


		if !*isMatch {
			http.Error(w, "Invalid Token", http.StatusForbidden)
			return
		}


		log.Println(h)
		ip, _,_ := net.SplitHostPort(r.RemoteAddr)

		userAuth := &api.Params{
			User: &api.User{
				ID:             userId,
			},
			Request: &api.Request{
				IPAdress: ip,
				Token:    token,
			},
		}
		ctx := context.WithValue(r.Context(), userCtxKey, userAuth)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}


func ForContext(ctx context.Context) *api.Params {
	raw, _ := ctx.Value(userCtxKey).(*api.Params)
	return raw
}

func GetTokenFromHeader(header string) string {
	token := ""
	authFields := strings.Fields(header);
	if len(authFields) > 1 && authFields[0] == "Bearer" {
		token = authFields[1]
	}
	token = header
	return token
}

func UserIdFromToken(token string) string {
	returnToken, err := JWTDecode(token)

	if err != nil {
		return ""
	}

	if claims, ok := returnToken.Claims.(*UserClaims); ok && returnToken.Valid {
		if claims == nil {
			return ""
		}
		return claims.UserId
	} else {
		return ""
	}
}
func HashPassword(p string)([]byte, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	if err != nil {
		return nil, err
	}

	return pw, nil
}

func ComparePassword(p, hash string) bool {
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil {
		return true
	}
	return false
}