# 分散型 Twitter 投票アプリケーション

## 起動方法

### `nsqlookupd`

```bash
$ nsqlookupd
```

### `nsqd`

```bash
$ nsqd --lookupd-tcp-address=localhost:4160
```

### `mongod`

```bash
$ mongod --dbpath ./db
```

### `main.go`

```bash
$ cd twittervotes
$ go build
$ ./twittervotes
```
