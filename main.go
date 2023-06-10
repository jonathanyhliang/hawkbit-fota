package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/jonathanyhliang/hawkbit-fota/backend"
	"github.com/jonathanyhliang/hawkbit-fota/docs"
	"github.com/jonathanyhliang/hawkbit-fota/frontend"
)

//	@title		Hawkbit FOTA Service API
//	@version	1.0

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host	localhost:port/hawkbit

func main() {
	var (
		BackendAddr  = flag.String("ba", "192.168.179.5:8080", "Backend HTTP listen address")
		FrontendAddr = flag.String("fa", ":8080", "Frontend HTTP listen address")
	)
	flag.Parse()

	errs := make(chan error)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var bs backend.BackendService
	{
		bs = backend.NewHawkbitBackendService()
		bs = backend.LoggingBackendMiddleware(logger)(bs)
	}

	var fs frontend.FrontendService
	{
		fs = frontend.NewHawkbitFrontendService()
		fs = frontend.LoggingFrontendMiddleware(logger)(fs)
	}

	var bh http.Handler
	{
		bh = backend.MakeBackendHTTPHandler(bs, log.With(logger, "component", "HTTP"))
	}

	var fh http.Handler
	{
		fh = frontend.MakeFrontendHTTPHandler(fs, log.With(logger, "component", "HTTP"))
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("backend", "HTTP", "addr", *BackendAddr)
		errs <- http.ListenAndServe(*BackendAddr, bh)
	}()

	go func() {
		docs.SwaggerInfo.BasePath = "/"
		logger.Log("frontend", "HTTP", "addr", *FrontendAddr)
		errs <- http.ListenAndServe(*FrontendAddr, fh)
	}()

	logger.Log("exit", <-errs)
}
