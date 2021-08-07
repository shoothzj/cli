package kafka

import (
	"github.com/Shoothzj/cli/cmd/commands"
	"github.com/Shoothzj/cli/pkg"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var (
	// hosts
	host          string
	topic         string
	port          int
	size          int
	tps           int
	consumerGroup string
	kafkaCommand  = cobra.Command{
		Use: "kafka",
		Run: func(cmd *cobra.Command, args []string) {
			klog.Info("host is %s", host)
		},
	}
)

// init ssh 子命令初始化
func init() {
	kafkaCommand.PersistentFlags().StringVar(&host, "host", "", "kafka hosts")
	commands.RootCmd.AddCommand(&kafkaCommand)
	addProducerCommand()
	addConsumerCommand()
}

func addProducerCommand() {
	producerCommand := cobra.Command{
		Use: "produce",
		Run: func(cmd *cobra.Command, args []string) {
			pkg.KafkaSend(host, topic, port, size, tps)
		},
	}
	producerCommand.Flags().StringVar(&topic, "topic", "", "kafka topic")
	producerCommand.Flags().IntVar(&port, "port", 9092, "kafka port")
	producerCommand.Flags().IntVar(&size, "size", 1024, "message size")
	producerCommand.Flags().IntVar(&tps, "tps", 10, "tps")

	kafkaCommand.AddCommand(&producerCommand)
}

func addConsumerCommand() {
	consumerCommand := cobra.Command{
		Use: "consume",
		Run: func(cmd *cobra.Command, args []string) {
			pkg.KafkaConsume(host, topic, consumerGroup, port)
		},
	}
	consumerCommand.Flags().StringVar(&topic, "topic", "", "kafka topic")
	consumerCommand.Flags().IntVar(&port, "port", 9092, "kafka port")
	consumerCommand.Flags().StringVar(&consumerGroup, "group", "group", "consumer group")

	kafkaCommand.AddCommand(&consumerCommand)
}
