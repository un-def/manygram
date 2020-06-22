project := 'manygram'
version := '0.1.0'

export GOOS := 'linux'
export GOARCH := 'amd64'

build_args := "-trimpath -ldflags='-s -w'"
build_output := project
archive_name := build_output + '-' + version + '-$GOOS-$GOARCH'

_list:
  @just --list

run +args='':
  go run ./cmd/{{project}} {{args}}

check-version:
  #!/bin/sh
  version=$(sed -n 's/.*manygramVersion = "\([^"]\+\)"/\1/p' internal/cli/cli.go)
  if [ -z "${version}" ]; then
    echo 'Cannot extract version from source code'
    exit 1
  fi
  if [ "${version}" != '{{version}}' ]; then
    echo "Version mismatch: ${version} != {{version}}"
    exit 1
  fi
  echo 'Version check passed'

build: check-version
  go build {{build_args}} -o {{build_output}} ./cmd/{{project}}

clean:
  rm -f ./{{build_output}}
  rm -f ./{{build_output}}*.tar.gz

test:
  go test -v ./...

archive: build
  tar -czf {{archive_name}}.tar.gz {{build_output}} LICENSE README.md CHANGELOG.md
