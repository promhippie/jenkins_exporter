package exporter

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	// namespace defines the Prometheus namespace for this exporter.
	namespace = "jenkins"
)

var (
	// isUp defines if the API response can get processed.
	isUp = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Check if Jenkins response can be processed",
		},
	)
)

// init just defines the initial state of the exports.
func init() {
	isUp.Set(0)
}

// NewExporter gives you a new exporter instance.
func NewExporter(address, username, password string) *Exporter {
	return &Exporter{
		address:  address,
		username: username,
		password: password,
	}
}

// Exporter combines the metric collector and descritions.
type Exporter struct {
	address  string
	username string
	password string
	mutex    sync.RWMutex
}

// Describe defines the metric descriptions for Prometheus.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- isUp.Desc()
}

// Collect delivers the metrics to Prometheus.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if err := e.scrape(); err != nil {
		log.Error(err)

		isUp.Set(0)
		ch <- isUp

		return
	}

	ch <- isUp
}

// scrape just starts the scraping loop.
func (e *Exporter) scrape() error {
	log.Debug("start scrape loop")

	isUp.Set(1)
	return nil
}
