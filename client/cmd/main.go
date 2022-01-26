package main

import (
	"client/internal/api"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

const natsSubject = "foo"

var natsClient *nats.Conn

func init() {
	var err error
	natsClient, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("nats connect error: ", err)
	}
}

func addItem() {
	var key, value string
	fmt.Println("Enter key: ")
	fmt.Scanf("%s", &key) //lint: no err check
	fmt.Println("Enter value: ")
	fmt.Scanf("%s", &value) //lint: no err check

	bytes, _ := json.Marshal(api.Msg{
		Action: api.ActionAddItem,
		Item: api.Item{
			Key:   api.ItemKey(key),
			Value: api.ItemValue(value),
		},
	})

	if err := natsClient.Publish(natsSubject, bytes); err != nil {
		log.Println("nats publish:", err)
	}
}

func delItem() {
	var key string
	fmt.Println("Enter item key to delete: ")
	fmt.Scanf("%s", &key) //lint: no err check

	bytes, _ := json.Marshal(api.Msg{
		Action: api.ActionDelItem,
		Item:   api.Item{Key: api.ItemKey(key)},
	})

	if err := natsClient.Publish(natsSubject, bytes); err != nil {
		log.Println("nats publish:", err)
	}
}

func getItem() {
	var key string
	fmt.Println("Enter item key to get: ")
	fmt.Scanf("%s", &key) //lint: no err check

	bytes, _ := json.Marshal(api.Msg{
		Action: api.ActionGetItem,
		Item:   api.Item{Key: api.ItemKey(key)},
	})

	if err := natsClient.Publish(natsSubject, bytes); err != nil {
		log.Println("nats publish:", err)
	}
}

func getAll() {
	bytes, _ := json.Marshal(api.Msg{
		Action: api.ActionGetAll,
	})

	if err := natsClient.Publish(natsSubject, bytes); err != nil {
		log.Println("nats publish:", err)
	}
}

func handle(action byte) {
	switch action {
	case api.ActionAddItem:
		addItem()
	case api.ActionDelItem:
		delItem()
	case api.ActionGetItem:
		getItem()
	case api.ActionGetAll:
		getAll()
	default:
		fmt.Println("unsupported action: ", action)
	}
}

func main() {
	defer natsClient.Close()

	var action byte
	for {
		fmt.Println(`Enter action:
   1 - AddItem
   2 - RemoveItem
   3 - GetItem
   4 - GetAllItems
   0 - for exit`,
		)
		fmt.Scanf("%d", &action) //lint: no err check
		if action == 0 {
			break
		}
		handle(action)
	}
}
