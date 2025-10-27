package service

import (
	"context"

	"load-tester/util/errs"

	"github.com/sirupsen/logrus"
)

type LoadIncrementalParam struct {
	BaseParam
	TotalReqs   int `json:"total_reqs"`   // Total number of requests to send
	StartingRPS int `json:"starting_rps"` // Starting requests per second
}

type LoadIncrementalResult = TestResult

// Gradually increasing RPS
func (service *Service) LoadIncremental(ctx context.Context, param *LoadIncrementalParam) (map[string]any, error) {
	const op errs.Op = "service/LoadIncremental"

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": param,
	}).Info()

	// TODO: Implement load incremental test

	return nil, nil
}
