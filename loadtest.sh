#!/bin/bash

# Load Testing Script for Backend Server
# Tests various scenarios and measures performance

BASE_URL="${BASE_URL:-http://localhost:8000}"
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Backend Load Testing ===${NC}\n"

# Function to get timestamp in nanoseconds (macOS compatible)
get_nanoseconds() {
    if command -v gdate >/dev/null 2>&1; then
        # Use GNU date if available (installed via Homebrew on macOS)
        gdate +%s%N
    elif command -v python3 >/dev/null 2>&1; then
        # Use Python3 (works on both macOS and Linux)
        python3 -c 'import time; print(int(time.time() * 1e9))'
    else
        # Try native date command (works on Linux, fails gracefully on macOS)
        result=$(date +%s%N 2>/dev/null)
        if [[ "$result" =~ ^[0-9]+$ ]]; then
            echo "$result"
        else
            # Final fallback: use milliseconds (less precise but works everywhere)
            echo "$(($(date +%s) * 1000000))000"
        fi
    fi
}

# Function to check if server is running
check_server() {
    if ! curl -s "${BASE_URL}/health" > /dev/null; then
        echo -e "${RED}❌ Server is not running at ${BASE_URL}${NC}"
        echo "Please start the server first: make run"
        exit 1
    fi
    echo -e "${GREEN}✓ Server is running${NC}\n"
}

# Function to make a single request and measure time
time_request() {
    local url=$1
    local label=$2
    local start=$(get_nanoseconds)
    local response=$(curl -s -w "\n%{http_code}" -o /tmp/response.json "$url")
    local end=$(get_nanoseconds)
    local http_code=$(echo "$response" | tail -n1)
    local duration=$(( (end - start) / 1000000 )) # Convert to milliseconds
    
    if [ "$http_code" = "200" ]; then
        echo -e "${GREEN}✓${NC} $label: ${duration}ms (HTTP $http_code)" >&2
        echo "$duration"
    else
        echo -e "${RED}✗${NC} $label: HTTP $http_code" >&2
        echo "0"
    fi
}

