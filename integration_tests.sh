#!/bin/sh
set -e
PKGARGS="$*"

rm -rf covdatafiles
mkdir covdatafiles
export GOCOVERDIR=covdatafiles

go build -cover $BUILDARGS .

./devprivops analyse user pass 127.0.0.1 3030 tmp 
echo "================== TEST DONE!"
./devprivops analyse user pass 127.0.0.1 3030 tmp  --local-dir test_files/test_1 || true
echo "================== TEST DONE!"
./devprivops analyse user pass 127.0.0.1 3030 tmp  --local-dir test_files/test_2 || true
echo "================== TEST DONE!"
./devprivops analyse user pass 127.0.0.1 3030 tmp  --local-dir test_files/test_3
echo "================== TEST DONE!"
# ./devprivops analyse user pass 127.0.0.1 3030 tmp

rm devprivops

go tool covdata percent -i=covdatafiles
