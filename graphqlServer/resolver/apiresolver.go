package resolver

import (
	"context"
	"fmt"
	"github.com/cryptopickle/invoicespace/app/api/psql"
	"github.com/cryptopickle/invoicespace/auth"
	"log"
	"time"

	"github.com/cryptopickle/invoicespace/graphqlServer"
	"github.com/cryptopickle/invoicespace/graphqlServer/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{
	*psql.Client
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
	r.Client.CreateUser(u);

	return u, nil
}
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*models.Token, error) {
	user, err := r.Client.GetUserByEmail(email);

	if err != nil {
		panic(err)
	}

	fmt.Println(user)
	return &models.Token{
		Token:     auth.JwtCrate("123sda", time.Now().Add(time.Hour *  1).Unix()),
		ExpiredAt: int(time.Now().Add(time.Hour * 1).Unix()),
	}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context) ([]*models.User, error) {
	panic("not implemented")
}
