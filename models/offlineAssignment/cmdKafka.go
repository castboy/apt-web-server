package offlineAssignment

import (
	"apt-web-server_v2/modules/mlog"
	"fmt"

	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
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

	_, produceErr := producer.Produce(TOPIC, int32(PARTITION), msg)
	if produceErr != nil {
		log := fmt.Sprintf("cannot produce message to %s:%d: %s", TOPIC, PARTITION, err)
		mlog.Debug(log)
	}

	return produceErr
}
