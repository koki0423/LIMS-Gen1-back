# 備品管理システムについて
## API
### 認証不要エンドポイント
- APIのベースグループ /api/v1
  - 備品登録エンドポイント
    -   GET /api/v1/ping
        -   応答性確認
    -   GET /api/v1/assets/:id
        -   備品ID(assetID)を基に詳細取得
    -   GET /api/v1/assetsAll
        -   全備品リスト取得
    -   POST /api/v1/assets
        -   備品登録
        -   ```json:リクエストボディ
              //個別管理
              {
                  "asset_master": {
                      "name": "Test Equipment",
                      "management_category_id": 1,
                      "genre_id": 1,
                      "manufacturer": "hogehoge.inc",
                      "model_number": "FUGA-1234 v10"
                  },
                  "asset": {
                      "item_master_id": 0,
                      "quantity": 1,
                      "serial_number": "SN-1234567890",
                      "status_id": 1,
                      "purchase_date": "2023-10-01",
                      "owner": "Noa Ushio",
                      "location": "seminar room 1",
                      "last_check_date": "2023-10-15",
                      "last_checker": "Yuka Hayase",
                      "notes": null
                  }
              }
              //全体管理
              {
                  "asset_master": {
                      "name": "Test Equipment",
                      "management_category_id": 2,
                      "genre_id": 5,
                      "manufacturer": "秋月電子通商",
                      "model_number": "白色LED 5mm　150W"
                  },
                  "asset": {
                      "item_master_id": 0,
                      "quantity": 6,
                      "serial_number": null,
                      "status_id": 1,
                      "purchase_date": "2023-10-01",
                      "owner": "Noa Ushio",
                      "location": "seminar room 1",
                      "last_check_date": "2023-10-15",
                      "last_checker": "Yuka Hayase",
                      "notes": null
                  }
              }
  - 備品貸出エンドポイント
    -  POST /api/v1/borrow/:id
       -  id はassetsテーブルのID（主キー）を指定
       -  ```json:リクエストボディ
            {
                "borrower": "Nozomi",
                "quantity": 12,
                "lend_date": "2025-01-20",
                "expected_return_date": null,
                "actual_return_date": null,
                "Notes": null
            }
    -  POST /api/v1/return/:id
       -  id はasset_lendsテーブルのID（主キー）を指定
    -  GET /api/v1/borrows
       -  貸出中一覧を取得するエンドポイント



### 認証必須エンドポイント


    