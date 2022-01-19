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
$ mongo
> use ballots
> db.polls.insert({"title":"調査のテスト１","options":["one","two","three"]})
> db.polls.find().pretty()
```

### `twittervote`

```bash
$ cd twittervotes
$ go build
$ ./twittervotes
```

### `counter`

```bash
$ cd counter
$ go build
$ ./counter
```

### `api`

```bash
$ cd api
$ go build
$ ./api
```
