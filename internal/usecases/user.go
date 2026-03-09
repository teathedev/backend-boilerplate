package usecases

import "github.com/teathedev/pkg/logger"

type userUseCase struct {
	log logger.Logger
}

var User userUseCase

func init() {
	User = userUseCase{
		log: logger.New("User"),
	}
}
