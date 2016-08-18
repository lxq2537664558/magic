# Mecury

Mecury is written in pure Go, target is  collecting metrics from the system it's
running on, or from other services, and writing them into InfluxDB、Nats or other
[outputs](https://github.com/corego/vgo/tree/master/mecury#supported-output-plugins).It is been 
developed  based on Telegraf: https://github.com/influxdata/telegraf ,but has changed
a lot,for more changes details ,please ses the head of main.go.

Features
--------------
    - Pure Go
    - Minimum Dependicies
    - Minimum system resources used(cpu and  memory、net bandwidth)
    - Fertile Plugins
    - Easy-Using Configurations
    - Once deployed,no more maintainance ,because of Configuration auto-update
    - For chinese users, we don't need any packages out of the wall
    - Extreme high performance
    - Vendoring Not third-party dependency manage

## Installation:
    ```bash
       $ go get -u github.com/corego/vgo
       $ cd $GOPATH/src/github.com/corego/vgo/mecury && go build
    ```

    Thats all,you even don't need to download any dependencies


## How to use it:
```bash
    $ ./mecury &
```

## How to close it
```bash
    $ pkill mecury
```

## Configuration

See the [configuration guide](docs/Configuration.md)

## Supported Input Plugins

See [input plugins](https://github.com/corego/vgo/tree/master/mecury/plugins/input)

## Supported Output Plugins
See [output plugins](https://github.com/corego/vgo/tree/master/mecury/plugins/output)

## Contributing
