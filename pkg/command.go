package pkg

import (
	"context"
	"github.com/Shoothzj/cli/pkg/util"
	"github.com/segmentio/kafka-go"
	"golang.org/x/time/rate"
	"k8s.io/klog/v2"
	"strconv"
	"time"
)

func KafkaSend(host, topic string, port, size, tps int) {
	// make a writer that produces to topic-A, using the least-bytes distribution
	w := &kafka.Writer{
		Addr:     kafka.TCP(host + ":" + strconv.Itoa(port)),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	rate.NewLimiter(rate.Limit(tps), tps)
	// construct kafka producer to produce messages
	for true {
		headers := make([]kafka.Header, 1)
		headers[0] = kafka.Header{Key: KafkaShTime, Value: util.TimeToBytes(time.Now())}
		err := w.WriteMessages(context.Background(), kafka.Message{Headers: headers, Value: util.FixLengthReadableByte(size)})
		if err != nil {
			panic(err)
		}
	}
}

func KafkaConsume(host, topic, group string, port int) {
	klog.Infoln("init a consumer")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{host + ":" + strconv.Itoa(port)},
		GroupID:  group,
		Topic:    topic,
		MinBytes: 10 * 1024 * 1024,
		MaxBytes: 50 * 1024 * 1024,
	})
	klog.Infoln("begin to message loop")
	for {
		message, err := r.ReadMessage(context.Background())
		if err != nil {
			panic(err)
		}
		for _, header := range message.Headers {
			if header.Key == KafkaShTime {
				sendTime := util.BytesToTimeNoError(header.Value)
				klog.Infof("consume message, cost is %d", time.Since(sendTime).Milliseconds())
				break
			}
		}
	}

}
