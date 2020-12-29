#!/bin/sh

# Timeout in seconds
TIMEOUT=99

# Location of LOCKFILE
LOCKFILE=./tmp.lock

# Try to acquire lockfile
while [ $TIMEOUT -gt 0 ]; do
  if { set -C; 2>/dev/null >$LOCKFILE; }; then

    # When lockfile is acquired set trap
    # to delete it when script finishes.
    trap "rm -f $LOCKFILE" EXIT

    # Execute the command
    echo "$1" | sh

    # After executon exit the script
    exit
  fi
  sleep 1
  TIMEOUT=$((TIMEOUT - 1))
done
