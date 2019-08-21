package exporter

import (
	"encoding/json"
	"github.com/demetrio108/monit-grafana/datastore"
	"io/ioutil"
	"log"
	"net/http"
)

type Exporter struct {
	name string
	ds   *datastore.DataStore
}

type GrafanaRequest struct {
	Targets []struct {
		RefID  string `json:"refId"`
		Target string `json:"target"`
		Type   string `json:"type"`
	} `json:"targets"`
}

type GrafanaResponseColumn struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type GrafanaResponse struct {
	Columns []GrafanaResponseColumn `json:"columns"`
	Rows    [][]string              `json:"rows"`
	Type    string                  `json:"type"`
}

func (exp *Exporter) RootHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func (exp *Exporter) SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`["system_info", "hosts_info"]`))
}

func (exp *Exporter) QueryHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(exp.Name(), ": Error reading request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := GrafanaRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Print(exp.Name(), ": Error unmarshalling request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(req.Targets) != 1 {
		log.Print(exp.Name(), ": Error: there should be one target only")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	table := req.Targets[0].Target

	var resp []byte
	grafanaResp := make([]GrafanaResponse, 1)
	grafanaResp[0] = GrafanaResponse{}
	grafanaResp[0].Type = "table"

	switch table {
	case "system_info":
		log.Print(exp.Name(), ": Serving system_info...")
		grafanaResp[0].Columns = make([]GrafanaResponseColumn, 12)

		grafanaResp[0].Columns[0] = GrafanaResponseColumn{"System", "string"}
		grafanaResp[0].Columns[1] = GrafanaResponseColumn{"Status", "string"}
		grafanaResp[0].Columns[2] = GrafanaResponseColumn{"Load 5", "string"}
		grafanaResp[0].Columns[3] = GrafanaResponseColumn{"Load 10", "string"}
		grafanaResp[0].Columns[4] = GrafanaResponseColumn{"Load 15", "string"}
		grafanaResp[0].Columns[5] = GrafanaResponseColumn{"CPU % us", "string"}
		grafanaResp[0].Columns[6] = GrafanaResponseColumn{"CPU % sy", "string"}
		grafanaResp[0].Columns[7] = GrafanaResponseColumn{"CPU % wa", "string"}
		grafanaResp[0].Columns[8] = GrafanaResponseColumn{"Memory %", "string"}
		grafanaResp[0].Columns[9] = GrafanaResponseColumn{"Memory bytes", "string"}
		grafanaResp[0].Columns[10] = GrafanaResponseColumn{"Swap %", "string"}
		grafanaResp[0].Columns[11] = GrafanaResponseColumn{"Swap bytes", "string"}

		sysInfo := exp.ds.ReadSystemInfo()

		var status string

		if sysInfo.Status == "OK" {
			status = "1"
		} else {
			status = "0"
		}

		grafanaResp[0].Rows = [][]string{
			[]string{
				sysInfo.System,
				status,
				sysInfo.Load5,
				sysInfo.Load10,
				sysInfo.Load15,
				sysInfo.CpuUs,
				sysInfo.CpuSy,
				sysInfo.CpuWa,
				sysInfo.MemoryPerc,
				sysInfo.MemoryBytes,
				sysInfo.SwapPerc,
				sysInfo.SwapBytes,
			},
		}

	case "hosts_info":
		log.Print(exp.Name(), ": Serving hosts_info...")
		grafanaResp[0].Columns = make([]GrafanaResponseColumn, 3)

		grafanaResp[0].Columns[0] = GrafanaResponseColumn{"Host", "string"}
		grafanaResp[0].Columns[1] = GrafanaResponseColumn{"Status", "string"}
		grafanaResp[0].Columns[2] = GrafanaResponseColumn{"Protocol", "string"}

		hostsInfo := exp.ds.ReadHostInfos()
		rows := make([][]string, len(hostsInfo))

		for i, _ := range rows {
			var status string

			if hostsInfo[i].Status == "OK" {
				status = "1"
			} else {
				status = "0"
			}

			rows[i] = []string{
				hostsInfo[i].Host,
				status,
				hostsInfo[i].Protocols,
			}
		}

		grafanaResp[0].Rows = rows
	}

	resp, err = json.Marshal(grafanaResp)
	if err != nil {
		log.Print(exp.Name(), ": Error marshalling response: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.Write(resp)
}

func New(name string, ds *datastore.DataStore) *Exporter {
	exp := Exporter{}
	exp.name = name
	exp.ds = ds

	return &exp
}

func (e *Exporter) Name() string {
	return e.name
}
