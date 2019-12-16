package resolver

import (
	"context"
	"errors"
  "fmt"
  "github.com/cryptopickle/invoicespace/auth"
  "github.com/cryptopickle/invoicespace/db/cache"
  "github.com/cryptopickle/invoicespace/db/psql"
	"log"
	"time"

	"github.com/cryptopickle/invoicespace/graphqlServer"
	"github.com/cryptopickle/invoicespace/graphqlServer/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{
	Psql  *psql.Client
	Redis *cache.Client
}

func (r *Resolver) Mutation() graphqlServer.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graphqlServer.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, input models.NewUser) (*models.User, error) {
	pass, err := auth.HashPassword(input.Password)
	if err != nil {
		log.Fatal("Cannot create user",err)
	}
	u := &models.User{
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Email:          input.Email,
		Password:       string(pass),
		OrganisationID: input.OrganisationID,
	}
	log.Println(u)
	r.Psql.CreateUser(u);

	return u, nil
}
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*models.Token, error) {
	user, err := r.Psql.GetUserByEmail(email);

	if err != nil {
		return nil, err
	}

	if ok := auth.ComparePassword(password, user.Password); !ok {
		return nil, errors.New("Incorrect Credentials")
	}

	refreshT := auth.JwtCrate(user.ID, time.Now().Add(time.Hour * 8760).Unix())

	r.Psql.SaveRefreshToken(refreshT, user.ID)

	accessT := auth.JwtCrate(user.ID, time.Now().Add(time.Hour *  1).Unix())

	ok := r.Redis.AddToken(user.ID, accessT)

  fmt.Println(*ok)
	return &models.Token{
		AccessToken:    accessT ,
		RefreshToken: refreshT,
		ExpiredAt: int(time.Now().Add(time.Hour * 1).Unix()),
	}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context) ([]*models.User, error) {
	panic("not implemented")
}
