# NATS Streaming Server example

This example illustrates how to run a local cluster of
NATS streaming servers doing some pub/sub.

## Prerequisites

The example has been tested with NATS streaming server 0.15.1
and Go 1.11.6.

## Nutshell

Open 1st console:

```
$ nats-streaming-server -c a.conf
...
```

Open 2nd console:

```
$ nats-streaming-server -c b.conf
...
```

Open 3rd console:

```
$ nats-streaming-server -c c.conf
...
```

Open 4th console to compile and run a producer of messages
being sent to the cluster:

```
$ make
$ ./pub
...
```

Open 5th console to run a first consumer of messages:

```
$ ./sub
...
```

Open 6th console to run a second consumer of messages:

```
$ ./sub
...
```

As the consumers work as a consumer group, only one of the
consumers gets a message. This is in contrast to a normal
subscription model where each consumer gets a copy of each
message.

# Licenses

MIT
