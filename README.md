# cli

# 命令介绍

## kafka

- Kafka发送时候携带头部sh-time

### produce

```bash
./cli kafka --host "localhost" produce --topic test-1 --size 10248000 --tps 1
```

### consume

```bash
./cli kafka --host "localhost" consume --topic test-1
```

## scp
```bash
./cli scp --hosts 127.0.0.1 --user root --password pwd --port 22 --dist /home/sh/testfile --content "I am Magneto"
```

## ssh

```bash
./cli ssh --hosts 127.0.0.1 --user root --password pwd --port 22 --commands "ls"
```

# glog 级别
- 10 every detail
- 4 info  
- 3 warn  
- 2 error  
- 1 fatal
