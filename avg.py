import re
import statistics

def parse_time_output(file_path):
    real_times = []
    
    # Define the regex pattern to match the real time output
    real_time_pattern = re.compile(r'^real\s+(\d+)m(\d+\.\d+)s$')
    
    with open(file_path, 'r') as file:
        for line in file:
            match = real_time_pattern.match(line.strip())
            if match:
                minutes = int(match.group(1))
                seconds = float(match.group(2))
                total_seconds = minutes * 60 + seconds
                real_times.append(total_seconds)
    
    if not real_times:
        print("No 'real' time entries found.")
        return

    # Calculate median and average
    median_time = statistics.median(real_times)
    average_time = statistics.mean(real_times)
    std_deviation = statistics.stdev(real_times)
    
    print(f"Median Time: {median_time:.4f} seconds")
    print(f"Average Time: {average_time:.4f} seconds")
    print(f"Standard Deviation: {std_deviation:.2f} seconds")

# Example usage
print("Analysis")
file_path = './times_analysis.txt'
parse_time_output(file_path)
print("Test")
file_path = './times_test.txt'
parse_time_output(file_path)
