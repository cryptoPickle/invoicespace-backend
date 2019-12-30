package users

import (
  "github.com/cryptopickle/invoicespace/db/psql"
  "github.com/cryptopickle/invoicespace/graphqlServer/models"
  "log"
)

type Client struct {
  *psql.Client
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

  err = stmt.QueryRow(email).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.OrganisationID, &u.Disabled, &u.Role);
  checkErr := psql.HandleError(err)
  if checkErr != nil {
    return nil, checkErr
  }
  return &u, nil
}

func (c *Client) UpdateUserRole(userId string, userRole int) (*models.User, error) {
  query := `UPDATE users SET role = $2 WHERE user_id = $1 RETURNING user_id, role`

  stmt, err := c.PgClient.Prepare(query)

  if err != nil {
    return nil, err
  }

  var user models.User
  err = stmt.QueryRow(userId, userRole).Scan(&user.ID, &user.Role)
  checkErr := psql.HandleError(err)

  if checkErr != nil {
    return nil, checkErr
  }

  return &user, nil
}
