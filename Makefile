
run:
	@go fmt
	@go run .

build-container:
	docker build -t pkgs.doc.network/m3stuff/kiosk_controller:latest .

build-app:
	@go fmt
	@go build -o bin/kiosk_controller .

push:
	docker push pkgs.doc.network/m3stuff/kiosk_controller:latest

build: build-app build-container push
