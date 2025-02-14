package http

import "github.com/sirupsen/logrus"

type SwipeAPIController struct {
	log *logrus.Logger
}

func NewSwipeAPIController(log *logrus.Logger) *SwipeAPIController {
	return &SwipeAPIController{
		log: log,
	}
}
