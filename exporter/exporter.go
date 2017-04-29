package exporter

import (
	"fmt"
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

	// jobColor defines a map to collect the build color codes.
	jobColor = map[string]prometheus.Gauge{}
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

	for _, metric := range jobColor {
		ch <- metric.Desc()
	}
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

	for _, metric := range jobColor {
		ch <- metric
	}
}

// scrape just starts the scraping loop.
func (e *Exporter) scrape() error {
	log.Debug("start scrape loop")

	var (
		root = &Root{}
	)

	if err := root.Fetch(e.address, e.username, e.password); err != nil {
		log.Debugf("%s", err)
		return fmt.Errorf("failed to fetch root data")
	}

	for _, job := range root.Jobs {
		log.Debugf("processing %s job", job.Name)

		if job.Color != "" {
			if _, ok := jobColor[job.Key()]; ok == false {
				jobColor[job.Key()] = prometheus.NewGauge(
					prometheus.GaugeOpts{
						Namespace: namespace,
						Name:      "job_color",
						Help:      "Color code of the Jenkins job",
						ConstLabels: prometheus.Labels{
							"name": job.Name,
						},
					},
				)
			}

			color := colorToGauge(job.Color)
			log.Debugf("setting color to %f for %s", color, job.Name)

			jobColor[job.Key()].Set(colorToGauge(job.Color))
		}
	}

	isUp.Set(1)
	return nil
}
