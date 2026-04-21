#!/usr/bin/env bash
# AGENT_ROLE-aware tasks.md completion check.
#
# For every openspec/changes/*/tasks.md (non-archive), locate the
# `## Stage: <name> (owner: $AGENT_ROLE)` section that belongs to the
# current agent, and fail if any `- [ ]` checkboxes remain unchecked.
#
# Runs only when AGENT_ROLE env is set (agent worktree). Human-run make
# ci-lint skips.
#
# Usage (implicit via make ci-lint):
#   AGENT_ROLE=dev-agent make ci-lint

set -euo pipefail

if [[ -z "${AGENT_ROLE:-}" ]]; then
  # human invocation — skip
  exit 0
fi

shopt -s nullglob
fail=0
for f in openspec/changes/*/tasks.md; do
  # skip archive (those are already done and landed)
  case "$f" in */archive/*) continue ;; esac

  # Extract the section whose owner matches AGENT_ROLE.
  # A section runs from its `## Stage:` heading to the next `## Stage:` or EOF.
  section=$(awk -v role="$AGENT_ROLE" '
    /^## Stage:/ {
      if (index($0, "(owner: " role ")")) { in_section=1 } else { in_section=0 }
      next
    }
    in_section { print }
  ' "$f")

  if [[ -z "$section" ]]; then
    # No section owned by this role in this file — nothing to check.
    continue
  fi

  unchecked=$(echo "$section" | grep -nE '^\s*-\s*\[\s*\]' || true)
  if [[ -n "$unchecked" ]]; then
    echo "FAIL: $f has unchecked tasks in your section (owner: $AGENT_ROLE):"
    echo "$unchecked"
    fail=1
  fi
done

if [[ $fail -eq 0 ]]; then
  echo "OK: all tasks in your section(s) (owner: $AGENT_ROLE) are checked."
fi
exit $fail
