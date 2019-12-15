package api

type RefreshToken struct {
  TokenId string `json:"token_id"`
  RefreshToken string `json:"refreshToken"`
  UserId string `json:"user_id"`
}