package main

import (

	//	. "github.com/ahmetb/go-linq" // From

	"github.com/jondot/lightscreen/pkg/admission"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	config     = kingpin.Flag("config", "Action mapping configuration file").Default("lightscreen.yaml").OverrideDefaultFromEnvar("LS_CONFIG").Short('c').String()
	input      = kingpin.Flag("check", "Check input file for admission").String()
	host       = kingpin.Flag("host", "HTTP host").OverrideDefaultFromEnvar("LS_HOST").Default("0.0.0.0").String()
	port       = kingpin.Flag("port", "HTTP port").OverrideDefaultFromEnvar("LS_PORT").Default("443").Int()
	metrics    = kingpin.Flag("metrics", "Metrics HTTP port").OverrideDefaultFromEnvar("LS_METRICSPORT").Default(":8080").String()
	production = kingpin.Flag("production", "Run in production mode").OverrideDefaultFromEnvar("LS_PROD").Short('p').Bool()
	certs      = kingpin.Flag("certs", "Certs dir").OverrideDefaultFromEnvar("LS_CERTS").Default("self-certs").String()
)

func main() {
	kingpin.Parse()

	development := !(*production)

	log_, _ := zap.NewProduction()
	if development {
		log_, _ = zap.NewDevelopment()
	}
	logger := log_.Sugar()
	server := admission.NewServer(admission.ServerOptions{
		CertDir:        *certs,
		Config:         *config,
		Development:    development,
		Port:           *port,
		Host:           *host,
		MetricsAddress: *metrics,
	}, logger)
	//server.Actions.
	logger.Infof("mutations=%v validations=%v", server.Actions.MutationMap, server.Actions.ValidationMap)
	logger.Infow("running", "port", *port, "host", *host, "config", *config, "development", development)
	if *input != "" {
		resp, err := server.Check(*input)
		if err != nil {
			logger.Errorf("error %v", err)
		} else {
			logger.Infow("", "response", resp)
		}

	} else {
		server.Serve()
	}
}
