package organisations

import (
  "github.com/invoice-space/is-backend/app/api/user_pool"
  "github.com/invoice-space/is-backend/db/psql"
  "github.com/invoice-space/is-backend/graphqlServer/models"
)

type Client struct {
  O *psql.Client
  U *user_pool.Client
}

//TODO create userpool and assing the user as admin

func (c *Client) CreateOrganisation(org models.Organisation, userId string) (*string, error) {
  query := `INSERT INTO organisations ("name", "description") VALUES ($1, $2) RETURNING id`
  stmt, err := c.O.PgClient.Prepare(query)

  if err != nil {
    return nil, err
  }
  var organisation = models.Organisation{}

  if err := stmt.QueryRow(org.Name, org.Description).Scan(&organisation.ID); err != nil {
    return nil, err
  }

  _, err = c.U.CreateUserPool(organisation.ID, userId)

  if err != nil {
    return nil, err
  }
  return &organisation.ID, nil
}

