package auth_services

import "go-gin-example/models"

type Auth struct {
	Username string
	Password string
}

func (auth *Auth) CheckAuth() (bool, error) {
	return models.CheckAuth(auth.Username, auth.Password)
}
