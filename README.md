# amqptools

## Installing

```
go get -u github.com/hassansin/amqptools
```

## Usage 

## amqptools consume

Consumes messages

### Synopsis


Consume messages

By default, it runs forever and waits for any new message in thequeue. Pass `--number` to consume certain number of messages and quit

Use comma-separated values for binding the same queue with multiple routing keys:

	amqptools consume --exchange logs --keys info,warning,debug
	
	

```
amqptools consume [flags]
```

### Examples

```
  ampqtool consume -H ampq.example.com -P 5672 --exchange amq.direct --durable-queue

```

### Options

```
  -H, --host string       specify host (default "localhost")
  -P, --port int          specify port (default 5672)
  -v, --vhost string      specify vhost (default "/")
  -u, --username string   specify username (default "guest")
  -p, --password string   specify password (default "guest")
  -e, --exchange string   exchange name (default "")
  -k, --key string        routing key (default "")
  -t, --type string       exchange type (default "direct")
      --durable           durable exchange
  -q, --queue string      specify queue (default auto-generated)
  -n, --number int        retrieve maximum n messages. 0 = forever (default 0)
      --passive           passive queue
      --exclusive         exclusive queue
      --durable-queue     durable queue
      --no-ack            don't send ack
  -h, --help              help for consume
```

## amqptools publish

Publishes a message

### Synopsis

Publish a message using exchange and routing key.

If an argument is passed, that is used as message. Otherwise message is read from STDIN

	echo 'hello world' | amqptools publish --exchange=logs --key=info


To pass headers and properites, use `--headers` & `--properties` any number of times in `key:value` format


```
amqptools publish [flags] [message]
```

### Examples

```
  ampqtools publish -H ampq.example.com -P 5672 --exchange=amq.direct --key=hello "hello world"
  amqptools publish "hello world" --properties="content-type:text/html" --properties="expiration:3000"	
	
```

### Options

```
  -H, --host string         specify host (default "localhost")
  -P, --port int            specify port (default 5672)
  -v, --vhost string        specify vhost (default "/")
  -u, --username string     specify username (default "guest")
  -p, --password string     specify password (default "guest")
  -e, --exchange string     exchange name (default "")
  -k, --key string          routing key (default "")
  -t, --type string         exchange type (default "direct")
      --durable             durable exchange
      --properties string   message properties, key:value format
      --headers string      message headers, key:value format
  -h, --help                help for publish
```

