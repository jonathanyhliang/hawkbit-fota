package frontend

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	PostUpload       endpoint.Endpoint
	PostDistribution endpoint.Endpoint
	PostDeployment   endpoint.Endpoint
}

func MakeFrontendServerEndpoints(s FrontendService) Endpoints {
	return Endpoints{
		PostUpload:       MakePostUpload(s),
		PostDistribution: MakePostDistribution(s),
		PostDeployment:   MakePostDeployment(s),
	}
}

func MakePostUpload(s FrontendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postUploadRequest)
		e := s.PostUpload(ctx, req.Name, req.Version, req.File)
		return postUploadResponse{Err: e}, nil
	}
}

func MakePostDistribution(s FrontendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postDistributionRequest)
		e := s.PostDistribution(ctx, req.Name, req.Version, req.Upload)
		return postDistributionResponse{Err: e}, nil
	}
}

func MakePostDeployment(s FrontendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postDeploymentRequest)
		e := s.PostDeployment(ctx, req.Target, req.Distribution)
		return postDeploymentResponse{Err: e}, nil
	}
}

type postUploadRequest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	File    string `json:"file"`
}

type postUploadResponse struct {
	Err error `json:"error,omitempty"`
}

type postDistributionRequest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Upload  string `json:"upload"`
}

type postDistributionResponse struct {
	Err error `json:"error,omitempty"`
}

type postDeploymentRequest struct {
	Target       string `json:"target"`
	Distribution string `json:"distribution"`
}

type postDeploymentResponse struct {
	Err error `json:"error,omitempty"`
}
