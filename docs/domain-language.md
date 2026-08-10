# ドメイン用語

このドキュメントは、賃貸物件を扱うRoom Finderで使用する英語のドメイン用語を定義する。
実装・API・DB・Agentのプロンプトでは、ここで定義した正規語を使用する。

## 正規語

### Listing

賃貸サイトに掲載されている、1件の募集情報を指す。

Room Finderが現在保存する対象は建物そのものではなく、LIFULL HOME'Sなどのサイトに掲載された募集情報であるため、正規語は`Listing`とする。

使用例:

- `Listing`
- `DBTableListing`
- `listings`
- `source_listing_id`
- `SaveListingRequest`
- `/rental-listings`

同じサイト上の同じ掲載IDを再送した場合は、同一Listingの現在情報を更新する。

### Rental Listing

文脈だけでは賃貸掲載情報であることが分かりにくい場合に使用する正式名称。
特にドキュメント、API仕様、Agentプロンプトでは`Rental Listing`を使用してよい。

コード上の型名・テーブル名では、Room Finder内の文脈が明確なため`Listing`を優先する。

## 使い分ける用語

### Unit

建物内の実際の部屋・住戸を指す。
複数の掲載情報が同じ住戸を指す場合の実体モデルが必要になったときに使用する。
現MVPでは、住戸の同一性管理は行わない。

### Building

マンション・アパートなどの建物を指す。
現MVPでは、建物情報と募集情報を分離しない。

### Search Profile

初回プロンプトと、そのプロンプトから解釈された絶対条件をまとめた検索設定を指す。
作成後は変更せず、条件を変える場合は新しいSearch Profileを作成する。

## 使用しない語

### Property

一般的な英語としては不動産を意味できるが、プログラミングではオブジェクトの属性も意味するため、Room Finderのドメイン用語としては使用しない。

既存のIssue #2実装には移行前の名前として、`properties`テーブルと`DBTableProperty`が残っている。今後のAPI・Agent・新規コードでは`Listing`を使用し、既存DB名を変更する場合は専用マイグレーションで段階的に移行する。

### Room

「部屋」の意味が広く、住戸・掲載情報・部屋タイプのどれを指すか曖昧になるため、単独のドメイン名には使用しない。

## 命名ルール

| 概念 | 正規語 | 例 |
| --- | --- | --- |
| 掲載情報 | Listing | `Listing`, `listings` |
| 取得元の掲載ID | Source Listing ID | `source_listing_id` |
| 実際の住戸 | Unit | `unit_id`（将来必要になった場合のみ） |
| 建物 | Building | `building_id`（将来必要になった場合のみ） |
| 検索設定 | Search Profile | `search_profile_id` |

新しい用語が必要になった場合は、既存語の流用で済ませず、このドキュメントへ定義を追加してから実装する。
