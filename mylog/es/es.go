package es

import (
	"context"
	"fmt"
	"strings"
	"test/logagent/mylog/config"
	"time"

	"github.com/olivere/elastic/v7"
)

var (
	client *elastic.Client
	ESchan chan *config.LogData
)

func Init(address string, size int) (err error) {
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		fmt.Println("Init ES failed,err:", err)
		return
	}
	fmt.Println("connect to es success")
	ESchan = make(chan *config.LogData, size)
	go SendToES()
	return
}

func SendToESChan(msg *config.LogData) {
	ESchan <- msg
}

func SendToES() {
	for {
		select {
		case msg := <-ESchan:
			put1, err := client.Index().Index(msg.Topic).BodyJson(msg).Do(context.Background())
			if err != nil {
				panic(err)
			}
			fmt.Printf("Index user:%s to index %s,type:%s\n", put1.Id, put1.Index, put1.Type)
		default:
			time.Sleep(time.Second)
		}
	}
}
