package admission

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/jondot/lightscreen/pkg/actions"
	"go.uber.org/zap"
)

type ServerOptions struct {
	Config         string
	Development    bool
	Port           int
	Host           string
	MetricsAddress string
	CertDir        string
}

type Server struct {
	opts    ServerOptions
	Actions *actions.Actions
	logger  *zap.SugaredLogger
}
type hookWrapper struct {
	client   client.Client
	decoder  *admission.Decoder
	workflow *AdmissionWorkflow
	logger   *zap.SugaredLogger
}

// hookWrapper admits a pod iff a specific annotation exists.
func (v *hookWrapper) Handle(ctx context.Context, req admission.Request) admission.Response {
	return v.workflow.Execute(ctx, req.Object.Raw)
}

// hookWrapper implements inject.Client.
// A client will be automatically injected.

// InjectClient injects the client.
func (v *hookWrapper) InjectClient(c client.Client) error {
	v.client = c
	return nil
}

// hookWrapper implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (v *hookWrapper) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
func NewServer(opts ServerOptions, logger *zap.SugaredLogger) *Server {
	return &Server{
		Actions: actions.Default(),
		opts:    opts,
		logger:  logger,
	}
}

func (s *Server) prepareWorkflow() (*AdmissionWorkflow, error) {

	logger := s.logger
	f, err := os.Open(s.opts.Config)
	if err != nil {
		logger.Fatalw("Cannot open file", "error", err)
	}

	return NewAdmissionWorkflow(s.Actions, f, logger)
}

func (s *Server) Check(fname string) (*admission.Response, error) {
	workflow, err := s.prepareWorkflow()
	if err != nil {
		s.logger.Fatal(err)
		return nil, err
	}
	out, err := ioutil.ReadFile(fname)
	if err != nil {
		s.logger.Fatal(err)
		return nil, err
	}
	resp := workflow.Execute(context.TODO(), out)
	return &resp, nil
}
func (s *Server) Serve() {
	workflow, err := s.prepareWorkflow()
	if err != nil {
		s.logger.Fatal(err)
		return
	}
	config, err := config.GetConfig()
	if err != nil {
		s.logger.Fatal(err)
		return
	}
	mgr, err := manager.New(config, manager.Options{
		MetricsBindAddress: s.opts.MetricsAddress,
		Port:               s.opts.Port,
		Host:               s.opts.Host,
	})
	if err != nil {
		s.logger.Error(err, "unable to create manager")
		os.Exit(1)
	}
	hookServer := mgr.GetWebhookServer()
	hookServer.CertDir = s.opts.CertDir
	hookServer.Register("/", &webhook.Admission{Handler: &hookWrapper{
		workflow: workflow,
		logger:   s.logger,
	}})
	hookServer.WebhookMux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	})
	s.logger.Infof("Lightscreen on %v:%v (metrics [%v/metrics]", s.opts.Host, s.opts.Port, s.opts.MetricsAddress)
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		s.logger.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
