package main

import (
	"assignmentapp/internal/api"
	"assignmentapp/internal/storage"
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

const (
	logFileName = "../server.log"
	natsSubject = "foo"
	nWorkers    = 4
)

var (
	natsClient   *nats.Conn
	actionLogger *log.Logger
	logFile      *os.File
	cache        *storage.Storage
	dataChan     chan *nats.Msg
	doneChan     chan struct{}
)

func init() {
	var err error
	log.SetLevel(log.DebugLevel)

	actionLogger = log.New()
	logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("fail to open log file: ", err)
	}
	actionLogger.SetOutput(logFile)

	natsClient, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("fail to nats connect: ", err)
	}

	cache = storage.New(100)
	doneChan = make(chan struct{}, 1)
	dataChan = make(chan *nats.Msg, 1)
}

func handleMsg(msg api.Msg) {
	switch msg.Action {
	case api.ActionAddItem:
		log.Debugf("Add item operation requested. Item: %+v", msg.Item)
		if err := cache.AddItem(msg.Item); err != nil {
			log.Error("Add item error: ", err)
		} else {
			log.Info("Add item success")
			actionLogger.Infof("AddItem %v", msg.Item)
		}
	case api.ActionDelItem:
		log.Debugf("Remove item operation requested. Key: %s", msg.Item.Key)
		if err := cache.RemoveItem(msg.Item.Key); err != nil {
			log.Error("Remove item error: ", err)
		} else {
			log.Info("Remove item success")
			actionLogger.Infof("RemoveItem %s", msg.Item.Key)
		}
	case api.ActionGetItem:
		log.Debugf("Get item operation requested. Key: %s", msg.Item.Key)
		if val, err := cache.GetItem(msg.Item.Key); err != nil {
			log.Error("Get item error: ", err)
		} else {
			log.Info("Get item success")
			actionLogger.Infof("GetItem %s: %s", msg.Item.Key, val)
		}
	case api.ActionGetAll:
		log.Debugf("Get all items operation requested")
		vals := cache.GetAllItems()
		log.Info("Get all items success")
		actionLogger.Infof("GetAllItems: %v", vals)

	default:
		log.Warnf("Unsupported action: %d", msg.Action)
	}
}

func worker(wg *sync.WaitGroup) {
	defer wg.Done()

	for msg := range dataChan {
		var decode api.Msg
		if err := json.Unmarshal(msg.Data, &decode); err != nil {
			log.Error("unmarshal msg")
			continue
		}
		handleMsg(decode)
	}
}

func subscriber() {
	sub, err := natsClient.ChanSubscribe(natsSubject, dataChan)
	if err != nil {
		log.Fatal("fail to nats subscribe: ", err)
	}
	log.Info("nats subscribe: ", natsSubject)

	<-doneChan

	sub.Unsubscribe() //lint: no err check
	log.Info("nats unsubscribe: ", natsSubject)
}

func main() {
	defer func() {
		if err := logFile.Close(); err != nil {
			println("unable close log file", err.Error())
		}
	}()
	defer natsClient.Close()

	log.Info("Server started")

	var wg sync.WaitGroup
	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go worker(&wg)
	}
	go subscriber()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	s := <-signalChan
	log.Info("got signal ", s.String())

	close(doneChan)
	close(dataChan)
	wg.Wait()
	log.Info("Terminated")
}
