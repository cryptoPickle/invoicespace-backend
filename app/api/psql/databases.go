package psql

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
	"log"
)


type Client struct {
	PgClient *sql.DB
}

func NewPostgress(credentials string) (*Client, error) {
	pg, err := sql.Open("postgres", credentials)

	if err != nil {
		return nil, err
	}

	err = pg.Ping()

	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(pg, &postgres.Config{});
	migration(credentials, driver);

	return &Client{PgClient: pg }, nil
}


type PostgressConnectionString struct {
	Host string
	Port int
	User string
	DatabaseName string
	SSLMode string
	Password string
}

func (pcs *PostgressConnectionString) ConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", pcs.Host, pcs.Port, pcs.User, pcs.Password, pcs.DatabaseName, pcs.SSLMode)
};

func migration(dbName string, driver database.Driver) {
	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		log.Fatal(err)
	}

	log.Println("Database Migrated")
}