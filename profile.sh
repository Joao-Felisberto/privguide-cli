for i in $(seq -f "%05g" $(ls profiles/cpu/ | wc -l) 10000); do
    ./devprivops analyse name pass localhost 3030 devprivops --profile cpu > /dev/null
    mv "cpu.prof" "profiles/cpu/cpu_$i.prof"
    go tool pprof -text "profiles/cpu/cpu_$i.prof" > "profiles/cpu/cpu_$i.txt"

    ./devprivops analyse name pass localhost 3030 devprivops --profile mem > /dev/null
    mv "mem.prof" "profiles/mem/mem_$i.prof"
    go tool pprof -text "profiles/mem/mem_$i.prof" > "profiles/mem/mem_$i.txt"

    echo $i
done