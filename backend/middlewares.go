package backend

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a BackendService (as opposed to endpoint) middleware.
type Middleware func(BackendService) BackendService

func LoggingBackendMiddleware(logger log.Logger) Middleware {
	return func(next BackendService) BackendService {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   BackendService
	logger log.Logger
}

func (mw loggingMiddleware) GetController(ctx context.Context, bid string) (c Controller, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetController", "bid", bid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetController(ctx, bid)
}

func (mw loggingMiddleware) PostCancelActionFeedback(ctx context.Context, bid string,
	fb CancelActionFeedback) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostCancelActionFeedback", "bid", bid, "took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.PostCancelActionFeedback(ctx, bid, fb)
}

func (mw loggingMiddleware) PutConfigData(ctx context.Context, bid string, cfg ConfigData) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PutConfigData", "bid", bid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PutConfigData(ctx, bid, cfg)
}

func (mw loggingMiddleware) GetDeplymentBase(ctx context.Context, bid string,
	acid string) (d DeploymentBase, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetDeplymentBase", "bid", bid, "acid", acid, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetDeplymentBase(ctx, bid, acid)
}

func (mw loggingMiddleware) PostDeploymentBaseFeedback(ctx context.Context, bid string,
	fb DeploymentBaseFeedback) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostDeploymentBaseFeedback", "bid", bid, "took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.PostDeploymentBaseFeedback(ctx, bid, fb)
}

func (mw loggingMiddleware) GetDownloadHttp(ctx context.Context, bid string, ver string) (f []byte, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetDownloadHttp", "bid", bid, "ver", ver, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetDownloadHttp(ctx, bid, ver)
}
