package http

import "github.com/sirupsen/logrus"

type LikeAPIController struct {
	log *logrus.Logger
}

func NewLikeAPIController(log *logrus.Logger) *LikeAPIController {
	return &LikeAPIController{
		log: log,
	}
}
