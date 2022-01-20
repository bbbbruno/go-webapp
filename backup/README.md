## ビルド

### `backup`コマンド

```bash
$ go build -o backup ./cmds/backup
```

### `backupd`コマンド

```bash
$ go build -o backupd ./cmds/backupd
```

## 使用方法

### バックアップ先ディレクトリを作成

```bash
$ mkdir backupdata
```

### バックアップ対象フォルダを追加

```bash
$ ./backup add ./test ./test2
+ ./test [まだアーカイブされてません]
+ ./test2 [まだアーカイブされてません]
```

### バックアップ対象フォルダを確認

```bash
$ ./backup list
{./test まだアーカイブされてません}
{./test2 まだアーカイブされてません}
```

### バックアップ対象フォルダを削除

```bash
$ ./backup remove
- {./test まだアーカイブされてません}
```

### バックアップデーモンを起動

```bash
$ ./backupd -db="./backupdata" -archive="./archive" -interval=5s
2022/01/20 14:14:18 チェックします...
2022/01/20 14:14:18  2個のディレクトリをアーカイブしました
2022/01/20 14:18:25 チェックします...
2022/01/20 14:18:25 変更はありません
```
