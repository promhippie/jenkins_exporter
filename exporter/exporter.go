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

	// queueItemStartTime defines a map to collect the queue items.
	queueItemStartTime = map[string]prometheus.Gauge{}
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

	for _, metric := range queueItemStartTime {
		ch <- metric
	}
}

// scrape just starts the scraping loop.
func (e *Exporter) scrape() error {
	log.Debug("start scrape loop")

	var (
		metrics = &Collector{}
	)

	if err := metrics.Fetch(e.address, e.username, e.password); err != nil {
		log.Debugf("%s", err)
		return fmt.Errorf("failed to fetch metrics")
	}

	for _, jobItem := range metrics.Jobs {
		log.Debugf("processing jobItem %s", jobItem.Name)

		if jobItem.Color != "" {
			if _, ok := jobColor[jobItem.Name]; ok == false {
				jobColor[jobItem.Name] = prometheus.NewGauge(
					prometheus.GaugeOpts{
						Namespace: namespace,
						Name:      "job_color",
						Help:      "Color code of the Jenkins job",
						ConstLabels: prometheus.Labels{
							"name": jobItem.Name,
						},
					},
				)
			}

			color := colorToGauge(jobItem.Color)
			log.Debugf("setting job_color to %f for %s", color, jobItem.Name)
			jobColor[jobItem.Name].Set(colorToGauge(jobItem.Color))
		}
	}

	for _, queueItem := range metrics.Queue {
		log.Debugf("processing queueItem %s", queueItem.Task.Name)

		queueItemStartTime[queueItem.Task.Name] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "queue_item_start_time_seconds",
				Help:      "Time the item entered the Jenkins queue since unix epoch in seconds.",
				ConstLabels: prometheus.Labels{
					"name": queueItem.Task.Name,
				},
			},
		)

		log.Debugf("setting queue_item_start_time_seconds to %f for %s", queueItem.InQueueSince, queueItem.Task.Name)
		queueItemStartTime[queueItem.Task.Name].Set(queueItem.InQueueSince)
	}

	isUp.Set(1)
	return nil
}
