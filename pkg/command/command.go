package command

import (
	"fmt"
	"os"
	"time"

	"github.com/promhippie/jenkins_exporter/pkg/action"
	"github.com/promhippie/jenkins_exporter/pkg/config"
	"github.com/promhippie/jenkins_exporter/pkg/version"
	"github.com/urfave/cli/v2"
)

// Run parses the command line arguments and executes the program.
func Run() error {
	cfg := config.Load()

	app := &cli.App{
		Name:    "jenkins_exporter",
		Version: version.String,
		Usage:   "Jenkins Exporter",
		Authors: []*cli.Author{
			{
				Name:  "Thomas Boerger",
				Email: "thomas@webhippie.de",
			},
		},
		Flags: RootFlags(cfg),
		Commands: []*cli.Command{
			Health(cfg),
		},
		Action: func(_ *cli.Context) error {
			logger := setupLogger(cfg)

			if cfg.Target.Address == "" {
				logger.Error("Missing required jenkins.url")
				return fmt.Errorf("missing required jenkins.url")
			}

			return action.Server(cfg, logger)
		},
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Show the help, so what you see now",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print the current version of that tool",
	}

	return app.Run(os.Args)
}

// RootFlags defines the available root flags.
func RootFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log.level",
			Value:       "info",
			Usage:       "Only log messages with given severity",
			EnvVars:     []string{"JENKINS_EXPORTER_LOG_LEVEL"},
			Destination: &cfg.Logs.Level,
		},
		&cli.BoolFlag{
			Name:        "log.pretty",
			Value:       false,
			Usage:       "Enable pretty messages for logging",
			EnvVars:     []string{"JENKINS_EXPORTER_LOG_PRETTY"},
			Destination: &cfg.Logs.Pretty,
		},
		&cli.StringFlag{
			Name:        "web.address",
			Value:       "0.0.0.0:9506",
			Usage:       "Address to bind the metrics server",
			EnvVars:     []string{"JENKINS_EXPORTER_WEB_ADDRESS"},
			Destination: &cfg.Server.Addr,
		},
		&cli.StringFlag{
			Name:        "web.path",
			Value:       "/metrics",
			Usage:       "Path to bind the metrics server",
			EnvVars:     []string{"JENKINS_EXPORTER_WEB_PATH"},
			Destination: &cfg.Server.Path,
		},
		&cli.BoolFlag{
			Name:        "web.debug",
			Value:       false,
			Usage:       "Enable pprof debugging for server",
			EnvVars:     []string{"JENKINS_EXPORTER_WEB_PPROF"},
			Destination: &cfg.Server.Pprof,
		},
		&cli.DurationFlag{
			Name:        "web.timeout",
			Value:       10 * time.Second,
			Usage:       "Server metrics endpoint timeout",
			EnvVars:     []string{"JENKINS_EXPORTER_WEB_TIMEOUT"},
			Destination: &cfg.Server.Timeout,
		},
		&cli.StringFlag{
			Name:        "web.config",
			Value:       "",
			Usage:       "Path to web-config file",
			EnvVars:     []string{"JENKINS_EXPORTER_WEB_CONFIG"},
			Destination: &cfg.Server.Web,
		},
		&cli.DurationFlag{
			Name:        "request.timeout",
			Value:       5 * time.Second,
			Usage:       "Timeout requesting Jenkins API",
			EnvVars:     []string{"JENKINS_EXPORTER_REQUEST_TIMEOUT"},
			Destination: &cfg.Target.Timeout,
		},
		&cli.StringFlag{
			Name:        "jenkins.url",
			Value:       "",
			Usage:       "URL to access the Jenkins to scrape",
			EnvVars:     []string{"JENKINS_EXPORTER_URL"},
			Destination: &cfg.Target.Address,
		},
		&cli.StringFlag{
			Name:        "jenkins.username",
			Value:       "",
			Usage:       "Username for the Jenkins authentication",
			EnvVars:     []string{"JENKINS_EXPORTER_USERNAME"},
			Destination: &cfg.Target.Username,
		},
		&cli.StringFlag{
			Name:        "jenkins.password",
			Value:       "",
			Usage:       "Password for the Jenkins authentication",
			EnvVars:     []string{"JENKINS_EXPORTER_PASSWORD"},
			Destination: &cfg.Target.Password,
		},
		&cli.BoolFlag{
			Name:        "collector.jobs",
			Value:       true,
			Usage:       "Enable collector for jobs",
			EnvVars:     []string{"JENKINS_EXPORTER_COLLECTOR_JOBS"},
			Destination: &cfg.Collector.Jobs,
		},
	}
}
