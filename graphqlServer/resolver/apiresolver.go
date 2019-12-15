package resolver

import (
	"context"
	"errors"
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
	psql *psql.Client
	redis *cache.Client
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
	r.psql.CreateUser(u);

	return u, nil
}
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*models.Token, error) {
	user, err := r.psql.GetUserByEmail(email);

	if err != nil {
		return nil, err
	}

	if ok := auth.ComparePassword(password, user.Password); !ok {
		return nil, errors.New("Incorrect Credentials")
	}

	refreshT := auth.JwtCrate(user.ID, time.Now().Add(time.Hour * 8760).Unix())

	r.psql.SaveRefreshToken(refreshT, user.ID)

	return &models.Token{
		AccessToken:     auth.JwtCrate(user.ID, time.Now().Add(time.Hour *  1).Unix()),
		RefreshToken: refreshT,
		ExpiredAt: int(time.Now().Add(time.Hour * 1).Unix()),
	}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context) ([]*models.User, error) {
	panic("not implemented")
}
