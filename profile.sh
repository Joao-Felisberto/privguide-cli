
for ex in examples/chat_app/ examples/file_sharing/ examples/iot/; do
    for i in $(seq -f "%05g" $(ls profiles/cpu/ | wc -l) 10000); do
    # for i in $(seq -f "%05g" $(ls profiles/cpu/ | wc -l) 10); do
        ./devprivops analyse name pass localhost 3030 devprivops \
            --local-dir $ex \
            --global-dir examples/global/ \
            --profile cpu > /dev/null
        mv "cpu.prof" "profiles/cpu/cpu_$i.prof"
        go tool pprof -text "profiles/cpu/cpu_$i.prof" > "profiles/cpu/cpu_$i.txt"

        ./devprivops analyse name pass localhost 3030 devprivops \
            --local-dir $ex \
            --global-dir examples/global/ \
            --profile mem > /dev/null
        mv "mem.prof" "profiles/mem/mem_$i.prof"
        go tool pprof -text "profiles/mem/mem_$i.prof" > "profiles/mem/mem_$i.txt"

        echo "$ex $i"
    done
    
    python3 extract_profile_data.py $(basename $ex)
    rm profiles/{cpu,mem}/*
done