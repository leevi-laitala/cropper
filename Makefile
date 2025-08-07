APPNAME=cropper

build: main.go area.go video.go view.go global.go seeker.go
	go build

lint: main.go area.go video.go view.go global.go seeker.go
	golangci-lint run --config ./golangci-config.yaml

clean:
	rm $(APPNAME)
