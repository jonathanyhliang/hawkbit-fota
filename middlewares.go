package main

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) GetController(ctx context.Context, bid string) (c Controller, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetController", "bid", bid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetController(ctx, bid)
}

func (mw loggingMiddleware) PostCancelActionFeedback(ctx context.Context, bid string, acid string,
	fb CancelActionFeedback) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostCancelActionFeedback", "bid", bid, "acid", acid, "took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.PostCancelActionFeedback(ctx, bid, acid, fb)
}

func (mw loggingMiddleware) PostConfigData(ctx context.Context, bid string, cfg ConfigData) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostConfigData", "bid", bid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostConfigData(ctx, bid, cfg)
}

func (mw loggingMiddleware) GetDeplymentBase(ctx context.Context, bid string) (d DeploymentBase, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetDeplymentBase", "bid", bid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetDeplymentBase(ctx, bid)
}

func (mw loggingMiddleware) PostDeploymentBaseFeedback(ctx context.Context, bid string, acid string,
	fb DeploymentBaseFeedback) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostDeploymentBaseFeedback", "bid", bid, "acid", acid, "took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.PostDeploymentBaseFeedback(ctx, bid, acid, fb)
}

func (mw loggingMiddleware) GetDownloadHttp(ctx context.Context, bid string, ver string) (f []byte, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetDownloadHttp", "bid", bid, "ver", ver, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetDownloadHttp(ctx, bid, ver)
}
