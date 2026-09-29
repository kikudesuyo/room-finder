#!/usr/bin/env bash

set -euo pipefail

usage() {
  printf 'usage: %s --source SOURCE --input FILE --state FILE [--batch-size N] -- COMMAND\n' "$0" >&2
}

input_file=''
state_file=''
source=''
batch_size=20

while (($# > 0)); do
  case "$1" in
    --input)
      (($# >= 2)) || { usage; exit 2; }
      input_file=$2
      shift 2
      ;;
    --state)
      (($# >= 2)) || { usage; exit 2; }
      state_file=$2
      shift 2
      ;;
    --source)
      (($# >= 2)) || { usage; exit 2; }
      source=$2
      shift 2
      ;;
    --batch-size)
      (($# >= 2)) || { usage; exit 2; }
      batch_size=$2
      shift 2
      ;;
    --)
      shift
      break
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      printf 'unknown argument: %s\n' "$1" >&2
      usage
      exit 2
      ;;
  esac
done

[[ -n "$input_file" && -f "$input_file" ]] || { printf 'input file not found\n' >&2; exit 2; }
[[ -n "$state_file" ]] || { printf 'state file is required\n' >&2; exit 2; }
[[ "$source" =~ ^[a-z0-9-]+$ ]] || { printf 'source must contain lowercase letters, digits, or hyphens\n' >&2; exit 2; }
[[ "$batch_size" =~ ^[1-9][0-9]*$ ]] || { printf 'batch size must be positive\n' >&2; exit 2; }
(($# > 0)) || { printf 'command is required\n' >&2; exit 2; }

command_template=$*
case "$command_template" in
  *'{item}'*) ;;
  *) printf 'command must contain the {item} placeholder\n' >&2; exit 2 ;;
esac

mkdir -p "$(dirname "$state_file")"
touch "$state_file"

processed=0
succeeded=0
retryable=0
rejected=0
failed=0
batch_count=0

while IFS= read -r item || [[ -n "$item" ]]; do
  [[ -z "$item" || "$item" == \#* ]] && continue
  rg -Fqx -- "$item" "$state_file" && continue

  if ((processed % batch_size == 0)); then
    batch_count=$((batch_count + 1))
    printf 'source=%s batch=%d processed=%d succeeded=%d retryable=%d rejected=%d pending=%s\n' \
      "$source" "$batch_count" "$processed" "$succeeded" "$retryable" "$rejected" "$((processed + 1))"
  fi

  command=${command_template//\{item\}/\"\$ROOM_FINDER_ITEM\"}
  set +e
  ROOM_FINDER_ITEM="$item" bash -lc "$command"
  exit_code=$?
  set -e
  processed=$((processed + 1))

  case "$exit_code" in
    0)
      succeeded=$((succeeded + 1))
      printf '%s\n' "$item" >> "$state_file"
      ;;
    10)
      retryable=$((retryable + 1))
      ;;
    20)
      rejected=$((rejected + 1))
      printf '%s\n' "$item" >> "$state_file"
      ;;
    *)
      failed=$((failed + 1))
      printf 'batch stopped: source=%s item=%s exit=%d processed=%d succeeded=%d retryable=%d rejected=%d pending=%s\n' \
        "$source" "$item" "$exit_code" "$processed" "$succeeded" "$retryable" "$rejected" "$((processed + 1))" >&2
      exit "$exit_code"
      ;;
  esac
done < "$input_file"

printf 'batch complete: source=%s processed=%d succeeded=%d retryable=%d rejected=%d failed=%d\n' \
  "$source" "$processed" "$succeeded" "$retryable" "$rejected" "$failed"
