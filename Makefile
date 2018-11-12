init:
	rm -rf vendor
	rm -f Gopkg.lock
	rm -f Gopkg.toml
	dep init -v

u:
	dep ensure -v -update
