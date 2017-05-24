package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/streadway/amqp"
)

var (
	// server options
	host     string
	port     int
	vhost    string
	username string
	password string

	// exchange options
	exchange        string
	routingkey      string
	durableExchange bool
	exchangeType    string

	// queue options
	queue        string
	number       int
	passive      bool
	exclusive    bool
	noAck        bool
	durableQueue bool
)

var valid_properties = map[string]string{
	"content-type":     "ContentType",
	"content-encoding": "ContentEncoding",
	"delivery-mode":    "DeliveryMode",
	"priority":         "Priority",
	"correlation-id":   "CorrelationId",
	"reply-to":         "ReplyTo",
	"expiration":       "Expiration",
	"message-id":       "MessageId",
	"type":             "Type",
	"user-id":          "UserId",
	"app-id":           "AppId",
}

type table struct {
	amqp.Table
	isProperty bool
}

func (t *table) String() string {
	return ""
}
func (t *table) Type() string {
	return "string"
}
func (t *table) Set(str string) error {
	i := 0
	if i = strings.IndexRune(str, ':'); i < 0 {
		return errors.New(`unable to parse value, must be in "key:value" format`)
	}
	key, val := str[:i], str[i+1:]
	if t.isProperty && valid_properties[key] == "" {
		return errors.New(`invalid property: ` + key)
	}

	if t.isProperty {
		(*t).Table[valid_properties[key]] = val
	} else {
		(*t).Table[key] = val
	}
	return nil
}

var headers = table{amqp.Table{}, false}
var properties = table{amqp.Table{}, true}

var RootCmd = &cobra.Command{
	Use:               "amqptools",
	DisableAutoGenTag: true,
	Short:             "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:`,
}

func commonFlagSet() *pflag.FlagSet {
	var fs = pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringVarP(&host, "host", "H", "localhost", "specify host")
	fs.IntVarP(&port, "port", "P", 5672, "specify port")
	fs.StringVarP(&vhost, "vhost", "v", "/", "specify vhost")
	fs.StringVarP(&username, "username", "u", "guest", "specify username")
	fs.StringVarP(&password, "password", "p", "guest", "specify password")

	fs.StringVarP(&exchange, "exchange", "e", "", `exchange name (default "")`)
	fs.StringVarP(&routingkey, "key", "k", "", `routing key (default "")`)
	fs.StringVarP(&exchangeType, "type", "t", "direct", "exchange type")
	fs.BoolVarP(&durableExchange, "durable", "", false, "durable exchange")
	fs.SortFlags = false
	return fs
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
