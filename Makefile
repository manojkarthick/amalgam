cleanup:
	rm -rf extract* fat* *.tar.gz *.zip

build: cleanup
	go build -o amalgam

dry-run: build
	./amalgam --owner manojkarthick --repo expenses --tag latest --amd64 "darwin_amd64" --arm64 "darwin_arm64" --compressed --overwrite

dry-run-rust: build
	./amalgam --owner manojkarthick --repo jreleaser-poc --tag latest --amd64 "x86_64-apple-darwin" --arm64 "aarch64-apple-darwin" --compressed --overwrite