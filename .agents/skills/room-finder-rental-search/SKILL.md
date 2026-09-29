---
name: room-finder-rental-search
description: "Room Finderで複数の掲載元から賃貸物件を検索・詳細確認・保存する定型業務を実行する。"
---

# Room Finder 物件検索

このスキルは、Room Finderで賃貸物件を取得・判定・保存するときに使用する。共通処理と掲載元ごとの取得処理を分離し、このリポジトリの `AGENTS.md`、`SPEC.md`、`api/crawler`、`api/runner`、`api/agent` の実装を正とする。

## 必須フロー

1. 掲載元を選び、その掲載元の source adapter の検索URL、対象ホスト、robots、取得間隔、候補抽出方法を確認する。
2. ユーザー条件を、掲載元の検索フィルタで適用する条件と詳細ページで確認する条件に分ける。
3. フィルタを設定して候補URLを取得する。掲載元別の候補件数を記録する。
4. 詳細ページで各必須条件を明示的な根拠テキスト付きで確認する。不明・掲載終了・空室状況不明は不一致として扱う。
5. 条件をすべて満たした物件だけを、共通の構造化形式でGo APIへ送信する。DBへ直接接続しない。
6. APIレスポンスで保存を確認し、Web UIで同じ物件を再表示できることを確認する。
7. 掲載元別に候補件数、一致件数、保存件数、失敗件数、未処理件数を分けて報告する。条件を自動で緩和しない。

掲載元ごとの差分と adapter の契約は [source-adapter-contract.md](references/source-adapter-contract.md) を読む。既存のLIFULL HOME'S対応は `api/crawler/lifullhomes` とそのテストを参照し、新しい掲載元を追加する場合も共通の保存・判定責務へ直接サイト固有処理を混ぜない。

## バッチ実行

候補URLの一覧を1行1URLで用意し、10〜20件のバッチに分けて、最後に成功したURLをcheckpointへ保存する。既存のrunnerが保存する状態と混同しないよう、用途ごとにstateファイルを分ける。

```bash
.agents/skills/room-finder-rental-search/scripts/batch-fetch.sh \
  --source lifullhomes \
  --input candidates.txt \
  --state .local/rental-search/lifullhomes/state \
  --batch-size 20 \
  -- 'go run ./cmd/rental-fetch --source lifullhomes --url {item}'
```

コマンドの終了コードは `0=成功`、`10=一時失敗（再試行対象）`、`20=不一致・確認不能（記録して次へ）` とする。その他の終了コードは実行エラーとしてその場で停止し、checkpointより後を未処理のままにする。出力やログに秘密情報を含めない。

## 境界

- robots.txt、掲載元の対象ホスト制限、リクエスト間隔、取得元URL・取得日時・根拠保存を守る。
- 詳細ページで根拠が取れない条件を推測で合格扱いにしない。
- shell scriptは候補の意味判定やDB保存を代替しない。判定と保存は既存のGo crawler/agent/APIの責務に従う。
- 外部サイトへの長時間処理は、10〜20件のバッチ、完了済み件数、再開位置を記録し、60秒以内の処理単位で進捗を確認する。
