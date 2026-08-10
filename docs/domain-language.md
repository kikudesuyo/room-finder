# ドメイン用語

このドキュメントは、賃貸物件を扱うRoom Finderで使用する英語のドメイン用語を定義する。
実装・API・DB・Agentのプロンプトでは、ここで定義した正規語を使用する。

## 正規語

### RentalOffer

賃貸サイトで募集されている、1件の賃貸情報を指す。

Room Finderが現在保存する対象は建物そのものではなく、LIFULL HOME'Sなどで提供されている賃貸募集情報であるため、正規語は`RentalOffer`とする。

使用例:

- `RentalOffer`
- `DBTableRentalOffer`
- `rental_offers`
- `source_offer_id`
- `SaveRentalOfferRequest`
- `/rental-offers`

同じサイト上の同じ募集IDを再送した場合は、同一RentalOfferの現在情報を更新する。

## 使い分ける用語

### Unit

建物内の実際の部屋・住戸を指す。
複数の掲載情報が同じ住戸を指す場合の実体モデルが必要になったときに使用する。
現MVPでは、住戸の同一性管理は行わない。

### Building

マンション・アパートなどの建物を指す。
現MVPでは、建物情報と募集情報を分離しない。

### SearchCondition

ユーザーが指定する1件分の検索条件を指す。作成後は変更せず、条件を変える場合は新しいSearchConditionを作成する。

## 使用しない語

### Property

一般的な英語としては不動産を意味できるが、プログラミングではオブジェクトの属性も意味するため、Room Finderのドメイン用語としては使用しない。

既存のIssue #2実装には移行前の名前として`properties`テーブルが存在する。専用マイグレーションで`rental_offers`へ移行する。

### Room

「部屋」の意味が広く、住戸・掲載情報・部屋タイプのどれを指すか曖昧になるため、単独のドメイン名には使用しない。

## 命名ルール

| 概念 | 正規語 | 例 |
| --- | --- | --- |
| 賃貸募集情報 | RentalOffer | `RentalOffer`, `rental_offers` |
| 取得元の募集ID | Source Offer ID | `source_offer_id` |
| 実際の住戸 | Unit | `unit_id`（将来必要になった場合のみ） |
| 建物 | Building | `building_id`（将来必要になった場合のみ） |
| 検索条件 | SearchCondition | `search_condition_id` |

新しい用語が必要になった場合は、既存語の流用で済ませず、このドキュメントへ定義を追加してから実装する。
