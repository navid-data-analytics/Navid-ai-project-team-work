#!/bin/bash

# Check that all .go files added are correctly formatted
status=0
gofiles=$(find . -type f -name '*.go' -not -path "./vendor/*")
if ! [ -z "$gofiles" ]; then
	unformatted=$(docker-compose run test_go_common gofmt -l $gofiles)
	if ! [ -z "$unformatted" ]; then
		status=1
		# Some files are not gofmt'd. Print message and fail.

		echo >&2 "Go files must be formatted with gofmt before commit. Please run:"
		for fn in $unformatted; do
			echo >&2 "  docker-compose run test_go_common gofmt -w $fn"
		done
	fi
fi

exit $status
