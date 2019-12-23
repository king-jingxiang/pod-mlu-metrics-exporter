package main

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

const mluMetricsEndpoint = "/mlu/metrics"

func newHttpServer(addr string) *http.Server {
	mux := http.NewServeMux()

	s := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  connectionTimeout,
		WriteTimeout: connectionTimeout,
	}

	mux.HandleFunc(mluMetricsEndpoint, getmlumetrics)
	return s
}

func startHttp(s *http.Server) {
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		glog.Fatalf("error listening requests at %v: %v", s.Addr, err)
	}
}

func stopHttp(s *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		glog.Errorf("error shutting down the http server: %v", err)
	} else {
		glog.Info("http server stopped")
	}
}

func getmlumetrics(resp http.ResponseWriter, req *http.Request) {
	metrics, err := ioutil.ReadFile(mluPodMetrics)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		glog.Errorf("error responding to %v%v: %v", req.Host, req.URL, err.Error())
		return
	}
	resp.Write(metrics)
	process_metrics, err := ioutil.ReadFile(deviceResourceMetrics)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		glog.Errorf("error responding to %v%v: %v", req.Host, req.URL, err.Error())
		return
	}
	resp.Write(process_metrics)
}
