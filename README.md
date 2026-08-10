# Room Finder

賃貸物件を取得し、自然言語で指定した必須条件に一致する物件を保存・表示するサービスです。

## アーキテクチャ

buildlogと同じ構成を採用します。

- Web: SvelteKit / TypeScript / Tailwind CSS
- API: Go / chi / GORM
- Local DB: Docker Compose / PostgreSQL
- Migration: Atlas
- Production: Cloudflare Pages / GCP Cloud Run / Neon PostgreSQL

詳細は [SPEC.md](./SPEC.md) と [docs/roadmap.md](./docs/roadmap.md) を参照してください。

## ローカル開発

### 1. 環境変数

```bash
cp .env.example .env
```

### 2. PostgreSQL起動

```bash
docker compose up -d
```

### 3. Go API起動

```bash
cd api
make dev
```

APIは http://localhost:8081 で起動します。

疎通確認:

```bash
curl http://localhost:8081/api/v1/health
```

### 4. Web起動

```bash
cd web
pnpm install
pnpm run dev
```

Webは http://localhost:5173 で起動します。

### 5. Atlas

DBスキーマを追加した後は、`sql/` 配下でAtlasを実行します。

```bash
cd sql
atlas migrate apply --env local
```

## 開発ルール

Issue、ブランチ、Pull Request、検証のルールは [AGENTS.md](./AGENTS.md) に従ってください。
