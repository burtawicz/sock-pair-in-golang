# Sock Pairing (in Golang)

[![CircleCI](https://circleci.com/gh/burtawicz/sock-pair-in-golang/tree/main.svg?style=svg)](https://circleci.com/gh/burtawicz/sock-pair-in-golang/tree/main)

## :question: Why?
Something I've always been curious about is defining algorithms for everyday activities.
A few years ago I read [Algorithms to Live By - The Computer Science of Human Decisions](https://algorithmstoliveby.com/) by Brian Christian and Tom Griffiths. 
One of the examples they explore in the book is how a college student pairs socks at random.
The random pairing approach bothered me, so I started to wonder:
* How many ways are there to pair socks? 
* Is there an optimal way to pair socks? 
* Is the optimal way to pair obvious?

## :alembic: Running the Tests
From the CLI, execute:
* `go test` to run the tests
* `go test -bench=.` to benchmark the strategies _(include the `-short` tag to skip the long-running benchmarks)_
