#!/bin/bash
if [ -z $1 ]; then
	echo "Usage: ./learn.sh <filename>"
else
	grep -h '^:' data/* | sed -e 's/[.:;!?] /\n/g' | tr '[:upper:]ÄÖÜ' '[:lower:]äöü' | sed -e 's/[^-a-z0-9_ äöüß]//g' | ngram --learn $1
fi
