package diskcollector

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const (
	MetricsPort = "METRICS_PORT"
	Address     = "METRICS_ADDRESS"
	DefaultPort = 8080
)

type serverSettings struct {
	port     *int
	addr     string
	certFile string
	keyFile  string
}

func (s *serverSettings) getAddr() string {
	if s.port == nil {
		s.port = new(int)
		*s.port = DefaultPort
	}
	return fmt.Sprintf("%s:%d", s.addr, *s.port)
}

func InitPrometheusRegistry(mountPath string) *prometheus.Registry {
	udf := InitPd(mountPath)
	ch := &CollectHandler{}
	ch.Register(udf)
	customRegistry := prometheus.NewRegistry()
	customRegistry.MustRegister(ch)
	return customRegistry
}

func ServeHTTP(registry *prometheus.Registry, settings *serverSettings) {
	addr := settings.getAddr()
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}
	log.Info("Starting server on ", addr)
	if err := s.ListenAndServe(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Webserver encountered an error")
	}
}

func ServeHTTPS(registry *prometheus.Registry, settings *serverSettings) {
	addr := settings.getAddr()
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}
	log.Info("Starting server on ", addr)
	if err := s.ListenAndServeTLS(settings.certFile, settings.keyFile); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Webserver encountered an error")
	}
}

func ServeCollector() {
	// Parse ENV flags
	port := os.Getenv(MetricsPort)
	addr := os.Getenv(Address)
	settings := &serverSettings{}
	if port != "" {
		convertedPort, _ := strconv.Atoi(port)
		settings.port = &convertedPort
	}
	if addr != "" {
		settings.addr = addr
	}
	// Parse user input
	mountPath := flag.String("path", "", "Path that is to be monitored with statfs")
	verboseOutput := flag.Bool("v", false, "Verbose output (debug)")
	certFile := flag.String("certFile", "", "If serving TLS provide certfile")
	keyFile := flag.String("keyFile", "", "If serving TLS provde keyfile")
	flag.Parse()

	// Setup log
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	if *verboseOutput {
		log.SetLevel(log.DebugLevel)
	}

	// Setup collector for prometheus registry and serve HTTP(S)
	registry := InitPrometheusRegistry(*mountPath)
	if *certFile != "" && *keyFile != "" {
		log.WithFields(log.Fields{
			"keyFile":  *keyFile,
			"certFile": *certFile,
		}).Debug("Key and cert file supplied, trying to serve TLS")

		settings.keyFile = *keyFile
		settings.certFile = *certFile
		ServeHTTPS(registry, settings)
	} else {
		ServeHTTP(registry, settings)
	}
}
