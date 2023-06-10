package backend

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jonathanyhliang/hawkbit-fota/deployment"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeBackendHTTPHandler(s BackendService, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeBackendServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// Get	/default/controller/v1/:Bid  							retrieves actions that need to be executed
	// Post	/default/controller/v1/:Bid/cancelAction/:Acid/feedback	Post action cancellation result
	// Post /default/controller/v1/:Bid/configData       			Post hardware level identification of the tarGet
	// Get 	/default/controller/v1/:Bid/deploymentBase/:Acid   		retrieves core resource for deployment operations
	// Get  /DEFAULT/controller/v1/:Bid/softwareModules/:ver   		retrieves artifcat for update

	r.Methods("Get").Path("/default/controller/v1/{Bid}").Handler(httptransport.NewServer(
		e.GetControllerEndpoint,
		decodeGetControllerEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("Post").Path("/default/controller/v1/{Bid}/cancelAction/{Acid}/feedback").Handler(httptransport.NewServer(
		e.PostCancelActionFeebackEndpoint,
		decodePostCancelActionFeebackEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("Post").Path("/default/controller/v1/{Bid}/configData").Handler(httptransport.NewServer(
		e.PostConfigDataEndpoint,
		decodePostConfigDataEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("Get").Path("/default/controller/v1/{Bid}/deploymentBase/{Acid}").Handler(httptransport.NewServer(
		e.GetDeploymentBaseEndpoint,
		decodeGetDeploymentBaseEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("Post").Path("/default/controller/v1/{Bid}/deploymentBase/{Acid}/feedback").Handler(httptransport.NewServer(
		e.PostDeploymentBaseFeedbackEndpoint,
		decodePostDeploymentBaseFeedbackEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("Get").Path("/DEFAULT/controller/v1/{Bid}/softwareModules/{ver}").Handler(httptransport.NewServer(
		e.GetDownloadHttpEndpoint,
		decodeGetDownloadHttpEndpoint,
		encodeDownloadHttpResponse,
		options...,
	))
	return r
}

func decodeGetControllerEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	Bid, ok := vars["Bid"]
	if !ok {
		return nil, ErrBadRouting
	}
	return GetControllerRequest{Bid: Bid}, nil
}

func decodePostCancelActionFeebackEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	Bid, ok := vars["Bid"]
	if !ok {
		return nil, ErrBadRouting
	}
	var fb CancelActionFeedback
	if e := json.NewDecoder(r.Body).Decode(&fb); e != nil {
		return nil, e
	}
	return PostCancelActionFeedbackRequest{Bid: Bid, Fb: fb}, nil
}

func decodePostConfigDataEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	Bid, ok := vars["Bid"]
	if !ok {
		return nil, ErrBadRouting
	}
	var cfg ConfigData
	if e := json.NewDecoder(r.Body).Decode(&cfg); e != nil {
		return nil, e
	}
	return PostConfigDataRequest{Bid: Bid, Cfg: cfg}, nil
}

func decodeGetDeploymentBaseEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	Bid, ok := vars["Bid"]
	if !ok {
		return nil, ErrBadRouting
	}
	Acid, ok := vars["Acid"]
	if !ok {
		return nil, ErrBadRouting
	}
	return GetDeplymentBaseRequest{Bid: Bid, Acid: Acid}, nil
}

func decodePostDeploymentBaseFeedbackEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	Bid, ok := vars["Bid"]
	if !ok {
		return nil, ErrBadRouting
	}
	var fb DeploymentBaseFeedback
	if e := json.NewDecoder(r.Body).Decode(&fb); e != nil {
		return nil, e
	}
	return PostDeploymentBaseFeedbackRequest{Bid: Bid, Fb: fb}, nil
}

func decodeGetDownloadHttpEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	Bid, ok := vars["Bid"]
	if !ok {
		return nil, ErrBadRouting
	}
	ver, ok := vars["ver"]
	if !ok {
		return nil, ErrBadRouting
	}
	return GetDownloadHttpRequest{Bid: Bid, Ver: ver}, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type filer interface {
	file() []byte
}

func encodeDownloadHttpResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	if f, ok := response.(filer); ok {
		w.Header().Set("Content-Type", "application/text; charset=utf-8")
		w.Write(f.file())
		return nil
	}
	return ErrBadRouting
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case deployment.ErrDeploymentNotFound:
		return http.StatusNotFound
	case ErrBackendBadRequest, ErrBackendDownload:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
