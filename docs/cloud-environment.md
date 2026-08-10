# Cloud環境設定

このドキュメントはIssue #12の設定手順を定義する。実際のID・URL・秘密値はリポジトリへ保存しない。

## GCP

作成するリソース:

- Artifact Registry repository: `room-finder`
- Cloud Run service: `room-finder-api`
- Secret Manager secret: `room-finder-database-url`
- GitHub Actions用Workload Identity Federation provider
- Cloud Build実行用サービスアカウント

Cloud Buildのサービスアカウントには、Artifact Registryへのpush、Cloud Runのデプロイ、対象Secretの参照権限だけを付与する。Cloud Runは`DATABASE_URL`だけをSecret Managerから受け取り、その他の公開設定は環境変数で渡す。

GitHubリポジトリ変数:

```text
GCP_PROJECT_ID
GCP_REGION=asia-northeast1
GCP_WIF_PROVIDER
GCP_SERVICE_ACCOUNT
ARTIFACT_REPOSITORY=room-finder
CLOUD_RUN_SERVICE=room-finder-api
DATABASE_SECRET_NAME=room-finder-database-url
ALLOWED_ORIGINS=https://<pages-domain>
```

実行順序:

1. `migrate-production.yml`でAtlas migrationを適用する
2. `deploy-api.yml`でCloud Runへデプロイする
3. Cloud RunのURLからhealth endpointを確認する

## Neon PostgreSQL

Neonでproduction databaseを作成し、接続URLを`room-finder-database-url`へ登録する。接続URLはGitHub Actions、Cloud Run、リポジトリへ直接記録しない。

## Cloudflare Pages

Workers & PagesからGitHub repositoryを接続する。SvelteKitのFramework presetを使用し、次を設定する。

```text
Production branch: main
Build command: pnpm run build
Build directory: .svelte-kit/cloudflare
Root directory: web
```

Pagesの環境変数:

```text
PUBLIC_API_BASE_URL=https://<cloud-run-api-domain>
```

Preview環境ではPreview用API URLを設定する。API側の`ALLOWED_ORIGINS`にはproductionとpreviewの必要なOriginだけを列挙する。

## 実行できない項目

GCP project、Neon database、Cloudflare account、GitHub Actions variablesはアカウント所有者の環境に依存するため、このリポジトリから自動作成しない。
