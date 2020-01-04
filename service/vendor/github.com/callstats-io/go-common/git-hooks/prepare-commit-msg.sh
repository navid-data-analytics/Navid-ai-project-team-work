#!/bin/bash
story_id=`git symbolic-ref --short -q HEAD | grep -o "^ch[[:digit:]]*"`
if [ ! -z "$story_id" -a "$story_id" != " " ]
then
  echo "$(cat $1) [$story_id]" > "$1"
fi
