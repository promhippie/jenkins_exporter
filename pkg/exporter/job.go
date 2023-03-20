package exporter

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/promhippie/jenkins_exporter/pkg/config"
	"github.com/promhippie/jenkins_exporter/pkg/internal/jenkins"
)

// JobCollector collects metrics about the servers.
type JobCollector struct {
	client   *jenkins.Client
	logger   log.Logger
	failures *prometheus.CounterVec
	duration *prometheus.HistogramVec
	config   config.Target

	Disabled              *prometheus.Desc
	Buildable             *prometheus.Desc
	Color                 *prometheus.Desc
	LastBuild             *prometheus.Desc
	LastCompletedBuild    *prometheus.Desc
	LastFailedBuild       *prometheus.Desc
	LastStableBuild       *prometheus.Desc
	LastSuccessfulBuild   *prometheus.Desc
	LastUnstableBuild     *prometheus.Desc
	LastUnsuccessfulBuild *prometheus.Desc
	NextBuild             *prometheus.Desc
}

// NewJobCollector returns a new JobCollector.
func NewJobCollector(logger log.Logger, client *jenkins.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *JobCollector {
	if failures != nil {
		failures.WithLabelValues("job").Add(0)
	}

	labels := []string{"name", "path", "class"}
	return &JobCollector{
		client:   client,
		logger:   log.With(logger, "collector", "job"),
		failures: failures,
		duration: duration,
		config:   cfg,

		Disabled: prometheus.NewDesc(
			"jenkins_job_disabled",
			"1 if the job is disabled, 0 otherwise",
			labels,
			nil,
		),
		Buildable: prometheus.NewDesc(
			"jenkins_job_buildable",
			"1 if the sjob is buildable, 0 otherwise",
			labels,
			nil,
		),
		Color: prometheus.NewDesc(
			"jenkins_job_color",
			"Color code of the jenkins job",
			labels,
			nil,
		),
		LastBuild: prometheus.NewDesc(
			"jenkins_job_last_build",
			"Builder number for last build",
			labels,
			nil,
		),
		LastCompletedBuild: prometheus.NewDesc(
			"jenkins_job_last_completed_build",
			"Builder number for last completed build",
			labels,
			nil,
		),
		LastFailedBuild: prometheus.NewDesc(
			"jenkins_job_last_failed_build",
			"Builder number for last failed build",
			labels,
			nil,
		),
		LastStableBuild: prometheus.NewDesc(
			"jenkins_job_last_stable_build",
			"Builder number for last stable build",
			labels,
			nil,
		),
		LastSuccessfulBuild: prometheus.NewDesc(
			"jenkins_job_last_successful_build",
			"Builder number for last successful build",
			labels,
			nil,
		),
		LastUnstableBuild: prometheus.NewDesc(
			"jenkins_job_last_unstable_build",
			"Builder number for last unstable build",
			labels,
			nil,
		),
		LastUnsuccessfulBuild: prometheus.NewDesc(
			"jenkins_job_last_unsuccessful_build",
			"Builder number for last unsuccessful build",
			labels,
			nil,
		),
		NextBuild: prometheus.NewDesc(
			"jenkins_job_next_build_number",
			"Next build number for the job",
			labels,
			nil,
		),
	}
}

// Metrics simply returns the list metric descriptors for generating a documentation.
func (c *JobCollector) Metrics() []*prometheus.Desc {
	return []*prometheus.Desc{
		c.Disabled,
		c.Buildable,
		c.Color,
		c.LastBuild,
		c.LastCompletedBuild,
		c.LastFailedBuild,
		c.LastStableBuild,
		c.LastSuccessfulBuild,
		c.LastUnstableBuild,
		c.LastUnsuccessfulBuild,
		c.NextBuild,
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *JobCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Disabled
	ch <- c.Buildable
	ch <- c.Color
	ch <- c.LastBuild
	ch <- c.LastCompletedBuild
	ch <- c.LastFailedBuild
	ch <- c.LastStableBuild
	ch <- c.LastSuccessfulBuild
	ch <- c.LastUnstableBuild
	ch <- c.LastUnsuccessfulBuild
	ch <- c.NextBuild
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *JobCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	now := time.Now()
	jobs, err := c.client.Job.All(ctx)
	c.duration.WithLabelValues("job").Observe(time.Since(now).Seconds())

	if err != nil {
		level.Error(c.logger).Log(
			"msg", "Failed to fetch jobs",
			"err", err,
		)

		c.failures.WithLabelValues("job").Inc()
		return
	}

	level.Debug(c.logger).Log(
		"msg", "Fetched jobs",
		"count", len(jobs),
	)

	for _, job := range jobs {
		var (
			disabled  float64
			buildable float64
		)

		labels := []string{
			job.Name,
			job.Path,
			job.Class,
		}

		if job.Disabled {
			disabled = 1.0
		}

		ch <- prometheus.MustNewConstMetric(
			c.Disabled,
			prometheus.GaugeValue,
			disabled,
			labels...,
		)

		if job.Buildable {
			buildable = 1.0
		}

		ch <- prometheus.MustNewConstMetric(
			c.Buildable,
			prometheus.GaugeValue,
			buildable,
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Color,
			prometheus.GaugeValue,
			colorToGauge(job.Color),
			labels...,
		)

		if job.LastBuild != nil {
			ch <- prometheus.MustNewConstMetric(
				c.LastBuild,
				prometheus.GaugeValue,
				float64(*job.LastBuild),
				labels...,
			)
		}

		if job.LastCompletedBuild != nil {
			ch <- prometheus.MustNewConstMetric(
				c.LastCompletedBuild,
				prometheus.GaugeValue,
				float64(*job.LastCompletedBuild),
				labels...,
			)
		}

		if job.LastFailedBuild != nil {
			ch <- prometheus.MustNewConstMetric(
				c.LastFailedBuild,
				prometheus.GaugeValue,
				float64(*job.LastFailedBuild),
				labels...,
			)
		}

		if job.LastStableBuild != nil {
			ch <- prometheus.MustNewConstMetric(
				c.LastStableBuild,
				prometheus.GaugeValue,
				float64(*job.LastStableBuild),
				labels...,
			)
		}

		if job.LastSuccessfulBuild != nil {
			ch <- prometheus.MustNewConstMetric(
				c.LastSuccessfulBuild,
				prometheus.GaugeValue,
				float64(*job.LastSuccessfulBuild),
				labels...,
			)
		}

		if job.LastUnstableBuild != nil {
			ch <- prometheus.MustNewConstMetric(
				c.LastUnstableBuild,
				prometheus.GaugeValue,
				float64(*job.LastUnstableBuild),
				labels...,
			)
		}

		if job.LastUnsuccessfulBuild != nil {
			ch <- prometheus.MustNewConstMetric(
				c.LastUnsuccessfulBuild,
				prometheus.GaugeValue,
				float64(*job.LastUnsuccessfulBuild),
				labels...,
			)
		}

		ch <- prometheus.MustNewConstMetric(
			c.NextBuild,
			prometheus.GaugeValue,
			float64(job.NextBuildNumber),
			labels...,
		)
	}
}

func colorToGauge(color string) float64 {
	switch color {
	case "blue":
		return 1.0
	case "blue_anime":
		return 1.5
	case "red":
		return 2.0
	case "red_anime":
		return 2.5
	case "yellow":
		return 3.0
	case "yellow_anime":
		return 3.5
	case "notbuilt":
		return 4.0
	case "notbuilt_anime":
		return 4.5
	case "disabled":
		return 5.0
	case "disabled_anime":
		return 5.5
	case "aborted":
		return 6.0
	case "aborted_anime":
		return 6.5
	case "grey":
		return 7.0
	case "grey_anime":
		return 7.5
	}

	return 0.0
}
