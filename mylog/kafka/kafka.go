package kafka

import (
	"fmt"
	"test/logagent/mylog/config"
	"test/logagent/mylog/es"
	"time"

	"github.com/Shopify/sarama"
)

var (
	client      sarama.SyncProducer
	logDataChan chan *config.LogData
	consumer    sarama.Consumer
	pc          sarama.PartitionConsumer
)

func Init(address []string, max_size int) (err error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.Return.Successes = true

	client, err = sarama.NewSyncProducer(address, cfg)
	if err != nil {
		fmt.Println("Produce error:", err)
		return
	}

	logDataChan = make(chan *config.LogData, max_size)

	consumer, err = sarama.NewConsumer(address, nil)
	if err != nil {
		fmt.Println("Init consumer failed,err:", err)
		return
	}

	go SendToMessage()

	return
}

func SendToChan(topic, data string) {
	var t = &config.LogData{
		Topic: topic,
		Data:  data,
	}
	logDataChan <- t
}

func SendToMessage() {
	for {
		select {
		case ld := <-logDataChan:
			msg := sarama.ProducerMessage{}
			msg.Topic = ld.Topic
			msg.Value = sarama.StringEncoder(ld.Data)
			pid, offset, err := client.SendMessage(&msg)
			if err != nil {
				fmt.Println("Send Message error:", err)
			}
			fmt.Printf("pid:%v offser:%v Topic:%v Value:%v\n", pid, offset, ld.Topic, ld.Data)
		default:
			time.Sleep(time.Second)
		}
	}
}

func Consumer(topic string) {
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Println("Get partitions failed,err:", err)
		return
	}
	for partition := range partitionList {
		pc, err = consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Println("failed to start consumer for partition,err:", err)
			return
		}
	}
	defer pc.AsyncClose()
	go func(sarama.PartitionConsumer) {
		for msg := range pc.Messages() {
			fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
			t := &config.LogData{
				Topic: topic,
				Data:  string(msg.Value),
			}
			es.SendToESChan(t)
		}
	}(pc)
	select {}
}

// func Consumer(topic string) (err error) {
// 	partitionList, err := consumer.Partitions(topic) //根据topic取到所有的分区
// 	if err != nil {
// 		fmt.Println("fail to get list of partition:", err)
// 		return
// 	}
// 	var pc sarama.PartitionConsumer
// 	fmt.Println(partitionList)
// 	for partition := range partitionList {
// 		pc, err = consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
// 		if err != nil {
// 			fmt.Printf("failed to start consumer for partition %d,err:%v", partition, err)
// 			return
// 		}
// 		defer pc.AsyncClose()
// 		//异步从每个分区消费消息
// 		go func(sarama.PartitionConsumer) {
// 			for msg := range pc.Messages() {
// 				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
// 				//直接发给ES
// 				var ld = config.LogData{
// 					Topic: topic,
// 					Data:  string(msg.Value),
// 				}
// 				es.SendToESChan(&ld) //函数调函数
// 				//优化一下，直接放到chann中
// 			}
// 		}(pc)
// 		select {}
// 	}
// 	return
// }
