package resolver

import (
	"context"
	"errors"
	"github.com/cryptopickle/invoicespace/app/api/organisations"
	"github.com/cryptopickle/invoicespace/app/api/refreshToken"
	"github.com/cryptopickle/invoicespace/app/api/users"
	"github.com/cryptopickle/invoicespace/auth"
	"github.com/cryptopickle/invoicespace/db/cache"
	"log"
	"time"

	"github.com/cryptopickle/invoicespace/graphqlServer"
	"github.com/cryptopickle/invoicespace/graphqlServer/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{
	Users *users.Client
	RefreshToken *refreshToken.Client
	Organisations *organisations.Client
	Redis *cache.Client
}

func (r *Resolver) Mutation() graphqlServer.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graphqlServer.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateOrganisation(ctx context.Context, input models.NewOrganisation) (*models.Organisation, error) {
	var disabled = false
	user := auth.GetUserFromContext(ctx);
	orgId, err := r.Organisations.CreateOrganisation(models.Organisation{
		Name:        input.Name,
		Description: input.Description,
		WorkerLimit: 1,
		UserLimit:   3,
		Disabled:    &disabled,
	}, user.ID)

	if err != nil {
		return nil, err
	}

	_, err = r.Users.UpdateUserRole(user.ID, 2)

	if err != nil {
		return nil, err
	}
	return  &models.Organisation{
		ID: *orgId,
	}, nil
}

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
		Role: &input.Role,
	}
	log.Println(u)
	r.Users.CreateUser(u);

	return u, nil
}
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*models.Token, error) {
	user, err := r.Users.GetUserByEmail(email);

	if err != nil {
		return nil, err
	}

	if ok := auth.ComparePassword(password, user.Password); !ok {
		return nil, errors.New("Incorrect Credentials")
	}

	refreshT := auth.JwtCrate(user, time.Now().Add(time.Hour * 8760).Unix())

	r.RefreshToken.SaveRefreshToken(refreshT, user.ID)

	accessT := auth.JwtCrate(user, time.Now().Add(time.Hour *  1).Unix())

  r.Redis.AddToken(user.ID, accessT)


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
