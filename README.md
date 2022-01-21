# amalgam
Tool to build macOS universal binaries and upload to a Github Release

## TODO

[ ] Integrate makefat library
[ ] Authenticate to Github
[ ] Support specifying file regexes for finding amd64 and arm64 binaries
[ ] Support specifying specific release tag
[ ] Fallback to latest github release tag
[ ] Support providing a name for the universal binary (or use name from amd64-regex)
[ ] Specify whether the binary is using a compressed archive
[ ] Specify path to the binary with the extracted archive
[ ] Make it available as a Github Actions Workflow

Flags:
* --owner
* --repo
* --tag
* --amd64-regex
* --arm64-regex
* --uses-archive
* --binary-path

Environment variables
* GITHUB_TOKEN
