package dto

import "github.com/pratyush934/sibling-bond-server/models"

type LoginResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}
