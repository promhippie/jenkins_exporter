package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/dannietjoh/jenkins_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"

	_ "net/http/pprof"
)

var (
	showVersion   = flag.Bool("version", false, "Print version information")
	listenAddress = flag.String("web.listen-address", ":9118", "Address to listen on for web interface and telemetry")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path to expose metrics of the exporter")
	address       = flag.String("jenkins.address", "http://localhost:8080", "Address where to reach Jenkins")
	username      = flag.String("jenkins.username", "", "Username to authenticate on Jenkins")
	password      = flag.String("jenkins.password", "", "Password to authenticate on Jenkins")
)

// init registers the collector version.
func init() {
	prometheus.MustRegister(version.NewCollector("jenkins_exporter"))
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("jenkins_exporter"))
		os.Exit(0)
	}

	log.Infoln("Starting Jenkins exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	e := exporter.NewExporter(*address, *username, *password)

	prometheus.MustRegister(e)

	http.Handle(*metricsPath, promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Infof("Listening on %s", *listenAddress)

	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatal(err)
	}
}
