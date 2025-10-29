#!/bin/bash

# Comprehensive test script for all examples, exercises, and projects

echo "╔════════════════════════════════════════════════════════════╗"
echo "║                                                            ║"
echo "║        Testing All Go Examples, Exercises & Projects      ║"
echo "║                                                            ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

TOTAL_PASS=0
TOTAL_FAIL=0

# Function to test a Go file
test_go_file() {
    local file=$1
    local name=$2
    
    if [ -f "$file" ]; then
        echo -n "Testing: $name ... "
        if go run "$file" > /dev/null 2>&1; then
            echo "✓ PASS"
            ((TOTAL_PASS++))
            return 0
        else
            echo "✗ FAIL"
            ((TOTAL_FAIL++))
            return 1
        fi
    fi
}

# Function to test a directory with main.go
test_directory() {
    local dir=$1
    local name=$2
    
    if [ -d "$dir" ] && [ -f "$dir/main.go" ]; then
        echo -n "Testing: $name ... "
        if cd "$dir" && go run main.go > /dev/null 2>&1; then
            echo "✓ PASS"
            ((TOTAL_PASS++))
            cd - > /dev/null
            return 0
        else
            echo "✗ FAIL"
            ((TOTAL_FAIL++))
            cd - > /dev/null
            return 1
        fi
    fi
}

echo "════════════════════════════════════════════════════════════"
echo "1. Testing Examples (12 examples)"
echo "════════════════════════════════════════════════════════════"
echo ""

for i in {01..12}; do
    if [ -d "basic/examples/${i}_"* ]; then
        dir=$(ls -d basic/examples/${i}_* 2>/dev/null | head -1)
        name=$(basename "$dir")
        test_directory "$dir" "Example: $name"
    fi
done

echo ""
echo "════════════════════════════════════════════════════════════"
echo "2. Testing Exercise Solutions"
echo "════════════════════════════════════════════════════════════"
echo ""

# Basic exercises
test_go_file "basic/exercises/01_basics/fizzbuzz_solution.go" "FizzBuzz Solution"
test_go_file "basic/exercises/01_basics/reverse_string_solution.go" "Reverse String Solution"

# Intermediate exercises
test_go_file "basic/exercises/02_intermediate/url_shortener_solution.go" "URL Shortener Solution"

# Advanced exercises
test_go_file "basic/exercises/03_advanced/web_crawler_solution.go" "Web Crawler Solution"

echo ""
echo "════════════════════════════════════════════════════════════"
echo "3. Testing Projects"
echo "════════════════════════════════════════════════════════════"
echo ""

# Test calculator (with automated input)
if [ -f "basic/projects/calculator/main.go" ]; then
    echo -n "Testing: Calculator Project ... "
    if (cd basic/projects/calculator && echo -e "1\n10\n5\nq" | go run main.go > /dev/null 2>&1); then
        echo "✓ PASS"
        ((TOTAL_PASS++))
    else
        echo "✗ FAIL"
        ((TOTAL_FAIL++))
    fi
fi

# Test todo list (with automated input)
if [ -f "basic/projects/todo_list/main.go" ]; then
    echo -n "Testing: Todo List Project ... "
    if (cd basic/projects/todo_list && echo -e "2\nq" | go run main.go > /dev/null 2>&1); then
        echo "✓ PASS"
        ((TOTAL_PASS++))
    else
        echo "✗ FAIL"
        ((TOTAL_FAIL++))
    fi
fi

echo ""
echo "════════════════════════════════════════════════════════════"
echo "Summary"
echo "════════════════════════════════════════════════════════════"
echo ""
echo "  Total Tests: $((TOTAL_PASS + TOTAL_FAIL))"
echo "  ✓ Passed:    $TOTAL_PASS"
echo "  ✗ Failed:    $TOTAL_FAIL"
echo ""

if [ $TOTAL_FAIL -eq 0 ]; then
    echo "════════════════════════════════════════════════════════════"
    echo "           🎉 All tests passed! 🎉"
    echo "════════════════════════════════════════════════════════════"
    exit 0
else
    echo "════════════════════════════════════════════════════════════"
    echo "           ⚠️  Some tests failed"
    echo "════════════════════════════════════════════════════════════"
    exit 1
fi

