#!/bin/bash

# Comprehensive test script for ALL examples in the repository

echo "╔════════════════════════════════════════════════════════════╗"
echo "║                                                            ║"
echo "║           Testing ALL Go Examples in Repository           ║"
echo "║                                                            ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

TOTAL_PASS=0
TOTAL_FAIL=0
TOTAL_SKIP=0

# Function to test a Go file
test_go_file() {
    local file=$1
    local name=$2
    
    if [ -f "$file" ]; then
        local dir=$(dirname "$file")
        echo -n "Testing: $name ... "
        if (cd "$dir" && timeout 5 go run "$(basename "$file")" > /dev/null 2>&1); then
            echo "✓ PASS"
            ((TOTAL_PASS++))
            return 0
        else
            echo "✗ FAIL"
            ((TOTAL_FAIL++))
            return 1
        fi
    else
        echo "⊘ SKIP: $name (file not found)"
        ((TOTAL_SKIP++))
        return 2
    fi
}

# Function to find and test all main.go files in a directory
test_directory_recursive() {
    local base_dir=$1
    local label=$2
    
    echo "════════════════════════════════════════════════════════════"
    echo "$label"
    echo "════════════════════════════════════════════════════════════"
    echo ""
    
    # Find all main.go files
    while IFS= read -r -d '' file; do
        # Get relative path for display
        local rel_path="${file#$base_dir/}"
        local dir_name=$(dirname "$rel_path")
        test_go_file "$file" "$dir_name"
    done < <(find "$base_dir" -name "main.go" -type f -print0 2>/dev/null | sort -z)
    
    echo ""
}

# Test basic examples (already tested, but include for completeness)
echo "════════════════════════════════════════════════════════════"
echo "1. Basic Examples (Interactive Learning)"
echo "════════════════════════════════════════════════════════════"
echo ""

for i in {01..12}; do
    if [ -d "basic/examples/${i}_"* ]; then
        dir=$(ls -d basic/examples/${i}_* 2>/dev/null | head -1)
        if [ -f "$dir/main.go" ]; then
            name=$(basename "$dir")
            test_go_file "$dir/main.go" "Example: $name"
        fi
    fi
done

echo ""

# Test advanced/learngo examples
if [ -d "advanced/learngo" ]; then
    test_directory_recursive "advanced/learngo" "2. Advanced LearnGo Examples"
fi

# Test advanced/learn-to-code-go-version-03 examples
if [ -d "advanced/learn-to-code-go-version-03" ]; then
    test_directory_recursive "advanced/learn-to-code-go-version-03" "3. Learn to Code Go v3 Examples"
fi

# Test course examples
if [ -d "course/code" ]; then
    test_directory_recursive "course/code" "4. Course Code Examples"
fi

echo "════════════════════════════════════════════════════════════"
echo "Summary"
echo "════════════════════════════════════════════════════════════"
echo ""
echo "  Total Tests: $((TOTAL_PASS + TOTAL_FAIL + TOTAL_SKIP))"
echo "  ✓ Passed:    $TOTAL_PASS"
echo "  ✗ Failed:    $TOTAL_FAIL"
echo "  ⊘ Skipped:   $TOTAL_SKIP"
echo ""

if [ $TOTAL_FAIL -eq 0 ]; then
    echo "════════════════════════════════════════════════════════════"
    echo "           🎉 All tests passed! 🎉"
    echo "════════════════════════════════════════════════════════════"
    exit 0
else
    echo "════════════════════════════════════════════════════════════"
    echo "           ⚠️  $TOTAL_FAIL test(s) failed"
    echo "════════════════════════════════════════════════════════════"
    exit 1
fi

