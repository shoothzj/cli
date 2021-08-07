package main

import (
	"flag"
	"github.com/Shoothzj/cli/cmd/commands"
	_ "github.com/Shoothzj/cli/cmd/commands/docker"
	_ "github.com/Shoothzj/cli/cmd/commands/kafka"
	_ "github.com/Shoothzj/cli/cmd/commands/scp"
	_ "github.com/Shoothzj/cli/cmd/commands/ssh"
	"k8s.io/klog/v2"
	"os"
)

func main() {
	flagSet := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(flagSet)
	klog.SetOutput(os.Stdout)
	commands.RootCmd.Execute()
}
