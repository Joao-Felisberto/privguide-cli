#!/bin/sh
set -e
PKGARGS="$*"

rm -rf covdatafiles
mkdir covdatafiles
export GOCOVERDIR=covdatafiles

#!/bin/bash

# Start HTTP server in the background
python3 -m http.server 8000 &
#python3 -m SimpleHTTPServer 8000 &>/dev/null &

# Store the PID of the server process
SERVER_PID=$!

# Wait for the server to start
# sleep 1

# Send a request to the server
# curl http://localhost:8000

go build -cover $BUILDARGS .

./devprivops schema attack-tree > res.json
cmp res.json schema/schemas/atk-tree-schema.json

./devprivops schema query > res.json
cmp res.json schema/schemas/query-schema.json

./devprivops schema report > res.json
cmp res.json schema/schemas/report_data-schema.json

./devprivops schema requirement > res.json
cmp res.json schema/schemas/requirement-schema.json

rm res.json

./devprivops analyse user pass 127.0.0.1 3030 tmp --report-endpoint http://localhost:8000
echo "================== TEST DONE!"
./devprivops test user pass 127.0.0.1 3030 tmp 
echo "================== TEST DONE!"
./devprivops analyse user pass 127.0.0.1 3030 tmp  --local-dir test_files/test_1 || true
echo "================== TEST DONE!"
./devprivops analyse user pass 127.0.0.1 3030 tmp  --local-dir test_files/test_2 || true
echo "================== TEST DONE!"
./devprivops analyse user pass 127.0.0.1 3030 tmp  --local-dir test_files/test_3
echo "================== TEST DONE!"
./devprivops analyse user pass 127.0.0.1 3030 tmp  --local-dir test_files/test_4 || true
echo "================== TEST DONE!"
./devprivops analyse user pass 127.0.0.1 3030 tmp  --local-dir test_files/test_5 || true
echo "================== TEST DONE!"
./devprivops analyse user pass 127.0.0.1 3030 tmp  --local-dir test_files/test_6 || true
echo "================== TEST DONE!"
./devprivops analyse user pass 127.0.0.1 3030 tmp  --local-dir test_files/test_7 || true
echo "================== TEST DONE!"
./devprivops test user pass 127.0.0.1 3030 tmp  --local-dir test_files/test_7 || true
echo "================== TEST DONE!"

# Close the server
kill $SERVER_PID

rm devprivops

go tool covdata percent -i=covdatafiles
# go tool covdata func -i=covdatafiles
