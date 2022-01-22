cleanup:
	rm -rf extract* fat* *.tar.gz

build: cleanup
	go build -o amalgam

dry-run: build
	./amalgam --owner manojkarthick --repo expenses --tag latest --amd64 "darwin_amd64" --arm64 "darwin_arm64" --compressed --overwrite --binary-path expenses
