#!/bin/bash
if [ -z $1 ]; then
	echo "Usage: ./learn.sh <filename>"
else
	grep -h '^:' data/* | sentsplit | ngram --learn $1
fi
