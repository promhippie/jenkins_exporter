package exporter

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/promhippie/jenkins_exporter/pkg/config"
	"github.com/promhippie/jenkins_exporter/pkg/internal/jenkins"
)

// JobCollector collects metrics about the servers.
type JobCollector struct {
	client   *jenkins.Client
	logger   *slog.Logger
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
	Duration              *prometheus.Desc
	StartTime             *prometheus.Desc
	EndTime               *prometheus.Desc
}

// NewJobCollector returns a new JobCollector.
func NewJobCollector(logger *slog.Logger, client *jenkins.Client, failures *prometheus.CounterVec, duration *prometheus.HistogramVec, cfg config.Target) *JobCollector {
	if failures != nil {
		failures.WithLabelValues("job").Add(0)
	}

	labels := []string{"name", "path", "class"}
	return &JobCollector{
		client:   client,
		logger:   logger.With("collector", "job"),
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
			"1 if the job is buildable, 0 otherwise",
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
		Duration: prometheus.NewDesc(
			"jenkins_job_duration",
			"Duration of last build in ms",
			labels,
			nil,
		),
		StartTime: prometheus.NewDesc(
			"jenkins_job_start_time",
			"Start time of last build as unix timestamp",
			labels,
			nil,
		),
		EndTime: prometheus.NewDesc(
			"jenkins_job_end_time",
			"Start time of last build as unix timestamp",
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
		c.Duration,
		c.StartTime,
		c.EndTime,
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
	ch <- c.Duration
	ch <- c.StartTime
	ch <- c.EndTime
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *JobCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	now := time.Now()
	jobs, err := c.client.Job.All(ctx)
	c.duration.WithLabelValues("job").Observe(time.Since(now).Seconds())

	if err != nil {
		c.logger.Error("Failed to fetch jobs",
			"err", err,
		)

		c.failures.WithLabelValues("job").Inc()
		return
	}

	c.logger.Debug("Fetched jobs",
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
				float64(job.LastBuild.Number),
				labels...,
			)

			build, err := c.client.Job.Build(ctx, job.LastBuild)

			if err != nil {
				c.logger.Error("Failed to fetch last build",
					"job", job.Path,
					"err", err,
				)

				c.failures.WithLabelValues("job").Inc()
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.Duration,
					prometheus.GaugeValue,
					float64(build.Duration),
					labels...,
				)

				ch <- prometheus.MustNewConstMetric(
					c.StartTime,
					prometheus.GaugeValue,
					float64(build.Timestamp),
					labels...,
				)

				ch <- prometheus.MustNewConstMetric(
					c.EndTime,
					prometheus.GaugeValue,
					float64(build.Timestamp+build.Duration),
					labels...,
				)
			}
		}

		if job.LastCompletedBuild != nil {
			ch <- prometheus.MustNewConstMetric(
				c.LastCompletedBuild,
				prometheus.GaugeValue,
				float64(job.LastCompletedBuild.Number),
				labels...,
			)
		}

		if job.LastFailedBuild != nil {
			ch <- prometheus.MustNewConstMetric(
				c.LastFailedBuild,
				prometheus.GaugeValue,
				float64(job.LastFailedBuild.Number),
				labels...,
			)
		}

		if job.LastStableBuild != nil {
			ch <- prometheus.MustNewConstMetric(
				c.LastStableBuild,
				prometheus.GaugeValue,
				float64(job.LastStableBuild.Number),
				labels...,
			)
		}

		if job.LastSuccessfulBuild != nil {
			ch <- prometheus.MustNewConstMetric(
				c.LastSuccessfulBuild,
				prometheus.GaugeValue,
				float64(job.LastSuccessfulBuild.Number),
				labels...,
			)
		}

		if job.LastUnstableBuild != nil {
			ch <- prometheus.MustNewConstMetric(
				c.LastUnstableBuild,
				prometheus.GaugeValue,
				float64(job.LastUnstableBuild.Number),
				labels...,
			)
		}

		if job.LastUnsuccessfulBuild != nil {
			ch <- prometheus.MustNewConstMetric(
				c.LastUnsuccessfulBuild,
				prometheus.GaugeValue,
				float64(job.LastUnsuccessfulBuild.Number),
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
