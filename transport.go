package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// GET	/default/controller/v1/:bid  							retrieves actions that need to be executed
	// POST	/default/controller/v1/:bid/cancelAction/:acid/feedback	post action cancellation result
	// POST /default/controller/v1/:bid/configData       			post hardware level identification of the target
	// GET 	/default/controller/v1/:bid/deploymentBase/:acid   		retrieves core resource for deployment operations
	// GET  /DEFAULT/controller/v1/:bid/softwareModules/:ver   		retrieves artifcat for update

	r.Methods("GET").Path("/default/controller/v1/{bid}").Handler(httptransport.NewServer(
		e.GetControllerEndpoint,
		decodeGetControllerEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/default/controller/v1/{bid}/cancelAction/{acid}/feedback").Handler(httptransport.NewServer(
		e.PostCancelActionFeebackEndpoint,
		decodePostCancelActionFeebackEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/default/controller/v1/{bid}/configData").Handler(httptransport.NewServer(
		e.PostConfigDataEndpoint,
		decodePostConfigDataEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/default/controller/v1/{bid}/deploymentBase/{acid}").Handler(httptransport.NewServer(
		e.GetDeploymentBaseEndpoint,
		decodeGetDeploymentBaseEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/default/controller/v1/{bid}/deploymentBase/{acid}/feedback").Handler(httptransport.NewServer(
		e.PostDeploymentBaseFeedbackEndpoint,
		decodePostDeploymentBaseFeedbackEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/DEFAULT/controller/v1/{bid}/softwareModules/{ver}").Handler(httptransport.NewServer(
		e.GetDownloadHttpEndpoint,
		decodeGetDownloadHttpEndpoint,
		encodeDownloadHttpResponse,
		options...,
	))
	return r
}

func decodeGetControllerEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	bid, ok := vars["bid"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getControllerRequest{BID: bid}, nil
}

func decodePostCancelActionFeebackEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	bid, ok := vars["bid"]
	if !ok {
		return nil, ErrBadRouting
	}
	acid, ok := vars["acid"]
	if !ok {
		return nil, ErrBadRouting
	}
	var fb CancelActionFeedback
	if e := json.NewDecoder(r.Body).Decode(&fb); e != nil {
		return nil, e
	}
	return postCancelActionFeedbackRequest{BID: bid, ACID: acid, Fb: fb}, nil
}

func decodePostConfigDataEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	bid, ok := vars["bid"]
	if !ok {
		return nil, ErrBadRouting
	}
	var cfg ConfigData
	if e := json.NewDecoder(r.Body).Decode(&cfg); e != nil {
		return nil, e
	}
	return postConfigDataRequest{BID: bid, Cfg: cfg}, nil
}

func decodeGetDeploymentBaseEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	bid, ok := vars["bid"]
	if !ok {
		return nil, ErrBadRouting
	}
	acid, ok := vars["acid"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getDeplymentBaseRequest{BID: bid, ACID: acid}, nil
}

func decodePostDeploymentBaseFeedbackEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	bid, ok := vars["bid"]
	if !ok {
		return nil, ErrBadRouting
	}
	acid, ok := vars["acid"]
	if !ok {
		return nil, ErrBadRouting
	}
	var fb DeploymentBaseFeedback
	if e := json.NewDecoder(r.Body).Decode(&fb); e != nil {
		return nil, e
	}
	return postDeploymentBaseFeedbackRequest{BID: bid, ACID: acid, Fb: fb}, nil
}

func decodeGetDownloadHttpEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	bid, ok := vars["bid"]
	if !ok {
		return nil, ErrBadRouting
	}
	ver, ok := vars["ver"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getDownloadHttpRequest{BID: bid, Ver: ver}, nil
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
	case ErrDeploymentNotFound:
		return http.StatusNotFound
	case ErrDownloadImage:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
