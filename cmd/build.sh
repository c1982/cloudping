${HOME}/go/bin/gox -osarch="linux/amd64 windows/amd64 darwin/amd64" -output="cping_{{.OS}}_{{.Arch}}" -ldflags="-s -w"