# monit-grafana

Exporting monit dashboard data in Grafanas's SimpleJSON compatible format.

*This is an alpha software*

## Motivation

In my work I have to deal with many docker-compose based installations. Looking for simple monitoring solutions I found 
[monit](https://mmonit.com/monit/), which provides simple and clean syntax for checks. Adding container with monit to 
docker-compose file allows to use monit's remote host checks to monitor containers health.

However, there is no free way to gather data from multiple monit instances, so I decided to make this project. It allows 
to export data from multiple monit dashboards to Grafana via SimpleJSON plugin.

## Quick start

For checking out how it works please refer to this [deployment example](https://github.com/monit-grafana-example).

## Running

You can build monit-grafana from source, or run as docker containerized version:
```
docker run -it demetrio108/monit-grafana
```

Run `monit-grafana` with following options:
* `-c <path-to-config>` - path to configuration file (default: `/etc/monit-grafana.yml`)
* `-l <[address]:port>` - bind to this address and port (default `:8080`)

### Configuration file

In configuration file you have to provide list of monit instances to export from in YAML format:
 ```
 ---
 - name: monit_test
   url: http://admin:monit@monit-host:2812/
   interval: 30

 - name: monit_test_2
   url: http://admin:monit@monit-host-2:2812/
   interval: 30
 ```

 `name` is user as root context fro Grafana datasource.

### Grafana setup

First you have to install SimpleJSON plugin into Grafana:
```
grafana-cli plugins install grafana-simple-json-datasource
```

Then you can add datasources from monit-grafana. For provided above config there wil be 2 datasources:
* `http://<monit-grafana>:8080/monit_test/`
* `http://<monit-grafana>:8080/monit_test_2/`

Two queries of type table are supported:
* `system_info` - System info from monit dashboards
* `hosts_info` - Remote host checks from monit dashboard

There is an example dashboard `monit-grafana-dashboard.json` provided in this repository.
