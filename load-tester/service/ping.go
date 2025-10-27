package service

import (
	"context"

	"load-tester/util/errs"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PingParam struct {
	Message string `json:"message" validate:"required"`
}

type PingResult struct {
	PingID string `json:"ping_id"`
}

func (service *Service) Ping(ctx context.Context, param *PingParam) (*PingResult, error) {
	const op errs.Op = "service/Ping"

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": param,
	}).Info()

	// initialize empty result
	serviceResult := &PingResult{}

	// tidy up result
	serviceResult.PingID = uuid.New().String()

	return serviceResult, nil
}
