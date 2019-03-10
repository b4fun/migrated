.PHONY: build
build: build-bindata
	go build -o migrated-cli main.go

.PHONY: build-bindata
build-bindata:
	cd migration && go-bindata -modtime 1500000000 \
	    -pkg=migration \
	    -ignore=.go \
	    .
