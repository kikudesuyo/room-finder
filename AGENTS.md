# AI Agent Instructions

## 最重要ルール

このプロジェクトは、`../buildlog` と同じアーキテクチャ、設計方針、実装規則、検証手順、GitHub運用を採用する。

作業開始前に必ず次を読むこと。

1. `../buildlog/AGENTS.md`
2. `../buildlog/README.md`
3. `https://github.com/kikudesuyo/dev-platform` の `README.md`
4. `https://github.com/kikudesuyo/dev-platform` の `dev-guideline/README.md`
5. このリポジトリの `README.md`、関連ドキュメント、既存コード

`../buildlog/AGENTS.md` とこのファイルの指示は、一般的な開発慣行より優先する。buildlogのルールと異なる独自方式を持ち込んではならない。

参照先が読めない場合、推測で作業を進めず、必要な判断をユーザーへ確認すること。

## アーキテクチャ

以下をbuildlogと同一にする。

- フロントエンド：SvelteKit、TypeScript、Tailwind CSS
- フロントエンドホスティング：Cloudflare Pages
- バックエンド：Go、chi、GORM
- バックエンドホスティング：GCP Cloud Run
- データベース：Neon PostgreSQL
- シークレット管理：GCP Secret Manager
- CI/CD：GitHub Actions + GCP Workload Identity Federation
- マイグレーション：Atlas
- ローカル開発DB：Docker ComposeのPostgreSQL

レイヤー構成、責務分離、命名、エラーハンドリング、APIレスポンス形式、テスト方針はbuildlogの既存実装を確認して踏襲すること。

AgentやブラウザからDBへ直接接続してはならない。DB操作はGo APIに集約すること。

## 実装ルール

- Issueを起点に実装する
- Issueごとに `issue/{issue_number}` ブランチを作る
- Issueの範囲を超える変更を行わない
- 1機能1ファイル原則と既存のレイヤー責務を守る
- 将来のためだけの抽象化、予備カラム、不要なモデルを追加しない
- 既存実装と異なる方式を採用する場合は、理由をIssueまたはPRに記録する
- 依存関係更新や大規模リファクタリングは別Issueに分ける
- 他Agentやユーザーの変更を勝手に削除・リセットしない
- 秘密情報をコード、ログ、コミットへ含めない

## Acceptance / Quality / Delivery Gate

実装完了は、コードを書いた時点ではなく、以下をすべて満たした時点とする。

1. IssueのAcceptance Criteriaを実際の挙動・出力で検証し、すべてPASSである
2. format、lint、type check、test、build、`git diff --check`を変更内容に応じて実行する
3. 差分、セキュリティ、影響範囲を自己レビューする
4. commit、push、Pull Request作成まで完了する
5. PRにはIssueリンク、Acceptance Criteria、設計判断、検証結果、既知の制約を記載する

UI変更時は、buildlogのルールに従い、実ブラウザでMobileとPCを確認し、スクリーンショット等の証跡をPRから参照可能にする。

検証できない項目がある場合、完了扱いにせず、未検証の理由と範囲を報告すること。

## 強制プロンプト

このプロジェクトで作業するAI Agentは、常に次の指示に従うこと。

> buildlogと同じアーキテクチャ・設計・コード規則・命名規則・レイヤー構成・エラーハンドリング・検証手順・Issue/PR運用を厳守せよ。既存のbuildlog実装を確認せずに独自設計を作るな。Issueの範囲を超える変更をするな。実装だけで完了と判断するな。Acceptance Criteria、Quality Gate、Delivery Gateをすべて満たし、実際の挙動を検証してから完了報告せよ。
