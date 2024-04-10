#!/bin/bash

make install
pid=0

kill_and_run() {
  if [ $pid -ne 0 ]; then
    kill -SIGTERM $pid
  else
    killall -SIGTERM sites
  fi
  sites run  &
  pid=$!
}

trap kill_and_run SIGHUP

kill_and_run

while inotifywait -P -r -e modify,move,create,delete --exclude '(.swp|.swx|.swo|.tmp|.bak|.sqlite3)|(^.\/.git\/)' .; do

  # we always need to rebuild the binary...
  make install

  kill -HUP $$ # Send SIGHUP to self
done
