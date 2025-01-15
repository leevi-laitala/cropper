APPNAME = cropper

$(APPNAME): main.go area.go video.go view.go const.go
	go build

lint: main.go area.go video.go view.go const.go
	golangci-lint run --enable-all --disable depguard,mnd,exhaustruct,err113,funlen

clean:
	rm $(APPNAME)
