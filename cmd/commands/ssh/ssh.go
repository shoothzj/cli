package ssh

import (
	"github.com/Shoothzj/cli/cmd/commands"
	"github.com/paashzj/gl/ssh"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

// hosts 地址列表，以逗号分隔
// user 用户名
// password 密码
// port 端口
// command 命令
var (
	hosts    []string
	user     string
	password string
	port     uint
	command  string
)

// init ssh 子命令初始化
func init() {
	sshCommand := cobra.Command{
		Use: "ssh",
		Run: func(cmd *cobra.Command, args []string) {
			klog.Infof("host is %s user is %s password is %s port is %", hosts, user, password, port)
			klog.Infof("execute ssh command %s", command)
			for _, host := range hosts {
				client, err := ssh.NewPasswordClient(host, user, password, port)
				if err != nil {
					klog.Errorf("create host %s ssh client error %s", host, err)
					continue
				}
				bytes, err := client.Run(command)
				if err != nil {
					klog.Errorf("execute %s command error %s", host, err)
					continue
				}
				klog.Infof("response from %s is %s", host, string(bytes))
			}
		},
	}
	sshCommand.Flags().StringSliceVar(&hosts, "hosts", nil, "ssh hosts separated by comma")
	sshCommand.Flags().StringVar(&user, "user", "", "ssh user")
	sshCommand.Flags().StringVar(&password, "password", "", "ssh password")
	sshCommand.Flags().UintVar(&port, "port", 22, "ssh port")
	sshCommand.Flags().StringVar(&command, "commands", "", "ssh commands")
	commands.RootCmd.AddCommand(&sshCommand)
}