# Function to run concurrent requests
run_concurrent() {
    local url=$1
    local count=$2
    local label=$3
    local pids=()
    local times=()
    
    echo -e "${YELLOW}Running $count concurrent requests: $label${NC}"
    
    for i in $(seq 1 $count); do
        (
            start=$(get_nanoseconds)
            curl -s -o /dev/null -w "%{http_code}" "$url" > /tmp/http_$i.code 2>/dev/null
            end=$(get_nanoseconds)
            duration=$(( (end - start) / 1000000 ))
            echo "$duration" > /tmp/time_$i.txt
        ) &
        pids+=($!)
    done
    
    # Wait for all to complete
    for pid in "${pids[@]}"; do
        wait $pid
    done
    
    # Collect results
    local total=0
    local success=0
    local times_array=()
    
    for i in $(seq 1 $count); do
        if [ -f /tmp/time_$i.txt ]; then
            time=$(cat /tmp/time_$i.txt)
            code=$(cat /tmp/http_$i.code 2>/dev/null || echo "000")
            if [ "$code" = "200" ] || [ "$code" = "503" ]; then
                times_array+=($time)
                total=$((total + time))
                if [ "$code" = "200" ]; then
                    success=$((success + 1))
                fi
            fi
            rm -f /tmp/time_$i.txt /tmp/http_$i.code
        fi
    done
    
    # Calculate percentiles
    IFS=$'\n' sorted=($(sort -n <<<"${times_array[*]}"))
    unset IFS
    
    local len=${#sorted[@]}
    if [ $len -gt 0 ]; then
        local p50_idx=$((len * 50 / 100))
        local p95_idx=$((len * 95 / 100))
        local p99_idx=$((len * 99 / 100))
        
        local p50=${sorted[$p50_idx]:-0}
        local p95=${sorted[$p95_idx]:-0}
        local p99=${sorted[$p99_idx]:-0}
        local avg=$((total / len))
        
        echo -e "  Success: $success/$count requests"
        echo -e "  Average: ${avg}ms"
        echo -e "  p50: ${p50}ms, p95: ${p95}ms, p99: ${p99}ms"
    fi
    echo ""
}

# Function to run sustained load
run_sustained_load() {
    local url=$1
    local duration=$2
    local rps=$3
    local label=$4
    # Calculate interval in microseconds for sleep
    local interval_us=$((1000000 / rps))
    
    echo -e "${YELLOW}Running sustained load: $label (${rps} RPS for ${duration}s)${NC}"
    
    local start_time=$(date +%s)
    local end_time=$((start_time + duration))
    local request_count=0
    local success_count=0
    local times=()
    
    while [ $(date +%s) -lt $end_time ]; do
        req_start=$(get_nanoseconds)
        http_code=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null)
        req_end=$(get_nanoseconds)
        duration_ms=$(( (req_end - req_start) / 1000000 ))
        
        request_count=$((request_count + 1))
        times+=($duration_ms)
        
        if [ "$http_code" = "200" ]; then
            success_count=$((success_count + 1))
        fi
        
        # Sleep in microseconds (divide by 1000000 for seconds)
        usleep $interval_us 2>/dev/null || sleep "0.$(printf "%06d" $interval_us)"
    done
    
    # Calculate metrics
    IFS=$'\n' sorted=($(sort -n <<<"${times[*]}"))
    unset IFS
    
    local len=${#sorted[@]}
    if [ $len -gt 0 ]; then
        local p50_idx=$((len * 50 / 100))
        local p95_idx=$((len * 95 / 100))
        local p99_idx=$((len * 99 / 100))
        
        local p50=${sorted[$p50_idx]:-0}
        local p95=${sorted[$p95_idx]:-0}
        local p99=${sorted[$p99_idx]:-0}
        
        local total=0
        for t in "${times[@]}"; do
            total=$((total + t))
        done
        local avg=$((total / len))
        
        echo -e "  Total Requests: $request_count"
        echo -e "  Successful: $success_count ($(( success_count * 100 / request_count ))%)"
        echo -e "  Average Latency: ${avg}ms"
        echo -e "  p50: ${p50}ms, p95: ${p95}ms, p99: ${p99}ms"
    fi
    echo ""
}

# Check server
check_server

echo -e "${BLUE}=== Test 1: Cold Cache (First Request) ===${NC}"
time_request "${BASE_URL}/api/component/welcome?lang=en" "Welcome EN (cold)"

echo -e "${BLUE}=== Test 2: Warm Cache (Second Request) ===${NC}"
time_request "${BASE_URL}/api/component/welcome?lang=en" "Welcome EN (warm)"

echo -e "${BLUE}=== Test 3: Different Languages ===${NC}"
time_request "${BASE_URL}/api/component/welcome?lang=es" "Welcome ES"
time_request "${BASE_URL}/api/component/welcome?lang=fr" "Welcome FR"
time_request "${BASE_URL}/api/component/welcome?lang=de" "Welcome DE"

echo -e "${BLUE}=== Test 4: Different Components ===${NC}"
time_request "${BASE_URL}/api/component/navigation?lang=en" "Navigation EN"
time_request "${BASE_URL}/api/component/user_profile?lang=en" "User Profile EN"
time_request "${BASE_URL}/api/component/footer?lang=en" "Footer EN"

echo -e "${BLUE}=== Test 5: Concurrent Requests (Within Limit) ===${NC}"
run_concurrent "${BASE_URL}/api/component/welcome?lang=en" 10 "10 concurrent (within default limit=20)"

echo -e "${BLUE}=== Test 6: Concurrent Requests (Beyond Limit) ===${NC}"
run_concurrent "${BASE_URL}/api/component/welcome?lang=en" 25 "25 concurrent (exceeds default limit=20, expect some 503s)"

echo -e "${BLUE}=== Test 7: Sustained Load (10 RPS for 10s) ===${NC}"
run_sustained_load "${BASE_URL}/api/component/welcome?lang=en" 10 10 "Welcome EN"

echo -e "${BLUE}=== Test 8: Mixed Component/Language Load ===${NC}"
echo "Running mixed requests..."
mixed_times=()
components=("welcome" "navigation" "user_profile" "footer")
langs=("en" "es" "fr" "de")

for i in {1..20}; do
    comp=${components[$RANDOM % ${#components[@]}]}
    lang=${langs[$RANDOM % ${#langs[@]}]}
    time=$(time_request "${BASE_URL}/api/component/${comp}?lang=${lang}" "${comp} ${lang}")
    # Filter to ensure only numeric values are stored
    if [[ "$time" =~ ^[0-9]+$ ]] && [ "$time" != "0" ]; then
        mixed_times+=($time)
    fi
done

if [ ${#mixed_times[@]} -gt 0 ]; then
    IFS=$'\n' sorted=($(sort -n <<<"${mixed_times[*]}"))
    unset IFS
    len=${#sorted[@]}
    p50=${sorted[$((len * 50 / 100))]:-0}
    p95=${sorted[$((len * 95 / 100))]:-0}
    p99=${sorted[$((len * 99 / 100))]:-0}
    total=0
    for t in "${mixed_times[@]}"; do
        # Ensure t is numeric before arithmetic
        if [[ "$t" =~ ^[0-9]+$ ]]; then
            total=$((total + t))
        fi
    done
    avg=$((total / len))
    
    echo -e "${GREEN}Mixed Load Summary:${NC}"
    echo -e "  Average: ${avg}ms"
    echo -e "  p50: ${p50}ms, p95: ${p95}ms, p99: ${p99}ms"
fi

echo -e "\n${BLUE}=== Test Complete ===${NC}"

