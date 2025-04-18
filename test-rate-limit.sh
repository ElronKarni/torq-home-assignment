#!/bin/bash

# Test the rate limiting by sending many requests in a short time
# The rate limit is set to 100 requests per second by default

echo "Testing rate limiting..."
echo "Sending 150 requests in rapid succession..."

RATE_LIMIT=100
TOTAL_REQUESTS=150

# Track counts of each HTTP status code
HTTP_200_COUNT=0
HTTP_429_COUNT=0

for i in $(seq 1 $TOTAL_REQUESTS); do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:8080/v1/find-country?ip=1.1.1.1")
  
  if [ "$STATUS" == "200" ]; then
    HTTP_200_COUNT=$((HTTP_200_COUNT + 1))
  elif [ "$STATUS" == "429" ]; then
    HTTP_429_COUNT=$((HTTP_429_COUNT + 1))
  fi
  
  # Display progress
  echo -ne "Progress: $i/$TOTAL_REQUESTS\r"
done

echo -e "\nResults:"
echo "Successful requests (200): $HTTP_200_COUNT"
echo "Rate limited requests (429): $HTTP_429_COUNT"

if [ $HTTP_429_COUNT -gt 0 ]; then
  echo "Rate limiting is working!"
else
  echo "Rate limiting did not trigger. Try increasing the number of requests."
fi 