package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/cryptopickle/invoicespace/app/api"
	"github.com/cryptopickle/invoicespace/graphqlServer/models"
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


func (c *Client) CreateUser(user  *models.User) models.User {
	query := `INSERT INTO 
		users ("first_name", "last_name", "email", "password", "organisation_id", "disabled", "role") 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING user_id, first_name, last_name, organisation_id, role `

	stmt, err := c.PgClient.Prepare(query);

	if err != nil {
		panic(err)
	}

	defer stmt.Close();

	var u models.User
	if err := stmt.QueryRow(
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.OrganisationID,
		user.Disabled,
		user.Role,
		).Scan(&u.FirstName, &u.LastName, &u.ID, &u.OrganisationID, &u.Role); err != nil {
		log.Println(err)
	}
	return u
}

func (c *Client) GetUserByEmail(email string) (*models.User, error) {
	var err error;
	query := `SELECT  user_id, first_name, last_name, email, password, organisation_id, disabled, role FROM users WHERE email=$1`

	stmt, err := c.PgClient.Prepare(query)

	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	var u models.User

	switch err := stmt.QueryRow(email).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.OrganisationID, &u.Disabled, &u.Role); err {
		case sql.ErrNoRows:
			return nil, errors.New("invalid credentials")
		case nil:
			return &u, nil
		default:
			panic(err)
	}
}


func (c *Client) SaveRefreshToken(refreshToken, userId string) *api.RefreshToken {
	query := `INSERT INTO token_revoke_list (refreshToken, user_id) VALUES ($1, $2) RETURNING token_id`

	stmt, err := c.PgClient.Prepare(query)

	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	var tkn api.RefreshToken

	if err := stmt.QueryRow(refreshToken, userId).Scan(&tkn.TokenId); err != nil{
		panic(err)
	}

	return &tkn
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

