## amqptools

A brief description of your application

### Synopsis


A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

### Options

```
  -h, --help   help for amqptools
```

### SEE ALSO
* [amqptools consume](#amqptools-consume)	 - Consumes messages
* [amqptools publish](#amqptools-publish)	 - Publishes a message

## amqptools consume

Consumes messages

### Synopsis


Consume messages
Uses the default exchange '', When no exchange is provided
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

### SEE ALSO
* [amqptools](amqptools.md)	 - A brief description of your application

## amqptools publish

Publishes a message

### Synopsis


Publish a message using exchange and routing key.
mesage can be string or stdin:
  echo 'hello world' | amqptools publish --exchange=logs --key=info



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

### SEE ALSO
* [amqptools](amqptools.md)	 - A brief description of your application

