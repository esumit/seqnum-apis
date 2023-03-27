#!/bin/bash

# Print the current date and time
echo "Running tests on $(date)"

# Change the directory to your Go project's seqnum_test folder
cd ./seqnum_test

# Run the test cases and save the output to a file
go test -v ./... > seqnumtest_output.txt

# Display the test output
cat seqnumtest_output.txt

