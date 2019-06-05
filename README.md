# Queue sub request/reply example

## Terminal Window 1

```text
$ cd replier
$ go run main.go
```

## Teminal Window 2

```text
$ cd requestor
$ go run main.go
```

## Terminal Window 3

When ready to switch instances, launch another replier and ctrl-C the previous
one (Terminal window 1)

```text
$ cd replier
$ go run main.go
```


