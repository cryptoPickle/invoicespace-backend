package roles

import (
  "context"
  "errors"
  "github.com/99designs/gqlgen/graphql"
  "github.com/cryptopickle/invoicespace/auth"
  "github.com/cryptopickle/invoicespace/graphqlServer/models"
  "strings"
)

type Role string

var RoleWrapper_Role_Value = map[string]int32{
  "baseuser" : 1,
  "organisationworker": 2,
  "organisationadmin": 3,
  "admin" : 4,
}


var RoleWrapper_Role_Name = map[int32]string{
  1 : "baseuser",
  2 : "organisationworker",
  3: "organisationadmin",
  4 : "admin",
}

func CheckRole(ctx context.Context, obj interface{}, next graphql.Resolver, role models.Role) (interface{}, error) {
  isValid := false

  apiParams, err := auth.GetApiParams(ctx)
  if err != nil {
    return nil, err
  }

  if apiParams.User == nil {
    return nil, errors.New("No user available")
  }



  if apiParams.User.Role == 0 {
    return nil, errors.New("No Role Has Been Set")
  }


  if apiParams.User.Role >= int(RoleWrapper_Role_Value[strings.ToLower(string(role))]) {
    isValid = true
  }


  if !isValid {
    return nil, errors.New("Action not allowed by role")
  }

  return next(ctx)
}