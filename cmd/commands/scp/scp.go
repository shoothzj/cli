package scp

import (
	"github.com/Shoothzj/cli/cmd/commands"
	"github.com/paashzj/gl/ssh"
	"github.com/pkg/sftp"
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"k8s.io/klog/v2"
	"os"
	"path/filepath"
	"strings"
)

var (
	// hosts 地址列表，以逗号分隔
	hosts []string
	// user 用户名
	user string
	// password 密码
	password string
	// port 端口
	port uint
	// dist 目标地址前缀
	dist string
	// content 内容
	content string
	// file 文件(夹)
	file string
	// excludeDirs 排除路径s
	excludeDirs []string
)

// init scp 子命令初始化
func init() {
	scpCommand := cobra.Command{
		Use: "scp",
		Run: func(cmd *cobra.Command, args []string) {
			klog.Infof("hosts is %s user is %s password is %s port is %", hosts, user, password, port)
			if len(content) > 0 && len(file) > 0 {
				klog.Errorf("You can't specify content with file.")
				return
			}
			klog.Infof("execute scp command")
			for _, host := range hosts {
				client, err := ssh.NewPasswordClient(host, user, password, port)
				if err != nil {
					klog.Errorf("create host %s scp client error %s", host, err)
					continue
				}
				sftpCli, err := client.NewSftp()
				if err != nil {
					klog.Errorf("create host %s sftp client error %s", host, err)
					continue
				}
				if len(content) > 0 {
					scpContent(sftpCli, host)
					continue
				}
				if len(file) > 0 {
					scpPosix(sftpCli, host)
				}
			}
		},
	}
	scpCommand.Flags().StringSliceVar(&hosts, "hosts", nil, "scp hosts separated by comma")
	scpCommand.Flags().StringVar(&user, "user", "", "scp user")
	scpCommand.Flags().StringVar(&password, "password", "", "scp password")
	scpCommand.Flags().UintVar(&port, "port", 22, "scp port")
	scpCommand.Flags().StringVar(&dist, "dist", "", "dist prefix")
	scpCommand.Flags().StringVar(&content, "content", "", "content")
	scpCommand.Flags().StringVar(&file, "file", "", "scp commands")
	scpCommand.Flags().StringSliceVar(&excludeDirs, "excludeDirs", nil, "excludeDirs separated by comma")
	commands.RootCmd.AddCommand(&scpCommand)
}

func scpContent(cli *sftp.Client, host string) {
	file, err := cli.Create(dist)
	if err != nil {
		klog.Errorf("create file error on host %s", host)
		return
	}
	defer file.Close()
	_, err = file.Write([]byte(content))
	if err != nil {
		klog.Errorf("write file error on host %s", host)
	}
}

func scpPosix(cli *sftp.Client, host string) {
	stat, err := os.Stat(file)
	if err != nil {
		klog.Errorf("local file not exists")
		return
	}
	if stat.IsDir() {
		scpDir(cli, host)
	} else {
		scpFile(cli, host, file, dist)
	}
}

func scpDir(cli *sftp.Client, host string) {
	klog.Info("begin to scp directory %s host %s", file, host)
	err := cli.Mkdir(dist)
	if err != nil {
		klog.Info("mkdir dir failed")
	}
	filepath.Walk(file, func(path string, info fs.FileInfo, err error) error {
		if len(path) == len(file) {
			return nil
		}
		relativePath := path[len(file)+1:]
		for _, excludeDir := range excludeDirs {
			if strings.HasPrefix(relativePath, excludeDir+"/") {
				return nil
			}
		}
		if info.IsDir() {
			err := cli.Mkdir(dist+path[len(file):])
			if err != nil {
				klog.Info("mkdir dir failed")
			}
		}
		if !info.IsDir() {
			scpFile(cli, host, path, dist+path[len(file):])
		}
		return nil
	})
}

func scpFile(cli *sftp.Client, host, srcFile, distFile string) {
	dist, err := cli.Create(distFile)
	if err != nil {
		klog.Errorf("create dist file %s error on host %s error %s", distFile, host, err)
		return
	}
	defer dist.Close()
	readFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		klog.Errorf("read local file %s error", srcFile)
		return
	}
	_, err = dist.Write(readFile)
	if err != nil {
		klog.Errorf("remote write file %s error on host %s", distFile, host)
		return
	}
}
