build:
	./scripts/build.sh

install:
	./scripts/install.sh

uninstall:
	./scripts/uninstall.sh

run:
	cd examples && go run ../main.go
