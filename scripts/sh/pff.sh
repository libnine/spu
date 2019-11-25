#! /usr/bin/env sh

js="./scripts/js/"

mongo $PFF < $js/pffclear.js
node "$js/app.js"
mongo $CUR --quiet $js/pfflt_onepct.js > pfflt_onepct.json
mongo $CUR --quiet $js/pffgt_onepct.js > pffgt_onepct.json
mongo $CUR --quiet $js/pff52wk_high.js > pff52wk_high.json
mongo $CUR --quiet $js/pff52wk_low.js > pff52wk_low.json
