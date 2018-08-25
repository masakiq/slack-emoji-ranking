# Slack emoji ranking

## get token of Slack

* [Slack API 推奨Tokenについて - Qiita](https://qiita.com/ykhirao/items/3b19ee6a1458cfb4ba21)

* Token need following scope

```
channels:read
chat:write:user
reactions:read
```

## set environment

```
$ export SLACK_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
$ export SLACK_CHANNEL=general
```

## execute

```
$ go run main.go
```

## go version

```
$ go version
go version go1.10.2 linux/amd64
```
