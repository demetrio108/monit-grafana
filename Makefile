current_dir = $(shell pwd)

release:
	docker run --rm -v $(current_dir):/go/src/github.com/demetrio108/monit-grafana -w /go/src/github.com/demetrio108/monit-grafana -e CGO_ENABLED=0 golang:1.20-alpine sh -c "apk add --no-cache --update alpine-sdk git make && go mod tidy && go build -a -installsuffix cgo -buildvcs=false -ldflags '-s -w -extldflags \"-static\"' ."

rpm:
	rm -rf .build
	mkdir -p .build/usr/bin
	mkdir -p .build/etc
	cp monit-grafana .build/usr/bin/
	cp monit-grafana.yml .build/etc/
	docker run -it --rm -v $(current_dir):/monit-grafana -w /monit-grafana lrdevops/fpm fpm -s dir -t rpm -C .build --name monit-grafana --version 0.0.1 --iteration 1

deb:
	rm -rf .build
	mkdir -p .build/usr/bin
	mkdir -p .build/etc
	cp monit-grafana .build/usr/bin/
	cp monit-grafana.yml .build/etc/
	docker run -it --rm -v $(current_dir):/monit-grafana -w /monit-grafana lrdevops/fpm fpm -s dir -t deb -C .build --name monit-grafana --version 0.0.1 --iteration 1
