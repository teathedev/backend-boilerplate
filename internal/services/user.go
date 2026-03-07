package services

import "github.com/teathedev/backend-boilerplate/pkg/logger"

type userService struct {
	log logger.Logger
}

var UserService userService

func init() {
	UserService = userService{
		log: logger.New("UserService"),
	}
}
