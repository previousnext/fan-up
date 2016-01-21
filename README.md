Fan Up
======

A daemon for ensuring the fan network bridge is up.

https://wiki.ubuntu.com/FanNetworking

## Usage

Basic

```bash
$ fan-up
```

Advanced

```bash
$ fan-up --interface=eth1 --overlay=10.0.0.0
```

## Building

We use a tool called gb. To install run:

```bash
$go get github.com/constabulary/gb/...
```

To build the project run:

```bash
$ make
```

