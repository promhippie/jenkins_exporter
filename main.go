package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.com/webhippie/jenkins_exporter/exporter"

	_ "net/http/pprof"
)

var (
	// showVersion is a flag to display the current version.
	showVersion = flag.Bool("version", false, "Print version information")

	// listenAddress defines the local address binding for the server.
	listenAddress = flag.String("web.listen-address", ":9103", "Address to listen on for web interface and telemetry")

	// metricsPath defines the path to access the metrics.
	metricsPath = flag.String("web.telemetry-path", "/metrics", "Path to expose metrics of the exporter")

	// address defines the URL to access the Jenkins instance.
	address = flag.String("jenkins.address", "", "Address where to reach Jenkins")

	// username defines the username for the Jenkins authentication.
	username = flag.String("jenkins.username", "", "Username to authenticate on Jenkins")

	// password defines the password for the Jenkins authentication.
	password = flag.String("jenkins.password", "", "Password to authenticate on Jenkins")
)

// init registers the collector version.
func init() {
	prometheus.MustRegister(version.NewCollector("jenkins_exporter"))
}

// main simply initializes this tool.
func main() {
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("jenkins_exporter"))
		os.Exit(0)
	}

	if *address == "" {
		fmt.Fprintln(os.Stderr, "Please provide a address for Jenkins")
		os.Exit(1)
	}

	log.Infoln("Starting Jenkins exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	e := exporter.NewExporter(*address, *username, *password)

	prometheus.MustRegister(e)
	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.Unregister(prometheus.NewProcessCollector(os.Getpid(), ""))

	http.Handle(*metricsPath, promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Infof("Listening on %s", *listenAddress)

	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatal(err)
	}
}
