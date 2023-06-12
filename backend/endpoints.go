package backend

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetControllerEndpoint              endpoint.Endpoint
	PostCancelActionFeebackEndpoint    endpoint.Endpoint
	PutConfigDataEndpoint              endpoint.Endpoint
	GetDeploymentBaseEndpoint          endpoint.Endpoint
	PostDeploymentBaseFeedbackEndpoint endpoint.Endpoint
	GetDownloadHttpEndpoint            endpoint.Endpoint
}

func MakeBackendServerEndpoints(s BackendService) Endpoints {
	return Endpoints{
		GetControllerEndpoint:              MakeGetControllerEndpoint(s),
		PostCancelActionFeebackEndpoint:    MakePostCancelActionFeedbackEndpoint(s),
		PutConfigDataEndpoint:              MakePutConfigDataEndpoint(s),
		GetDeploymentBaseEndpoint:          MakeGetDeploymentBaseEndpoint(s),
		PostDeploymentBaseFeedbackEndpoint: MakePostDeploymentBaseFeedbackEndpoint(s),
		GetDownloadHttpEndpoint:            MakeGetDownloadHttpEndpoint(s),
	}
}

func MakeGetControllerEndpoint(s BackendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetControllerRequest)
		c, e := s.GetController(ctx, req.Bid)
		return GetControllerResponse{Ctrlr: c, Err: e}, nil
	}
}

func MakePostCancelActionFeedbackEndpoint(s BackendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PostCancelActionFeedbackRequest)
		e := s.PostCancelActionFeedback(ctx, req.Bid, req.Fb)
		return PostCancelActionFeedbackResponse{Err: e}, nil
	}
}

func MakePutConfigDataEndpoint(s BackendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PutConfigDataRequest)
		e := s.PutConfigData(ctx, req.Bid, req.Cfg)
		return PutConfigDataResponse{Err: e}, nil
	}
}

func MakeGetDeploymentBaseEndpoint(s BackendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetDeplymentBaseRequest)
		d, e := s.GetDeplymentBase(ctx, req.Bid, req.Acid)
		return GetDeplymentBaseResponse{Dp: d, Err: e}, nil
	}
}

func MakePostDeploymentBaseFeedbackEndpoint(s BackendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PostDeploymentBaseFeedbackRequest)
		e := s.PostDeploymentBaseFeedback(ctx, req.Bid, req.Fb)
		return PostDeploymentBaseFeedbackResponse{Err: e}, nil
	}
}

func MakeGetDownloadHttpEndpoint(s BackendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetDownloadHttpRequest)
		f, e := s.GetDownloadHttp(ctx, req.Bid, req.Ver)
		return GetDownloadHttpResponse{File: f, Err: e}, nil
	}
}

type GetControllerRequest struct {
	Bid string
}

type GetControllerResponse struct {
	Ctrlr Controller `json:"controller"`
	Err   error      `json:"err,omitempty"`
}

func (r GetControllerResponse) error() error { return r.Err }

type PostCancelActionFeedbackRequest struct {
	Bid string
	Fb  CancelActionFeedback `json:"cancelActionFeedback,omitempty"`
}

type PostCancelActionFeedbackResponse struct {
	Err error `json:"err,omitempty"`
}

func (r PostCancelActionFeedbackResponse) error() error { return r.Err }

type PutConfigDataRequest struct {
	Bid string
	Cfg ConfigData `json:"configData,omitempty"`
}

type PutConfigDataResponse struct {
	Err error `json:"err,omitempty"`
}

func (r PutConfigDataResponse) error() error { return r.Err }

type GetDeplymentBaseRequest struct {
	Bid  string
	Acid string
}

type GetDeplymentBaseResponse struct {
	Dp  DeploymentBase `json:"deploymentBase,omitempty"`
	Err error          `json:"err,omitempty"`
}

type PostDeploymentBaseFeedbackRequest struct {
	Bid string
	Fb  DeploymentBaseFeedback `json:"deploymentBaseFeedback,omitempty"`
}

type PostDeploymentBaseFeedbackResponse struct {
	Err error `json:"err,omitempty"`
}

func (r GetDeplymentBaseResponse) error() error { return r.Err }

type GetDownloadHttpRequest struct {
	Bid string
	Ver string
}

type GetDownloadHttpResponse struct {
	File []byte
	Err  error
}

func (r GetDownloadHttpResponse) file() []byte { return r.File }

func (r GetDownloadHttpResponse) error() error { return r.Err }
