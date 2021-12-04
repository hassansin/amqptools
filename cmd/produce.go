package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"time"

	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

// produceCmd represents the produce command
var produceCmd = &cobra.Command{
	Use:     "publish [flags] [message]",
	Aliases: []string{"produce", "send"},
	Short:   "Publishes a message",
	Long: `Publish a message using exchange and routing key.
mesage can be string or stdin:

	echo 'hello world' | amqptools publish --exchange=logs --key=info

To pass headers and properites, use '--headers' & '--properties' any number of times in 'key:value' format

`,
	Example: `  ampqtools publish -H ampq.example.com -P 5672 --exchange=amq.direct --key=hello "hello world"
  amqptools publish "hello world" --properties="content-type:text/html" --properties="expiration:3000"	
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		var message string
		if len(args) > 0 {
			message = args[0]
		} else {
			bytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("read file: %v", err)
			}
			message = string(bytes)
		}

		uri := getUri()
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
		}

		msg := amqp.Publishing{
			Headers:   headers.Table,
			Timestamp: time.Now(),
			Body:      []byte(message),
		}

		for key, val := range properties.Table {
			v := reflect.ValueOf(&msg).Elem().FieldByName(key)
			if v.IsValid() {
				v.SetString(val.(string))
			}
		}

		if err := ch.Confirm(false); err != nil {
			return fmt.Errorf("confirm.select destination: %s", err)
		}

		returns := ch.NotifyReturn(make(chan amqp.Return, 1))
		confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

		fmt.Printf("Sending message with Exchange: '%s' and Key: '%s'\n", exchange, routingkey)
		if err = ch.Publish(
			exchange,   // exchange
			routingkey, // routing-key
			true,       // mandatory
			false,      // immediate
			msg,        // message
		); err != nil {
			return fmt.Errorf("basic.publish: %v", err)
		}
		select {
		case r := <-returns:
			fmt.Println("Message not delivered: " + r.ReplyText)
		case c := <-confirms:
			if c.Ack {
				fmt.Println("Message published")
			} else {
				fmt.Println("Message failed to publish")
			}
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(produceCmd)
	produceCmd.Flags().AddFlagSet(commonFlagSet())
	produceCmd.Flags().VarP(&properties, "properties", "", "message properties, key:value format")
	produceCmd.Flags().VarP(&headers, "headers", "", "message headers, key:value format")
	produceCmd.Flags().SortFlags = false
}
