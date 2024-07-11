import matplotlib.pyplot as plt
import json


def clean(num):
    if num.endswith("ms") or num.endswith("B"):
        return float(num[:-2])
    if num.endswith("%"):
        return float(num[:-1])

    return float(num)

def cpu():
    with open("cpu_res.json") as f:
        raw = json.loads(f.read())

    duration = [float(d["duration"][:-3]) for d in raw]

    for e in raw:
        e.pop("duration")
        e.pop("total samples")

    methods = {}

    for e in raw:
        for m in e:
            if m not in methods:
                methods[m] = {
                    "flat": [], 
                    "flat%": [], 
                    "sum%": [], 
                    "cum": [], 
                    "cum%": [], 
                }

            for k in ("flat", "flat%", "sum%", "cum", "cum%"):
                methods[m][k].append(clean(e[m][k]))
    
    n_methods = len(methods)
    top_flat = sorted(
        [(k, methods[k]["flat"]) for k in methods], 
        key = lambda m_l: sum(methods[m_l[0]]["flat"])/n_methods, 
        reverse = True
    )

    fig, (ax1, ax2, ax3, ax4) = plt.subplots(1, 4)

    ax1.violinplot(duration, showmedians=True)
    ax1.set_title("Execution time")
    ax1.set_ylabel("time (ms)")

    for (m, flat), ax in zip(top_flat[:3], (ax2, ax3, ax4)):
        ax.violinplot(flat, showmedians=True)
        ax.set_title(f"Flat '{m}'")
        ax.set_ylabel("time (ms)")

    plt.show()

def mem():
    with open("mem_res_2.json") as f:
        raw = json.loads(f.read())

    total = [clean(d["total"]) for d in raw]

    for e in raw:
        e.pop("total")

    methods = {}

    for e in raw:
        for m in e:
            if m not in methods:
                methods[m] = {
                    "flat": [], 
                    "flat%": [], 
                    "sum%": [], 
                    "cum": [], 
                    "cum%": [], 
                }

            for k in ("flat", "flat%", "sum%", "cum", "cum%"):
                methods[m][k].append(clean(e[m][k]))
    
    n_methods = len(methods)
    top_flat = sorted(
        [(k, methods[k]["flat"]) for k in methods], 
        key = lambda m_l: sum(methods[m_l[0]]["flat"])/n_methods, 
        reverse = True
    )

    fig, (ax1, ax2, ax3, ax4) = plt.subplots(1, 4)

    ax1.violinplot(total, showmedians=True)
    ax1.set_title("Total memory usage")
    ax1.set_ylabel("memory (kB)")

    for (m, flat), ax in zip(top_flat[:3], (ax2, ax3, ax4)):
        ax.violinplot(flat, showmedians=True)
        ax.set_title(f"Flat '{m}'")
        ax.set_ylabel("time (ms)")

    plt.show()

if __name__ == "__main__":
    cpu()
    mem()