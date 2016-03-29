createandupdate: *.go
	go build

version=0.1.0

release: ./createandupdate
	rm -f pkg/*_amd64 pkg/*.exe
	ghr -u flozano v$(version) pkg

compress:
	gzip -fk pkg/*_amd64
	zip pkg/createandupdate_windows_amd64.zip pkg/createandupdate_windows_amd64.exe

build:
	gox -os "linux darwin windows" -arch "amd64"  -output "pkg/createandupdate_{{.OS}}_{{.Arch}}" .

clean:
	rm -f createandupdate

distclean: clean
	rm -rf pkg
