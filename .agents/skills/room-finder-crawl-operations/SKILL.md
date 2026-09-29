---
name: room-finder-crawl-operations
description: "Room Finderの外部データ取得をbatch・checkpoint・再試行付きで安全に実行する。"
---

# Room Finder Crawl Operations

物件候補を掲載元ごとに分割実行するときに使用する。掲載元の検索条件や物件の意味判定は [Room Finder rental search skill](../room-finder-rental-search/SKILL.md) に従い、このスキルは実行単位の制御だけを扱う。

## 実行

候補URLを1行1件で用意し、10〜20件単位で実行する。sourceごとにstateを分ける。

```bash
.agents/skills/room-finder-crawl-operations/scripts/batch-fetch.sh \
  --source lifullhomes \
  --input candidates.txt \
  --state .local/rental-search/lifullhomes/state \
  --batch-size 20 \
  -- 'go run ./cmd/rental-fetch --source lifullhomes --url {item}'
```

コマンドの終了コードは `0=成功`、`10=一時失敗（再試行対象）`、`20=不一致・確認不能（処理済み）` とする。その他の終了コードは実行エラーとして停止する。

## 不変条件

- 成功または不一致として完了した候補だけをcheckpointへ追加する。
- 終了コード10の候補は次回実行へ残す。
- その他のエラーでは後続候補を処理せず、未処理位置を出力する。
- 実行ごとにsource、処理数、成功数、再試行数、不一致数、未処理位置を記録する。
- 外部サイトへのアクセスは掲載元adapterのrobots、対象ホスト、リクエスト間隔、タイムアウトに従う。
- 60秒以内のboundedな処理単位に分割し、長時間実行中は少なくとも60秒ごとに進捗を確認する。

このscriptはローカル実行・fixture検証・手動再開用であり、本番SchedulerやDB保存の代替ではない。
