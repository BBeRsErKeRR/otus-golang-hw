package integration_test

import (
	"flag"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	grpcAddr, amqpAddr, rootHTTPURL string
	sleepDuration                   time.Duration
)

func init() {
	flag.StringVar(&rootHTTPURL, "http-addr", "http://localhost:5000", "Address of the http server to smoke-check")
	flag.StringVar(&grpcAddr, "grpc-addr", "0.0.0.0:5080", "Address of the grpc server to smoke-check")
	flag.StringVar(&amqpAddr, "mq-addr", "amqp://guest:guest@localhost:5672/", "Address of the mq server to smoke-check")
	flag.DurationVar(&sleepDuration, "scheduler-duration", time.Second, "Scheduler await timeout to smoke-check")
}

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTest Suite")
}
