.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -rf build

.PHONY: run_dev
run_dev:
	go run cmd/env-up-main.go testdata/simple_config/env-up-conf.yaml

.PHONY: bundle_static_assets
bundle_static_assets:
	go run vendor/github.com/rakyll/statik/statik.go -src=frontend-web -dest=build/client

.PHONY: build_prod_linux_x86_64
build_prod_linux_x86_64: clean bundle_static_assets
	mkdir -p build/bin/linux_x86_64
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags "purego prod" -o build/bin/linux_x86_64/env-up-app cmd/env-up-main.go

.PHONY: build_prod_macos
build_prod_macos: clean bundle_static_assets
	mkdir -p build/bin/macos
	env GOOS=macos GOARCH=amd64 CGO_ENABLED=0 go build -tags "purego prod" -o build/bin/macos/env-up-app cmd/env-up-main.go
