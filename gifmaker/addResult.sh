#!/bin/bash

width=`identify -format %w $1`; \
convert -background '#0005' -fill white -gravity center -size ${width}x90 \
caption:"$2" \
$1 +swap -gravity center -composite  $1
