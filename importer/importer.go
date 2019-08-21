package importer

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/demetrio108/monit-grafana/datastore"
	"log"
	"regexp"
	"strconv"
)

var (
	loadAvgRe = regexp.MustCompile(`\[(?P<load5>[0-9.]+)\]\xA0\[(?P<load10>[0-9.]+)\]\xA0\[(?P<load15>[0-9.]+)\]`)
	cpuRe     = regexp.MustCompile(`(?P<us>[0-9.]+)%us,\xA0(?P<sy>[0-9.]+)%sy,\xA0(?P<wa>[0-9.]+)%wa`)
	memRe     = regexp.MustCompile(`(?P<perc>[0-9.]+)% \[(?P<val>[0-9.]+) (?P<units>\w+)\]`)
)

type Importer struct {
	url  string
	name string
	ds   *datastore.DataStore
}

func New(name string, url string, ds *datastore.DataStore) *Importer {
	imp := Importer{}
	imp.url = url
	imp.name = name
	imp.ds = ds

	return &imp
}

func (imp *Importer) Import() error {
	sysInfo := datastore.SystemInfo{}

	//TODO: Make dynamic slice resize on growth
	hostsInfo := make([]*datastore.HostInfo, 0, 1024)

	doc, err := htmlquery.LoadURL(imp.url)
	if err != nil {
		return err
	}

	tables := htmlquery.Find(doc, "//table[@id='header-row']")
	if len(tables) < 2 {
		return fmt.Errorf("Error parsing. Must be 2 monit tables")
	}

	for i, table := range tables {
		if i == 0 {
			// Parsing System Info
			for j, td := range htmlquery.Find(table, "//td") {
				switch j {
				case 0:
					sysInfo.System = htmlquery.InnerText(td)
				case 1:
					sysInfo.Status = htmlquery.InnerText(td)
				case 2:
					sysInfo.Load5, sysInfo.Load10, sysInfo.Load15, err = imp.parseLoad(htmlquery.InnerText(td))
					if err != nil {
						log.Print(imp.Name(), ": Importer: Error parsing load: ", err)
						log.Print(imp.Name(), ": Text is: ", htmlquery.InnerText(td))
					}
				case 3:
					sysInfo.CpuUs, sysInfo.CpuSy, sysInfo.CpuWa, err = imp.parseCpu(htmlquery.InnerText(td))
					if err != nil {
						log.Print(imp.Name(), ": Importer: Error parsing CPU: ", err)
						log.Print(imp.Name(), ": Text is: ", htmlquery.InnerText(td))
					}
				case 4:
					sysInfo.MemoryPerc, sysInfo.MemoryBytes, err = imp.parseMem(htmlquery.InnerText(td))
					if err != nil {
						log.Print(imp.Name(), ": Importer: Error parsing Memory: ", err)
						log.Print(imp.Name(), ": Text is: ", htmlquery.InnerText(td))
					}
				case 5:
					sysInfo.SwapPerc, sysInfo.SwapBytes, err = imp.parseMem(htmlquery.InnerText(td))
					if err != nil {
						log.Print(imp.Name(), ": Importer: Error parsing Swap: ", err)
						log.Print(imp.Name(), ": Text is: ", htmlquery.InnerText(td))
					}
				}
			}

			imp.ds.WriteSystemInfo(&sysInfo)
		} else {
			// Parsing Hosts
			for i, tr := range htmlquery.Find(table, "//tr") {
				if i > 0 {
					hostInfo := datastore.HostInfo{}
					for j, td := range htmlquery.Find(tr, "//td") {
						switch j {
						case 0:
							hostInfo.Host = htmlquery.InnerText(td)
						case 1:
							hostInfo.Status = htmlquery.InnerText(td)
						case 2:
							hostInfo.Protocols = htmlquery.InnerText(td)
						}
					}
					hostsInfo = append(hostsInfo, &hostInfo)
				}
			}

			imp.ds.WriteHostInfos(hostsInfo)
		}
	}

	return nil
}

func (imp *Importer) parseLoad(load string) (string, string, string, error) {
	var load5, load10, load15 string
	var err error

	for i, match := range loadAvgRe.FindStringSubmatch(load) {
		switch i {
		case 1:
			load5 = match
		case 2:
			load10 = match
		case 3:
			load15 = match
		}
	}

	return load5, load10, load15, err
}

func (imp *Importer) parseCpu(cpu string) (string, string, string, error) {
	var us, sy, wa string
	var err error

	for i, match := range cpuRe.FindStringSubmatch(cpu) {
		switch i {
		case 1:
			us = match
		case 2:
			sy = match
		case 3:
			wa = match
		}
	}

	return us, sy, wa, err
}

func (imp *Importer) parseMem(mem string) (string, string, error) {
	var perc, bytes, val, units string
	var err error

	for i, match := range memRe.FindStringSubmatch(mem) {
		switch i {
		case 1:
			perc = match
		case 2:
			val = match
		case 3:
			units = match
		}
	}

	val_float, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return perc, "err", err
	}

	switch units {
	case "KB":
		bytes = fmt.Sprintf("%.0f", val_float*1024)
	case "MB":
		bytes = fmt.Sprintf("%.0f", val_float*1024*1024)
	case "GB":
		bytes = fmt.Sprintf("%.0f", val_float*1024*1024*1024)
	}

	return perc, bytes, err
}

func (imp *Importer) Name() string {
	return imp.name
}
