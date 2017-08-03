package models

import (
	//	"encoding/json"
	"fmt"
	//	"log"

	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"

	"apt-web-server/modules/mlog"
)

var kafkaAddrs = []string{KAFKA + ":9092", KAFKA + ":9093"}

func SendOfflineMsg(bytes []byte) error {
	conf := kafka.NewBrokerConf("test-client")
	conf.AllowTopicCreation = true

	broker, err := kafka.Dial(kafkaAddrs, conf)
	if err != nil {
		mlog.Debug("cannot connect to kafka cluster")
	}
	defer broker.Close()

	producer := broker.Producer(kafka.NewProducerConf())

	msg := &proto.Message{Value: bytes}

	_, produceErr := producer.Produce(TOPIC, PARTITION, msg)
	if produceErr != nil {
		log := fmt.Sprintf("cannot produce message to %s:%d: %s", TOPIC, PARTITION, err)
		mlog.Debug(log)
	}

	return produceErr
}
