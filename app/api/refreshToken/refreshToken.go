package refreshToken

import (
  "database/sql"
  "errors"
  "github.com/invoice-space/is-backend/app/api"
  "github.com/invoice-space/is-backend/db/psql"
)

type Client struct {
  *psql.Client
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



func (c *Client) RevokeRefreshToken(userid string) (*string, error) {
  query := `DELETE FROM token_revoke_list WHERE user_id=$1 RETURNING token_id`

  stmt, err := c.PgClient.Prepare(query)

  if err != nil {
    panic(err)
  }

  defer stmt.Close();

  var TokenId string

  switch err := stmt.QueryRow(userid).Scan(&TokenId); err {
  case sql.ErrNoRows:
    return nil, errors.New("invalid credentials")
  case nil:
    return &TokenId, nil
  default:
    panic(err)
  }
}
