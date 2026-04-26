#!/bin/bash

# Define a cleanup function
cleanup() {
  # Only reset if we actually staged something with -N
  # and didn't finish the commit.
  git reset . > /dev/null 2>&1
}

# 'trap' catches interrupts (Ctrl+C) or script exits
# and runs the cleanup function.
trap cleanup EXIT

git add -N .
DIFF=$(git diff HEAD)

if [ -z "$DIFF" ]; then
  echo "No changes detected."
  exit 0
fi

MESSAGE=$(echo "$DIFF" | opencode run "Generate a git commit message. Return ONLY the text.")

if [ -n "$MESSAGE" ]; then
  echo -e "\nProposed Message: $MESSAGE"
  read -p "Apply this commit? (y/n): " CONFIRM
  if [ "$CONFIRM" = "y" ]; then
    git add .
    git commit -m "$MESSAGE"
    # Disable the trap so it doesn't reset our successful commit
    trap - EXIT 
  fi
fi