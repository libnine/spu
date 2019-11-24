#! /usr/bin/env sh

js="./scripts/js/"

mongo $PFF < "$js/pffclear.js" &> sh.log 
node "$js/app.js" &> sh.log
mongo $PFF < "$js/pfflt_onepct.js"

