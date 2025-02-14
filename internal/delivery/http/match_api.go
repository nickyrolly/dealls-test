package http

import "github.com/sirupsen/logrus"

type MatchAPICongtroller struct {
	log *logrus.Logger
}

func NewMatchAPICongtroller(log *logrus.Logger) *MatchAPICongtroller {
	return &MatchAPICongtroller{
		log: log,
	}
}
