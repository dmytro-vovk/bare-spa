#!/bin/bash

set -eu

npm install

npm run build-guest
npm run build-user

rm -rf ./internal/webserver/handlers/home/css

cp -rf ./frontend/styles ./internal/webserver/handlers/home/css

gzip -c ./frontend/index.html > ./internal/webserver/handlers/home/index.html.gz

gsed -i 's/=guest.js.map/=js.map/gI' ./internal/webserver/handlers/home/guest.js
gzip -f ./internal/webserver/handlers/home/guest.js
gzip -f ./internal/webserver/handlers/home/guest.js.map

gsed -i 's/=user.js.map/=js.map/gI' ./internal/webserver/handlers/home/user.js
gzip -f ./internal/webserver/handlers/home/user.js
gzip -f ./internal/webserver/handlers/home/user.js.map
