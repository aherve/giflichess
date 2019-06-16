#!/bin/bash

convert $1 \
  -gravity SouthEast -splice 0x18 \
  -annotate +0+2 'gifchess.com  ' $1 && \
  convert $1 \
  -gravity SouthWest \
  -annotate +0+2 "$(printf '\xa0')$(printf '\xa0') $2 vs $3" $1
