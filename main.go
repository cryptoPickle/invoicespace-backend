package main

import (
	"fmt"
	"github.com/99designs/gqlgen/handler"
	"github.com/cryptopickle/invoicespace/app/api/organisations"
	"github.com/cryptopickle/invoicespace/app/api/refreshToken"
	"github.com/cryptopickle/invoicespace/app/api/roles"
	"github.com/cryptopickle/invoicespace/app/api/user_pool"
	"github.com/cryptopickle/invoicespace/app/api/users"
	"github.com/cryptopickle/invoicespace/auth"
	"github.com/cryptopickle/invoicespace/db/cache"
	"github.com/cryptopickle/invoicespace/db/psql"
	"github.com/cryptopickle/invoicespace/graphqlServer"
	"github.com/cryptopickle/invoicespace/graphqlServer/resolver"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main(){
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from main", r)
		}
	}()
	pcs := psql.PostgressConnectionString{
		Host: "localhost",
		Port: 5432,
		SSLMode: "disable",
		DatabaseName: "users",
		Password: "test",
		User: "postgres",
	}

	newPostgress, err := psql.NewPostgress(pcs.ConnString());

	if err != nil {
		panic(err)
	}

	 defer newPostgress.PgClient.DB.Close()

	options := &redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	}

	rc, err := cache.NewClient(options);

	if err != nil {
	  panic(err)
  }

	defer rc.Client.Close()


	r := mux.NewRouter()
	mw := auth.NewAuthMiddleware(rc)
	r.Use(mw.HTTPMiddleware)
	r.Handle("/", handler.Playground("Graphql Playground", "/graphql"))

	rslv  := newResolvers(newPostgress, rc);


	cfg := graphqlServer.Config{Resolvers: &rslv}

	cfg.Directives.Authorize = auth.Authorise
	cfg.Directives.Role = roles.CheckRole

	schemas := graphqlServer.NewExecutableSchema(cfg)

	r.Handle("/graphql", handler.GraphQL(schemas))

	//cos := cors.New(cors.Options{
	//	AllowedMethods: []string{"GET", "POST", "HEAD", "OPTIONS", "PUT", "DELETE"},
	//	AllowedOrigins: []string{"*"},
	//	AllowedHeaders: []string{"*"},
	//})
	//r.Use(cos.Handler)

	http.Handle("/", r)
	log.Println("Server Started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func newResolvers(psql *psql.Client, cache *cache.Client) resolver.Resolver {
	u := users.Client{psql}
	up := user_pool.Client{psql}
	o := organisations.Client{O: psql, U: &up }
	t := refreshToken.Client{psql}

	return resolver.Resolver{
		Users:        &u,
		RefreshToken: &t,
		Redis:        cache,
		Organisations: &o,
		UserPool: &up,
	}

}