# amalgam

> [!WARNING]
> **amalgam is deprecated and the repository is archived.** It is no longer maintained.
> Existing release binaries, Docker images, and the Homebrew formula remain available.
>
> Create universal binaries as part of your build instead:
>
> * **Go:** use GoReleaser's built-in [`universal_binaries`](https://goreleaser.com/customization/builds/universalbinaries/). It runs on Linux runners.
>   ```yaml
>   # .goreleaser.yaml
>   universal_binaries:
>     - replace: true
>   ```
> * **On macOS (any language):** `lipo` ships with Xcode.
>   ```sh
>   lipo -create -output myapp myapp-arm64 myapp-amd64
>   ```
> * **On Linux (any language):** use [konoui/lipo](https://github.com/konoui/lipo), a `lipo`-compatible tool written in Go, or `llvm-lipo`.
>
> Sign and notarize the merged binary after you create it.

* CLI Tool to build macOS universal binaries and upload to a Github Release

## Installation

### Release binaries

You can download release binaries for linux/macOS/windows [here](https://github.com/manojkarthick/amalgam/releases).

### Docker

Pre-built docker images are built and stored in [Docker Hub](https://hub.docker.com/r/manojkarthick/amalgam) and the [Github container registry](https://github.com/manojkarthick/amalgam/pkgs/container/amalgam).

### Homebrew

```shell
brew tap manojkarthick/tap
brew install amalgam
```

## Usage

```shell
❯ amalgam --help
NAME:
   amalgam - Create macOS Universal binaries from Github releases

USAGE:
   amalgam [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --owner value        Github repo owner username
   --repo value         Github repository name
   --tag value          Github repository name (default: "latest")
   --amd64 value        Substring for the amd64 binary
   --arm64 value        Substring for the arm64 binary
   --compressed         Do the releases use compressed archives (default: false)
   --binary-path value  Path to the binary inside the archive
   --overwrite          Delete pre-existing universal asset? (default: false)
   --help, -h           show help (default: false)
```

## Troubleshooting

On macOS, if you encounter the error `“amalgam” cannot be opened because the developer cannot be verified`, run the following:

```shell
xattr -d com.apple.quarantine ./amalgam
```

## Credits

Thanks to Keith Randall for the [makefat](https://github.com/randall77/makefat) library this is based on.
