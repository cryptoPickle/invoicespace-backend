package psql

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)


type Client struct {
	PgClient *sqlx.DB
}

func NewPostgress(credentials string) (*Client, error) {
	pg, err := sqlx.Connect("postgres", credentials)

	if err != nil {
		return nil, err
	}
	err = pg.Ping()

	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(pg.DB, &postgres.Config{});
	err = migration(driver)

	if err != nil {
		panic(err)
	}

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

func migration(driver database.Driver) error {

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations/", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange{
			log.Println("Database Migration No Change!")
			return nil
		}
		return err
	}

	log.Println("Database Migrated")
	return nil
}

