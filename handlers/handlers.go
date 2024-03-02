package handlers

import "github.com/sirupsen/logrus"

type Handlers struct {
	logger *logrus.Logger
}

func NewHandlers() *Handlers {
	return &Handlers{
		logger: logrus.New(),
	}
}
