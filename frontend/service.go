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
	GetUpload(ctx context.Context, n string) (deployment.Upload, error)
	// DeleteUpload(ctx context.Context, n string) error
	PostDistribution(ctx context.Context, n string, v string, u string) error
	GetDistribution(ctx context.Context, n string) (deployment.Distribution, error)
	// DeleteDistribution(ctx context.Context, n string) error
	PostDeployment(ctx context.Context, t string, d string) error
	GetDeployment(ctx context.Context, t string) (deployment.Deployment, error)
}

type hawkbitFrontendService struct{}

func NewHawkbitFrontendService() FrontendService {
	return &hawkbitFrontendService{}
}

// PostUpload godoc
//
//	@Summary	Upload new image
//	@Schemes
//	@Description	Upload new image profile which is to be added to a distribution
//	@Tags			Hawkbit FOTA
//	@Param			array	body	frontend.postUploadRequest	false	"New image profile"
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/hawkbit/upload [post]
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

// GetMessage godoc
//
//	@Summary	Retrieve existing upload
//	@Schemes
//	@Description	Retrieve existing upload by specifying upload name
//	@Tags			Hawkbit FOTA
//	@Param			string	path	int	true	"Upload name"
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	deployment.Upload
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/hawkbit/upload/{name} [get]
func (h *hawkbitFrontendService) GetUpload(ctx context.Context, n string) (deployment.Upload, error) {
	u, err := deployment.GetUpload(n)
	if err != nil {
		return deployment.Upload{}, err
	}
	return u, nil
}

// PostDistribution godoc
//
//	@Summary	Create new distribution
//	@Schemes
//	@Description	Create new distribution with upload specified which is to be added to a deployment
//	@Tags			Hawkbit FOTA
//	@Param			array	body	frontend.postDistributionRequest	false	"New distribution"
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/hawkbit/dist [post]
func (h *hawkbitFrontendService) PostDistribution(ctx context.Context, n string, v string, u string) error {
	var d deployment.Distribution
	d.Name = n
	d.Version = v
	if err := deployment.SetDistribution(d, u); err != nil {
		return ErrFrontendDistribution
	}
	return nil
}

// GetDistribution godoc
//
//	@Summary	Retrieve existing distribution
//	@Schemes
//	@Description	Retrieve existing distribution by specifying distribution name
//	@Tags			Hawkbit FOTA
//	@Param			string	path	int	true	"Distribution name"
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	deployment.Distribution
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/hawkbit/dist/{name} [get]
func (h *hawkbitFrontendService) GetDistribution(ctx context.Context, n string) (deployment.Distribution, error) {
	d, err := deployment.GetDistribution(n)
	if err != nil {
		return deployment.Distribution{}, err
	}
	return d, nil
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

func (h *hawkbitFrontendService) GetDeployment(ctx context.Context, t string) (deployment.Deployment, error) {
	dp, err := deployment.GetDeployment(t)
	if err != nil {
		return deployment.Deployment{}, err
	}
	return dp, nil
}
