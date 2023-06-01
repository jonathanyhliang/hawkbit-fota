package backend

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetControllerEndpoint              endpoint.Endpoint
	PostCancelActionFeebackEndpoint    endpoint.Endpoint
	PostConfigDataEndpoint             endpoint.Endpoint
	GetDeploymentBaseEndpoint          endpoint.Endpoint
	PostDeploymentBaseFeedbackEndpoint endpoint.Endpoint
	GetDownloadHttpEndpoint            endpoint.Endpoint
}

func MakeBackendServerEndpoints(s BackendService) Endpoints {
	return Endpoints{
		GetControllerEndpoint:              MakeGetControllerEndpoint(s),
		PostCancelActionFeebackEndpoint:    MakePostCancelActionFeedbackEndpoint(s),
		PostConfigDataEndpoint:             MakePostConfigDataEndpoint(s),
		GetDeploymentBaseEndpoint:          MakeGetDeploymentBaseEndpoint(s),
		PostDeploymentBaseFeedbackEndpoint: MakePostDeploymentBaseFeedbackEndpoint(s),
		GetDownloadHttpEndpoint:            MakeGetDownloadHttpEndpoint(s),
	}
}

func MakeGetControllerEndpoint(s BackendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getControllerRequest)
		c, e := s.GetController(ctx, req.BID)
		return getControllerResponse{Ctrlr: c, Err: e}, nil
	}
}

func MakePostCancelActionFeedbackEndpoint(s BackendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postCancelActionFeedbackRequest)
		e := s.PostCancelActionFeedback(ctx, req.BID, req.Fb)
		return postCancelActionFeedbackResponse{Err: e}, nil
	}
}

func MakePostConfigDataEndpoint(s BackendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postConfigDataRequest)
		e := s.PostConfigData(ctx, req.BID, req.Cfg)
		return postConfigDataResponse{Err: e}, nil
	}
}

func MakeGetDeploymentBaseEndpoint(s BackendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getDeplymentBaseRequest)
		d, e := s.GetDeplymentBase(ctx, req.BID, req.ACID)
		return getDeplymentBaseResponse{Dp: d, Err: e}, nil
	}
}

func MakePostDeploymentBaseFeedbackEndpoint(s BackendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postDeploymentBaseFeedbackRequest)
		e := s.PostDeploymentBaseFeedback(ctx, req.BID, req.Fb)
		return postDeploymentBaseFeedbackResponse{Err: e}, nil
	}
}

func MakeGetDownloadHttpEndpoint(s BackendService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getDownloadHttpRequest)
		f, e := s.GetDownloadHttp(ctx, req.BID, req.Ver)
		return getDownloadHttpResponse{File: f, Err: e}, nil
	}
}

type getControllerRequest struct {
	BID string
}

type getControllerResponse struct {
	Ctrlr Controller `json:"controller"`
	Err   error      `json:"err,omitempty"`
}

func (r getControllerResponse) error() error { return r.Err }

type postCancelActionFeedbackRequest struct {
	BID string
	Fb  CancelActionFeedback `json:"cancelActionFeedback,omitempty"`
}

type postCancelActionFeedbackResponse struct {
	Err error `json:"err,omitempty"`
}

func (r postCancelActionFeedbackResponse) error() error { return r.Err }

type postConfigDataRequest struct {
	BID string
	Cfg ConfigData `json:"configData,omitempty"`
}

type postConfigDataResponse struct {
	Err error `json:"err,omitempty"`
}

func (r postConfigDataResponse) error() error { return r.Err }

type getDeplymentBaseRequest struct {
	BID  string
	ACID string
}

type getDeplymentBaseResponse struct {
	Dp  DeploymentBase `json:"deploymentBase,omitempty"`
	Err error          `json:"err,omitempty"`
}

type postDeploymentBaseFeedbackRequest struct {
	BID string
	Fb  DeploymentBaseFeedback `json:"deploymentBaseFeedback,omitempty"`
}

type postDeploymentBaseFeedbackResponse struct {
	Err error `json:"err,omitempty"`
}

func (r getDeplymentBaseResponse) error() error { return r.Err }

type getDownloadHttpRequest struct {
	BID  string
	ACID string
	Ver  string
}

type getDownloadHttpResponse struct {
	File []byte
	Err  error
}

func (r getDownloadHttpResponse) file() []byte { return r.File }

func (r getDownloadHttpResponse) error() error { return r.Err }
