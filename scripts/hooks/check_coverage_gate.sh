#!/usr/bin/env bash
# check_coverage_gate.sh — fail if total Go coverage is below the threshold.
# Usage: check_coverage_gate.sh <coverage.out> [threshold-percent]
set -euo pipefail

PROFILE="${1:-coverage.out}"
THRESHOLD="${2:-75}"

if [[ ! -f "${PROFILE}" ]]; then
  echo "coverage-gate: profile ${PROFILE} not found" >&2
  exit 1
fi

# `go tool cover -func` prints a trailing "total:" line like "total: (statements) 81.2%".
total="$(go tool cover -func="${PROFILE}" | awk '/^total:/ {gsub(/%/,"",$3); print $3}')"

if [[ -z "${total}" ]]; then
  echo "coverage-gate: could not parse coverage from ${PROFILE}" >&2
  exit 1
fi

# Compare with awk (no bc dependency).
if awk -v t="${total}" -v thr="${THRESHOLD}" 'BEGIN { exit !(t+0 < thr+0) }'; then
  echo "coverage-gate: FAIL — ${total}% < ${THRESHOLD}% required"
  exit 1
fi

echo "coverage-gate: OK — ${total}% >= ${THRESHOLD}%"
