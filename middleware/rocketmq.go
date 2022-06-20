package middleware

import (
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var p rocketmq.Producer

func InitProducer() {
	var err error
	p, err = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"10.214.150.171:9876"})),
		producer.WithRetry(2),
		producer.WithNamespace("MQ_INST_8149062485579066312_2586445845"),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: "rocketmq2",
			SecretKey: "12345678",
		}),
		producer.WithGroupName("GIT_Group"),
	)
	if err != nil {
		panic(err)
	}
	err = p.Start()
	if err != nil {
		panic(err)
	}
	fmt.Println("Producer started successfully!")
}

func GetProducer() rocketmq.Producer {
	return p
}

func CloseProducer() {
	err := p.Shutdown()
	if err != nil {
		panic(err)
	}
}
