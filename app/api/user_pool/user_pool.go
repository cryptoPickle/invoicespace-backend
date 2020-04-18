package user_pool

import "github.com/invoice-space/is-backend/db/psql"

type Client struct {
  *psql.Client
}

type UserPool struct {
  UserId string
  User_role string
  UserPoolId string
}

type UserPools struct {
  ID string
  OrganisationId string
}

func (c *Client) CreateUserPool(organisationId, userId string) (*UserPools, error){
  query := `INSERT INTO user_pools (organisation_id) VALUES ($1) RETURNING id`

  stmt, err := c.PgClient.Prepare(query)

  if err != nil {
    return nil, err
  }
  defer stmt.Close()

  var userPool = UserPools{}

  if err := stmt.QueryRow(organisationId).Scan(&userPool.ID); err != nil {
    return nil, err
  }

  _, err = c.AssignUser(userId, userPool.ID, 2)

  if err != nil {
    return nil, err
  }

  return &userPool, err

}

func (c *Client) AssignUser(userid, userPoolId string, role int) (*UserPool, error) {
  query := `INSERT INTO user_pool (userId, user_role, id) VALUES ($1, $2, $3) RETURNING userId, user_role`

  stmt, err := c.PgClient.Prepare(query)

  if err != nil {
    return nil, err
  }
  defer stmt.Close()

  var user = UserPool{}
  if err := stmt.QueryRow(userid, role, userPoolId).Scan(&user.UserId, &user.User_role); err != nil {
    return nil,err
  }
  return &user,nil
}

func (c *Client) IsUserInOrganisation(userId, organisationId string) (*bool, error) {
  query := `SELECT userid FROM  user_pool LEFT JOIN user_pools ON user_pools.id = user_pool.id WHERE organisation_id = $1 AND userid = $2 `

  stmt, err := c.PgClient.Prepare(query)

  if err != nil {
    return nil, err
  }
  var up = UserPool{}
  if err := stmt.QueryRow(organisationId, userId).Scan(&up.UserId); err != nil {
    return nil, err
  }

  isValid := false

  if up.UserId != "" {
    isValid = true
    return &isValid, nil
  }

  return &isValid, nil

}
