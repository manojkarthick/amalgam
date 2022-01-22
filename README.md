# amalgam

* CLI Tool to build macOS universal binaries and upload to a Github Release
* You can download release binaries for linux/macOS/windows [here](https://github.com/manojkarthick/amalgam/releases).

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
