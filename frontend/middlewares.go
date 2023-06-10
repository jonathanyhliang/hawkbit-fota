package frontend

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jonathanyhliang/hawkbit-fota/deployment"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(FrontendService) FrontendService

func LoggingFrontendMiddleware(logger log.Logger) Middleware {
	return func(next FrontendService) FrontendService {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   FrontendService
	logger log.Logger
}

func (mw loggingMiddleware) PostUpload(ctx context.Context, n string, v string, f string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostUpload", "name", n, "version", v, "file", f,
			"took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostUpload(ctx, n, v, f)
}

func (mw loggingMiddleware) GetUpload(ctx context.Context, n string) (u deployment.Upload, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetUpload", "name", n, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetUpload(ctx, n)
}

func (mw loggingMiddleware) PostDistribution(ctx context.Context, n string, v string, u string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostDistribution", "name", n, "version", v, "upload", u,
			"took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostDistribution(ctx, n, v, u)
}

func (mw loggingMiddleware) GetDistribution(ctx context.Context, n string) (d deployment.Distribution, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetDistribution", "name", n, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetDistribution(ctx, n)
}

func (mw loggingMiddleware) PostDeployment(ctx context.Context, t string, d string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostDeployment", "target", t, "distribution", d,
			"took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostDeployment(ctx, t, d)
}

func (mw loggingMiddleware) GetDeployment(ctx context.Context, t string) (dp deployment.Deployment, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetDeployment", "target", t, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetDeployment(ctx, t)
}
