package frontend

import (
	"context"
	"errors"

	"github.com/jonathanyhliang/hawkbit-fota/deployment"
)

var (
	ErrFrontendUpload       = errors.New("Frontend: image upload failed")
	ErrFrontendDistribution = errors.New("Frontend: distribution set failed")
	ErrFrontendDeployment   = errors.New("Frontend: deployment set failed")
	ErrFrontendBadRequest   = errors.New("Frontend: bad request")
)

type FrontendService interface {
	PostUpload(ctx context.Context, n string, v string, f string) error
	PostDistribution(ctx context.Context, n string, v string, u string) error
	PostDeployment(ctx context.Context, t string, d string) error
}

type hawkbitFrontendService struct {
}

func NewHawkbitFrontendService() FrontendService {
	return &hawkbitFrontendService{}
}

func (h *hawkbitFrontendService) PostUpload(ctx context.Context, n string, v string, f string) error {
	var u deployment.Upload
	u.Name = n
	u.Version = v
	u.File = f
	if err := deployment.SetUpload(u); err != nil {
		return ErrFrontendUpload
	}
	return nil
}

func (h *hawkbitFrontendService) PostDistribution(ctx context.Context, n string, v string, u string) error {
	var d deployment.Distribution
	d.Name = n
	d.Version = v
	if err := deployment.SetDistribution(d, u); err != nil {
		return ErrFrontendDistribution
	}
	return nil
}

func (h *hawkbitFrontendService) PostDeployment(ctx context.Context, t string, d string) error {
	if t == "" {
		return ErrFrontendBadRequest
	}
	if err := deployment.SetDeployment(t, d); err != nil {
		return ErrFrontendDeployment
	}
	return nil
}
