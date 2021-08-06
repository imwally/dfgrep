# dfgrep
Dumb, fast grep. Not dumb fast grep.

## What?
I had a recent technical interview where I was asked to write a method that can search a large dataset for a substring. It should be fast and memory efficient. Well at first I wrote something terrible, and after staring at the mess of code I decided to clean it up a bit. This was the result. It's still terrible, but less so.

## Dumb
It's dumb because it doesn't support the 're' in grep.

## Fast
It's fast because goroutines go brrrr.

## Benchmarks
```
% du -sh maildir 
2.5G	maildir
% time grep -r "Wally" maildir 1>/dev/null
grep -r "Wally" maildir > /dev/null  7.21s user 13.09s system 36% cpu 54.902 total
% time dfgrep "Wally" maildir 1>/dev/null
dfgrep "Wally" maildir > /dev/null  6.64s user 15.09s system 217% cpu 9.988 total
```