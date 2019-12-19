package cache

import (
  "github.com/go-redis/redis/v7"
  "time"
)


type Client struct {
  *redis.Client
}


func NewClient(options *redis.Options) ( *Client, error ){
  client := redis.NewClient(options)
  _, err := client.Ping().Result();

  if err != nil {
    return nil, err
  }

  return &Client{client}, nil
}



func(c *Client) AddToken(userId, token string) *bool {
  set, err := c.Client.SetNX(token, userId, time.Hour * 1).Result()

  if err != nil {
    panic(err)
  }

  return &set
}

func(c *Client)IsTokenMatches(userId, token string)(*bool){
  ui, err := c.Client.Get(token).Result()
  var isValid  = false

  if err != nil {
    return &isValid
  }

  if userId == ui {
    isValid = true
  }
  return &isValid
}