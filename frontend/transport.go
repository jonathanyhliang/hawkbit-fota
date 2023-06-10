package frontend

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
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeFrontendHTTPHandler(s FrontendService, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeFrontendServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/hawkbit/upload").Handler(httptransport.NewServer(
		e.PostUpload,
		decodePostUploadEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/hawkbit/upload/{name}").Handler(httptransport.NewServer(
		e.GetUpload,
		decodeGetUploadEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/hawkbit/dist").Handler(httptransport.NewServer(
		e.PostDistribution,
		decodePostDistributionEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/hawkbit/dist/{name}").Handler(httptransport.NewServer(
		e.GetDistribution,
		decodeGetDistributionEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/hawkbit/deploy").Handler(httptransport.NewServer(
		e.PostDeployment,
		decodePostDeploymentEndpoint,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/hawkbit/deploy/{target}").Handler(httptransport.NewServer(
		e.GetDeployment,
		decodeGetDeploymentEndpoint,
		encodeResponse,
		options...,
	))
	r.PathPrefix("/hawkbit/docs").Handler(httpSwagger.WrapHandler)
	return r
}

func decodePostUploadEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	var u postUploadRequest
	e := json.NewDecoder(r.Body).Decode(&u)
	if e != nil {
		return nil, e
	}
	return postUploadRequest{Name: u.Name, Version: u.Version, File: u.File}, nil
}

func decodeGetUploadEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	n, ok := vars["name"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getUploadRequest{Name: n}, nil
}

func decodePostDistributionEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	var d postDistributionRequest
	e := json.NewDecoder(r.Body).Decode(&d)
	if e != nil {
		return nil, e
	}
	return postDistributionRequest{Name: d.Name, Version: d.Version, Upload: d.Upload}, nil
}

func decodeGetDistributionEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	n, ok := vars["name"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getDistributionRequest{Name: n}, nil
}

func decodePostDeploymentEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	var d postDeploymentRequest
	e := json.NewDecoder(r.Body).Decode(&d)
	if e != nil {
		return nil, e
	}
	return postDeploymentRequest{Target: d.Target, Distribution: d.Distribution}, nil
}

func decodeGetDeploymentEndpoint(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	t, ok := vars["target"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getDeploymentRequest{Target: t}, nil
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
	case deployment.ErrDeploymentNotFound,
		deployment.ErrDeploymentUploadNotFound,
		deployment.ErrDeploymentDistNotFound:
		return http.StatusNotFound
	case ErrFrontendUpload,
		ErrFrontendDistribution,
		ErrFrontendDeployment:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
