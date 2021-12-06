package cmd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/oleiade/reflections"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

// consumeCmd represents the consume command
var consumeCmd = &cobra.Command{
	Use:     "consume [flags]",
	Aliases: []string{"receive"},
	Short:   "Consumes messages",
	Long: `Consume messages
By default, it runs forever and waits for any new message in thequeue. Pass '--number' to consume certain number of messages and quit

Use comma-separated values for binding the same queue with multiple routing keys:

	amqptools consume --exchange logs --keys info,warning,debug
	
	`,
	Example: `  ampqtool consume -H ampq.example.com -P 5672 --exchange amq.direct --durable-queue
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("unknown arg: %s", args[0])
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		uri := getUri()
		// Dial amqp server
		conn, err := amqp.Dial(uri)
		if err != nil {
			return fmt.Errorf("connection.open: %v", err)
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			return fmt.Errorf("channel.open: %v", err)
		}
		defer ch.Close()

		closes := ch.NotifyClose(make(chan *amqp.Error, 1))

		q, err := ch.QueueDeclare(
			queue,         // queue
			durableQueue,  // durable
			!durableQueue, // auto-delete
			exclusive,     // exclusive
			false,         // no-wait
			nil,           // args
		)
		if err != nil {
			return fmt.Errorf("queue.declare: %v", err)
		}
		// bind queue for non-default exchange
		if exchange != "" {
			if err = ch.ExchangeDeclare(
				exchange,        // name
				exchangeType,    // type
				durableExchange, // durable
				false,           // auto-delete
				false,           // internal
				false,           // no-wait
				nil,             // args
			); err != nil {
				return fmt.Errorf("exchange.declare: %v", err)
			}
			// bind the queue to all routingkeys
			for _, key := range strings.Split(routingkey, ",") {
				if err = ch.QueueBind(
					queue,    // queue
					key,      // routing-key
					exchange, //exchange
					false,    // no-wait
					nil,      // args
				); err != nil {
					return fmt.Errorf("queue.bind: %v", err)
				}
			}
		}

		prefetchCount := 10
		if number > 0 {
			prefetchCount = number
		}
		if err = ch.Qos(
			prefetchCount, // prefetch count
			0,             // prefetch size
			false,         // global
		); err != nil {
			return fmt.Errorf("basic.qos: %v", err)
		}

		msgs, err := ch.Consume(
			queue,     // queue
			"",        // consumer id
			false,    // auto-ack
			exclusive, // exclusive
			false,     // no-local
			false,     // no-wait
			nil,       // args
		)
		if err != nil {
			return fmt.Errorf("basic.consume: %v", err)
		}
		count := 0

		go func() {
			for msg := range msgs {
				fmt.Printf("Timestamp: %s\n", msg.Timestamp)
				fmt.Printf("Exchange: %s\n", msg.Exchange)
				fmt.Printf("RoutingKey: %s\n", msg.RoutingKey)
				fmt.Printf("Queue: %s\n", q.Name)
				fmt.Printf("Redelivered: %v\n", msg.Redelivered)
				fmt.Printf("Headers: %v\n", msg.Headers)
				fmt.Printf("ConsumerTag: %v\n", msg.ConsumerTag)
				fmt.Printf("Payload: \n%s\n\n", msg.Body)

				fmt.Println("Properties:")
				for _, key := range valid_properties {
					if val, err := reflections.GetField(msg, key); err == nil && val != reflect.Zero(reflect.TypeOf(val)).Interface() {
						fmt.Printf("%s: %q\n", key, val)
					}
				}

				if !noAck {
					ch.Ack(msg.DeliveryTag, false)
				}

				count++
				if number != 0 && count >= number {
					ch.Close()
				}
			}
		}()
		fmt.Printf("Waiting for messages. Queue: %s. To exit press CTRL+C\n\n", q.Name)
		<-closes
		fmt.Println("Connection Closed")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(consumeCmd)
	consumeCmd.Flags().AddFlagSet(commonFlagSet())
	consumeCmd.Flags().StringVarP(&queue, "queue", "q", "", "specify queue (default auto-generated)")
	consumeCmd.Flags().IntVarP(&number, "number", "n", 0, "retrieve maximum n messages. 0 = forever (default 0)")
	consumeCmd.Flags().BoolVarP(&passive, "passive", "", false, "passive queue")
	consumeCmd.Flags().BoolVarP(&exclusive, "exclusive", "", false, "exclusive queue")
	consumeCmd.Flags().BoolVarP(&durableQueue, "durable-queue", "", false, "durable queue")
	consumeCmd.Flags().BoolVarP(&noAck, "no-ack", "", false, "don't send ack")
	consumeCmd.Flags().SortFlags = false
}
