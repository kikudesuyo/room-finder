#!/usr/bin/env bash

set -euo pipefail

script_dir=$(cd "$(dirname "$0")" && pwd)
fixture=$(mktemp -d "${TMPDIR:-/tmp}/room-finder-rental-search.XXXXXX")
trap 'rm -rf "$fixture"' EXIT

printf '%s\n' one two three four > "$fixture/candidates.txt"

  "$script_dir/batch-fetch.sh" \
  --source test-source \
  --input "$fixture/candidates.txt" \
  --state "$fixture/state" \
  --batch-size 2 \
  -- 'case "{item}" in one|three) exit 0 ;; two) exit 20 ;; four) exit 10 ;; esac'

test "$(wc -l < "$fixture/state" | tr -d ' ')" = 3

  "$script_dir/batch-fetch.sh" \
  --source test-source \
  --input "$fixture/candidates.txt" \
  --state "$fixture/state" \
  --batch-size 2 \
  -- 'test "{item}" = four && exit 0'

test "$(wc -l < "$fixture/state" | tr -d ' ')" = 4

set +e
"$script_dir/batch-fetch.sh" \
  --source test-source \
  --input "$fixture/candidates.txt" \
  --state "$fixture/error-state" \
  -- 'test "{item}" = one && exit 1'
exit_code=$?
set -e
test "$exit_code" = 1

printf 'rental search batch fixture tests passed\n'
