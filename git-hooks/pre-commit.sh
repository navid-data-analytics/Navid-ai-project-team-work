#!/bin/bash

#Check that all .py files added are correctly formatted
status=0
pyfiles=$(git diff --cached --name-only --diff-filter=ACM | grep '\.py$')
if ! [ -z "$pyfiles" ]; then
	#Select unique files that contain formatting problem
	unformatted=$(pylama)
	if ! [ -z "$unformatted" ]; then
		status=1
		# Some files are not yapf'd. Print message and fail.
		echo >&2 "Python files must be formatted. Try to run yapf to fix these issues:"
		for fn in $pyfiles; do
			echo >&2 "yapf -i $fn"
		done
	fi
fi

exit $status
