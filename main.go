package main

import (
  "fmt"
  "github.com/99designs/gqlgen/handler"
  "github.com/cryptopickle/invoicespace/auth"
  "github.com/cryptopickle/invoicespace/db/cache"
  "github.com/cryptopickle/invoicespace/db/psql"
  "github.com/cryptopickle/invoicespace/graphqlServer"
  "github.com/cryptopickle/invoicespace/graphqlServer/resolver"
  "github.com/gorilla/mux"
  "github.com/rs/cors"
  "log"
  "net/http"
)

func main(){
	var lh http.Handler
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

	psql, err := psql.NewPostgress(pcs.ConnString());

	if err != nil {
		panic(err)
	}

	 defer psql.PgClient.DB.Close()

	rc, err := cache.NewClient();

	if err != nil {
	  panic(err)
  }

	defer rc.Client.Close()


	r := mux.NewRouter()
	r.Handle("/", handler.Playground("Graphql Playground", "/graphql"))

	rslv  := resolver.Resolver{psql, rc}
	cfg := graphqlServer.Config{Resolvers: &rslv}

	schemas := graphqlServer.NewExecutableSchema(cfg)

	r.Handle("/graphql", handler.GraphQL(schemas))

	lh = auth.HTTPMiddleware(r)

	cos := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "HEAD", "OPTIONS", "PUT", "DELETE"},
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
	})

	lh = cos.Handler(lh);

	http.Handle("/", lh)
	log.Println("Server Started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}