# horizon

Horizon is a new blockchain. It is written in golang and uses a novel networking stack built from two primitives: channels and extensible data notation.
On this lower layer an economic system is built - the proof-of-asset algorithm. The lower layer is created such that in principle any sensible economic incentive and consensus algorithm
can be plugged in and more generically a distributed system can be built. Horizon allows for arbitrary message encoding and signing, using new primitives for communication between nodes. This makes it more general as a transaction and communcation platform. 

## runing a node

install golang and git

run node:
```go build -o bin/horizon
./bin/horizon```

see [install docs](https://github.com/stackledger/docs/blob/master/install.md)
see also [telnet](https://github.com/stackledger/docs/blob/master/telnet.md)

## client functions

create keys

```cd client && go run client.go -option=createkeys```

 verify signature
 
 ```go run client.go -option=verify```


## testing

```go test ./...```

## architecture

* netio
* transaction multiplexing
* immutable datastructures on the language level

## contributions

contributions, such as pull requests, bug reports and comments are welcome


License: MIT license
