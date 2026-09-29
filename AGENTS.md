# Room Finder AI Agent Instructions

このファイルは、共有ルールを複製せず、Room Finder固有の補足だけを定義する。

## 作業開始時に読むもの

次の順に確認する。

1. `../buildlog/AGENTS.md`
2. `../buildlog/README.md`
3. `https://github.com/kikudesuyo/dev-platform` の `README.md` と `dev-guideline/README.md`
4. このリポジトリの `README.md`、`SPEC.md`、関連ドキュメント、既存コード
5. 物件検索・Web fetch・スクレイピングを行う場合は [Room Finder物件検索skill](.agents/skills/room-finder-rental-search/SKILL.md)

参照先を読めない場合は、共有ルールや仕様を推測して実装しない。必要な判断をユーザーへ確認する。

Issue、branch、Acceptance / Quality / Delivery Gate、GitHub運用、一般的な実装・レビュー・納品手順は `../buildlog/AGENTS.md` と dev-platform のルールに従う。

## Room Finder固有の境界

- AgentやブラウザからPostgreSQLへ直接接続しない。物件情報・検索プロファイルの読み書きはGo APIへ集約する。
- Web/API/DBの構成、責務分離、エラーハンドリング、API形式、テスト方針は `README.md`、`SPEC.md`、既存の `api/` / `web/` 実装を正とする。
- 物件検索のフィルタ適用、詳細ページの根拠確認、候補・一致・保存件数の分離、バッチ再開、Go API保存、Web UI再表示は、上記のrepository-local skillに従う。
- skillが定める確認を実行できない場合、保存済み・UI反映済み・条件一致済みとは報告せず、未検証範囲と理由を明記する。
- ユーザーが「検索だけ」「DBへ保存しない」と明示した場合は、保存手順よりその指定を優先する。
