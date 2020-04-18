package auth

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/dgrijalva/jwt-go"
	"github.com/invoice-space/is-backend/app/api"
	"github.com/invoice-space/is-backend/db/cache"
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
	User *api.User
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
		user, err := ValidateToken(token)

		if err != nil {
			http.Error(w, "Token error", http.StatusInternalServerError)
			return
		}
		isMatch := m.cache.IsTokenMatches(user.ID, token)


		if !*isMatch {
			http.Error(w, "Invalid Token", http.StatusForbidden)
			return
		}


		log.Println(h)
		ip, _,_ := net.SplitHostPort(r.RemoteAddr)

		userAuth := &api.Params{
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

func authorise(ctx context.Context) (*api.Params, error){
	apiParams, ok := ctx.Value(userCtxKey).(*api.Params)

	if !ok {
		return nil, errors.New("Wrong params")
	}

	if apiParams.Request.Token == "" {
		return nil, errors.New("Unauthorised")
	}

	usr, err := ValidateToken(apiParams.Request.Token);

	if err != nil {
		return nil, errors.New("Token Error")
	}

	if usr == nil {
		return nil, errors.New("No user")
	}

	apiParams.User = usr

	return apiParams, nil

}

func Authorise(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	log.Println("Checking Authorisation")
	params, err := authorise(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ctx = context.WithValue(ctx, userCtxKey, params)
	return next(ctx)
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

func ValidateToken(token string) (*api.User, error) {
	returnToken, err := JWTDecode(token)

	if err != nil {
		return nil, err
	}

	if claims, ok := returnToken.Claims.(*UserClaims); ok && returnToken.Valid {
		if claims == nil {
			return nil, errors.New("Invalid Token")
		}
		return  claims.User, nil
	} else {
		return nil, errors.New("Invalid Token")
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

func GetUserFromContext(ctx context.Context) (*api.User, error) {
	raw, found := ctx.Value(userCtxKey).(*api.Params)
	if !found {
		return nil, errors.New("No params provided")
	}
	return raw.User, nil
}

func GetApiParams(ctx context.Context) (*api.Params, error) {
	raw, found := ctx.Value(userCtxKey).(*api.Params)
	if !found {
		return nil, errors.New("No params provided")
	}
	return raw, nil
}