# fly-dist-sys

## echo

```zsh
./maelstrom test -w echo --bin ~/go/bin/echo --node-count 1 --time-limit 10
```

## unique-ids

```zsh
./maelstrom test -w unique-ids --bin ~/go/bin/unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition 
```