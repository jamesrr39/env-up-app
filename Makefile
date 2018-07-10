.PHONY: clean
clean:
	rm -rf build

.PHONY: run_dev
run_dev:
	go run cmd/env-up-main.go testdata/simple_config/env-up-conf.yaml

.PHONY: build_prod_x86_64
build_prod_x86_64: clean
	go run vendor/github.com/rakyll/statik/statik.go -src=frontend-web -dest=build/client
	mkdir -p build/bin/x86_64
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags "purego prod" -o build/bin/x86_64/env-up-app cmd/env-up-main.go
