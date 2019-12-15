package cache

import (
  "github.com/go-redis/redis/v7"
  "time"
)

type Client struct {
  *redis.Client
}


func NewClient() ( *Client, error ){
  client := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
    Password: "",
    DB: 0,
  })
  _, err := client.Ping().Result();

  if err != nil {
    return nil, err
  }

  return &Client{client}, nil
}

func(c *Client) AddToken(userId, token string) (*bool, error) {
  set, err := c.Client.SetNX(userId, token, time.Hour * 1).Result()

  if err != nil {
    return nil, err
  }

  return &set, nil
}