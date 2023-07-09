package frontend

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/jonathanyhliang/hawkbit-fota/deployment"
)

type Endpoints struct {
	PostUpload       endpoint.Endpoint
	GetUpload        endpoint.Endpoint
	PostDistribution endpoint.Endpoint
	GetDistribution  endpoint.Endpoint
	PostDeployment   endpoint.Endpoint
	GetDeployment    endpoint.Endpoint
}

func MakeFrontendServerEndpoints(s FrontendService) Endpoints {
	return Endpoints{
		PostUpload:       MakePostUpload(s),
		GetUpload:        MakeGetUpload(s),
		PostDistribution: MakePostDistribution(s),
		GetDistribution:  MakeGetDistribution(s),
		PostDeployment:   MakePostDeployment(s),
		GetDeployment:    MakeGetDeployment(s),
	}
}

func MakePostUpload(s FrontendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postUploadRequest)
		e := s.PostUpload(ctx, req.Name, req.Version, req.File)
		return postUploadResponse{Err: e}, nil
	}
}

func MakeGetUpload(s FrontendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getUploadRequest)
		u, e := s.GetUpload(ctx, req.Name)
		return getUploadResponse{Upload: u, Err: e}, nil
	}
}

func MakePostDistribution(s FrontendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postDistributionRequest)
		e := s.PostDistribution(ctx, req.Name, req.Version, req.Upload)
		return postDistributionResponse{Err: e}, nil
	}
}

func MakeGetDistribution(s FrontendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getDistributionRequest)
		d, e := s.GetDistribution(ctx, req.Name)
		return getDistributionResponse{Distribution: d, Err: e}, nil
	}
}

func MakePostDeployment(s FrontendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postDeploymentRequest)
		e := s.PostDeployment(ctx, req.Target, req.Distribution)
		return postDeploymentResponse{Err: e}, nil
	}
}

func MakeGetDeployment(s FrontendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getDeploymentRequest)
		dp, e := s.GetDeployment(ctx, req.Target)
		return getDeploymentResponse{Deployment: dp, Err: e}, nil
	}
}

type postUploadRequest struct {
	Name    string `json:"name" example:"zephyr_cc3220sf_signed"`
	Version string `json:"version" example:"1.0.0+1"`
	File    string `json:"file" example:"/workdir/build/artifact.bin"`
}

type postUploadResponse struct {
	Err error `json:"error,omitempty"`
}

type getUploadRequest struct {
	Name string `json:"name"`
}

type getUploadResponse struct {
	Upload deployment.Upload `json:"upload,omitempty"`
	Err    error             `json:"error,omitempty"`
}

type postDistributionRequest struct {
	Name    string `json:"name" example:"hawkbit"`
	Version string `json:"version" example:"1.0.0+1"`
	Upload  string `json:"upload" example:"zephyr_cc3220sf_signed"`
}

type postDistributionResponse struct {
	Err error `json:"error,omitempty"`
}

type getDistributionRequest struct {
	Name string `json:"name"`
}

type getDistributionResponse struct {
	Distribution deployment.Distribution `json:"distribution,omitempty"`
	Err          error                   `json:"error,omitempty"`
}

type postDeploymentRequest struct {
	Target       string `json:"target" example:"ti_cc3200wf_12345"`
	Distribution string `json:"distribution" example:"hawkbit"`
}

type postDeploymentResponse struct {
	Err error `json:"error,omitempty"`
}

type getDeploymentRequest struct {
	Target string `json:"target"`
}

type getDeploymentResponse struct {
	Deployment deployment.Deployment `json:"deployment,omitempty"`
	Err        error                 `json:"error,omitempty"`
}
