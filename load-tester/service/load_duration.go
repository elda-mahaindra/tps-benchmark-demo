package service

import (
	"context"

	"load-tester/util/errs"

	"github.com/sirupsen/logrus"
)

type LoadDurationParam struct {
	BaseParam
	Duration string `json:"duration"`
}

type LoadDurationResult = TestResult

func (service *Service) LoadDuration(ctx context.Context, param *LoadDurationParam) (map[string]any, error) {
	const op errs.Op = "service/LoadDuration"

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": param,
	}).Info()

	// TODO: Implement load duration test

	return nil, nil
}
