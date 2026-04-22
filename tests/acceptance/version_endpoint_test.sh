#!/bin/bash
# Integration test for /version endpoint
# Uses Docker Compose for real-link testing with actual server instance

set -euo pipefail

TEST_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DOCKER_COMPOSE_FILE="${TEST_DIR}/docker-compose.test.yml"
CONTAINER_NAME="ubox-crosser-version-test"
ADMIN_ADDR="localhost:8080"
TIMEOUT_SECS=10

echo "=== Starting /version Endpoint Integration Test ==="

# Create docker-compose test file if it doesn't exist
if [ ! -f "$DOCKER_COMPOSE_FILE" ]; then
  echo "[INFO] Creating docker-compose test file..."
  cat > "$DOCKER_COMPOSE_FILE" << 'EOF'
version: '3.8'
services:
  server:
    container_name: ubox-crosser-version-test
    build:
      context: .
      dockerfile: Dockerfile
    command: /app/server --admin-addr :8080
    ports:
      - "8080:8080"
    networks:
      - test-net
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/version"]
      interval: 2s
      timeout: 5s
      retries: 5
      start_period: 2s

networks:
  test-net:
    driver: bridge
EOF
fi

# Cleanup function
cleanup() {
  echo "[INFO] Cleaning up test containers..."
  docker-compose -f "$DOCKER_COMPOSE_FILE" down || true
}

trap cleanup EXIT

echo "[INFO] Building and starting server..."
docker-compose -f "$DOCKER_COMPOSE_FILE" up -d

echo "[INFO] Waiting for server to be ready..."
for i in $(seq 1 $TIMEOUT_SECS); do
  if curl -sf http://$ADMIN_ADDR/version > /dev/null 2>&1; then
    echo "[OK] Server is ready"
    break
  fi
  if [ "$i" -eq "$TIMEOUT_SECS" ]; then
    echo "[ERROR] Server failed to start within ${TIMEOUT_SECS}s"
    docker-compose -f "$DOCKER_COMPOSE_FILE" logs
    exit 1
  fi
  sleep 1
done

echo ""
echo "=== Test Suite: /version Endpoint ==="
echo ""

PASS=0
FAIL=0

# Test 1: GET /version returns 200
echo "[TEST 1] GET /version returns 200 OK"
RESPONSE=$(curl -s -w "\n%{http_code}" http://$ADMIN_ADDR/version)
STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$STATUS" = "200" ]; then
  echo "[PASS] Status code: $STATUS"
  ((PASS++))
else
  echo "[FAIL] Expected 200, got $STATUS"
  echo "Response: $BODY"
  ((FAIL++))
fi

# Test 2: Response is valid JSON
echo "[TEST 2] Response is valid JSON with 'commit' field"
if echo "$BODY" | jq -e '.commit' > /dev/null 2>&1; then
  COMMIT=$(echo "$BODY" | jq -r '.commit')
  echo "[PASS] commit field found: $COMMIT"
  ((PASS++))
else
  echo "[FAIL] Response is not valid JSON or missing 'commit' field"
  echo "Response: $BODY"
  ((FAIL++))
fi

# Test 3: Commit format validation
echo "[TEST 3] Commit field is either 'unknown' or valid hex hash"
if [[ "$COMMIT" == "unknown" ]]; then
  echo "[PASS] commit is 'unknown' (expected when not injected)"
  ((PASS++))
elif [[ "$COMMIT" =~ ^[0-9a-f]{40}$ ]]; then
  echo "[PASS] commit is valid 40-char hex hash"
  ((PASS++))
else
  echo "[FAIL] commit format invalid: $COMMIT"
  ((FAIL++))
fi

# Test 4: POST /version returns 405
echo "[TEST 4] POST /version returns 405 Method Not Allowed"
POST_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST http://$ADMIN_ADDR/version)
POST_STATUS=$(echo "$POST_RESPONSE" | tail -n1)

if [ "$POST_STATUS" = "405" ]; then
  echo "[PASS] POST returns 405"
  ((PASS++))
else
  echo "[FAIL] Expected 405, got $POST_STATUS"
  ((FAIL++))
fi

# Test 5: Content-Type is JSON
echo "[TEST 5] Response Content-Type is application/json"
CONTENT_TYPE=$(curl -s -I http://$ADMIN_ADDR/version | grep -i "content-type" | cut -d' ' -f2-)
if [[ "$CONTENT_TYPE" == *"application/json"* ]]; then
  echo "[PASS] Content-Type: $CONTENT_TYPE"
  ((PASS++))
else
  echo "[FAIL] Expected application/json, got: $CONTENT_TYPE"
  ((FAIL++))
fi

# Test 6: Response time is acceptable (< 100ms)
echo "[TEST 6] Response time is under 100ms"
START=$(date +%s%N)
curl -s http://$ADMIN_ADDR/version > /dev/null
END=$(date +%s%N)
ELAPSED_MS=$(( (END - START) / 1000000 ))

if [ "$ELAPSED_MS" -lt 100 ]; then
  echo "[PASS] Response time: ${ELAPSED_MS}ms"
  ((PASS++))
else
  echo "[WARN] Response time: ${ELAPSED_MS}ms (expected < 100ms)"
  # Don't fail on timing, just warn
fi

echo ""
echo "=== Test Results ==="
echo "PASSED: $PASS"
echo "FAILED: $FAIL"
echo ""

if [ "$FAIL" -gt 0 ]; then
  exit 1
fi

echo "[SUCCESS] All /version endpoint tests passed!"
exit 0
