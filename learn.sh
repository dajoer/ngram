#!/bin/zsh
grep -h '^:' data/* | sentsplit | ngram --learn langModel.json
