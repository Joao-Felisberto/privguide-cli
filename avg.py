import statistics
import json

def process_cpu_profile(name: str):
    with open(f"profiles/res/cpu_{name}.json") as f:
        data = json.loads(f.read())
    
    durations = {float(d["duration"][:-3]) for d in data}

    median_time = statistics.median(durations)
    average_time = statistics.mean(durations)
    std_deviation = statistics.stdev(durations)
    
    print(f"Median Time: {median_time:.4f}ms")
    print(f"Average Time: {average_time:.4f}ms")
    print(f"Standard Deviation: {std_deviation:.2f}ms")

def to_kb(space: str) -> float:
    if space.endswith("MB"):
        return 1000 * float(space[:-2])
    if space.endswith("kB"):
        return float(space[:-2])
    return float(space)

def process_mem_profile(name: str):
    with open(f"profiles/res/mem_{name}.json") as f:
        data = json.loads(f.read())
    
    memories = {to_kb(d["total"]) for d in data}

    median_mem = statistics.median(memories)
    average_mem = statistics.mean(memories)
    std_deviation = statistics.stdev(memories)
    
    print(f"Median Memory Consumption: {median_mem:.4f}kB")
    print(f"Average Memory Consumption: {average_mem:.4f}kB")
    print(f"Standard Deviation: {std_deviation:.2f}kB")

if __name__ == "__main__":
    print("Chat App")
    # process_cpu_profile("chat_app")
    process_mem_profile("chat_app")
    
    print("File Sharing")
    # process_cpu_profile("file_sharing")
    process_mem_profile("file_sharing")
    
    print("IoT")
    # process_cpu_profile("iot")
    process_mem_profile("iot")

# import re
# import statistics
# 
# def parse_time_output(file_path):
#     real_times = []
#     
#     # Define the regex pattern to match the real time output
#     real_time_pattern = re.compile(r'^real\s+(\d+)m(\d+\.\d+)s$')
#     
#     with open(file_path, 'r') as file:
#         for line in file:
#             match = real_time_pattern.match(line.strip())
#             if match:
#                 minutes = int(match.group(1))
#                 seconds = float(match.group(2))
#                 total_seconds = minutes * 60 + seconds
#                 real_times.append(total_seconds)
#     
#     if not real_times:
#         print("No 'real' time entries found.")
#         return
# 
#     # Calculate median and average
#     median_time = statistics.median(real_times)
#     average_time = statistics.mean(real_times)
#     std_deviation = statistics.stdev(real_times)
#     
#     print(f"Median Time: {median_time:.4f} seconds")
#     print(f"Average Time: {average_time:.4f} seconds")
#     print(f"Standard Deviation: {std_deviation:.2f} seconds")
# 
# # Example usage
# print("Analysis")
# file_path = './times_analysis.txt'
# parse_time_output(file_path)
# print("Test")
# file_path = './times_test.txt'
# parse_time_output(file_path)
