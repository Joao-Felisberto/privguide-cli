import os
import sys
import json

CPU_DIR = "profiles/cpu"
MEM_DIR = "profiles/mem"

if __name__ == "__main__":
#    cpu_data = []
#    
#    for profile_f in [f for f in os.listdir(CPU_DIR) if f.endswith(".txt")]:
#        profile_f = f"{CPU_DIR}/{profile_f}"
#        print(profile_f)
#        with open(profile_f) as f:
#            data = [ln.split() for ln in f.read().splitlines()]
#        
#        print(data[3])
#        if len(data[3]) == 7:
#            tot_samples = data[3][6][1:-1]
#        elif len(data[3]) == 8:
#            tot_samples = data[3][7][:-1]
#        else:
#            tot_samples = "0"
#        datum = {
#            "duration": data[3][1],
#            "total samples": (data[3][5], tot_samples),
#            **{
#                ' '.join(d[5:]): {
#                    "flat": d[0],
#                    "flat%": d[1],
#                    "sum%": d[2],
#                    "cum": d[3],
#                    "cum%": d[4],
#                } for d in data[6:]
#            }
#        }
#
#        # print(json.dumps(datum, indent=2))
#        cpu_data.append(datum)
#
#    with open("cpu_res.json", 'w') as f:
#        f.write(json.dumps(cpu_data))
#        print(len(cpu_data))
#
#    del(cpu_data)
    mem_data = []
    
    for profile_f in [f for f in os.listdir(MEM_DIR) if f.endswith(".txt")]:
        profile_f = f"{MEM_DIR}/{profile_f}"
        print(profile_f)
        with open(profile_f) as f:
            data = [ln.split() for ln in f.read().splitlines()]
        
        # print(data[3])
        datum = {
            "total": data[3][-2],
            **{
                ' '.join(d[5:]): {
                    "flat": d[0],
                    "flat%": d[1],
                    "sum%": d[2],
                    "cum": d[3],
                    "cum%": d[4],
                } for d in data[6:]
            }
        }

        # print(json.dumps(datum, indent=2))
        mem_data.append(datum)

    with open("mem_res.json", 'w') as f:
        f.write(json.dumps(mem_data))
        print(len(mem_data))