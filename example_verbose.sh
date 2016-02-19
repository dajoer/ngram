#!/bin/bash
if [ -z $1 ]; then
	echo "Usage: ./example_verbose.sh <filename>"
else
	heapsAlg our point of second meeting | ngram -v $1 | sort -gr | less
fi
