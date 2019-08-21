package main

import (
	"flag"
	"github.com/demetrio108/monit-grafana/config"
	"github.com/demetrio108/monit-grafana/datastore"
	"github.com/demetrio108/monit-grafana/exporter"
	"github.com/demetrio108/monit-grafana/importer"
	"github.com/demetrio108/monit-grafana/router"
	"log"
	"time"
)

var (
	config_flag = flag.String("c", "/etc/monit-grafana.yml", "path to config")
	listen_flag = flag.String("l", ":8080", "host and port, format [host]:port")
)

func main() {
	flag.Parse()

	cfg, err := config.Load(*config_flag)
	if err != nil {
		log.Fatal(err)
	}

	if len(cfg.Instances) < 1 {
		log.Fatal("No monit instances to check")
	}

	exporters := make([]*exporter.Exporter, len(cfg.Instances))
	for i, instance := range cfg.Instances {
		ds := datastore.New()

		imp_ticker := time.NewTicker(time.Duration(instance.Interval) * time.Second)
		imp := importer.New(instance.Name, instance.URL, ds)

		importAndLog(imp)
		go func() {
			for _ = range imp_ticker.C {
				importAndLog(imp)
			}
		}()

		exp := exporter.New(instance.Name, ds)
		exporters[i] = exp
	}

	r := router.New(*listen_flag, exporters)

	log.Fatal(r.ListenAndServe())
}

func importAndLog(imp *importer.Importer) {
	log.Print(imp.Name(), ": Importing data from monit...")
	err := imp.Import()
	if err != nil {
		log.Print(imp.Name(), ": Error importing from monit: ", err)
	}
}
