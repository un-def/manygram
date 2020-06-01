project := 'manygram'

build_args := "-trimpath"

_list:
  @just --list

run +args='':
  go run ./cmd/{{project}} {{args}}

build:
  go build {{build_args}} ./cmd/{{project}}

clean:
  rm -f ./{{project}}

test:
  go test -v ./...
