package main

import (
	"database/sql"

	// "equipmentManager/assets"
	"equipmentManager/internal/database/db"
	// "equipmentManager/lends"
	// "equipmentManager/disposal"
	// "equipmentManager/developer"
	"equipmentManager/internal/database/model/tables"
	"equipmentManager/internal/database/utils"
	_ "github.com/go-sql-driver/mysql"
	// "log"
)

var assetsIndivisual model.AllDataOfAsset
var assetsAllManagement model.AllDataOfAsset

var disposalAssets model.AssetsDisposal
var disposalAssetsAllManagement model.AssetsDisposal

var lendAssetsAllManagement model.AssetsLend

func main() {
	db, _ := db.ConnectDB()
	defer db.Close()

	initAssets()
	initDisposalAssets()
	initLendAssets()

	AssetTest(db)     // 資産登録のテスト
	DisporsalTest(db) // 資産廃棄登録のテスト
	LendTest(db)      // 資産貸出登録のテスト
	DeveloperTest(db) // 開発者登録と認証のテスト
}

func initAssets() {
	assetsIndivisual = model.AllDataOfAsset{
		Master: model.AssetsMaster{
			Name:                 "Test Equipment",
			ManagementCategoryID: 1, // 個別管理の管理区分ID
			GenreID:              sql.NullInt64{Int64: 1, Valid: true},
			Manufacturer:         sql.NullString{String: "hogehoge.inc", Valid: true},
			ModelNumber:          sql.NullString{String: "FUGA-1234 v10", Valid: true},
		},
		Asset: model.Asset{
			ItemMasterID:  0, // 初期値は0、登録後に更新される
			Quantity:      1,
			SerialNumber:  sql.NullString{String: "SN-1234567890", Valid: true},
			StatusID:      1,
			PurchaseDate:  sql.NullTime{Time: utils.MustParseDate("2023-10-01"), Valid: true},
			Owner:         sql.NullString{String: "Noa Ushio", Valid: true},
			Location:      sql.NullString{String: "seminar room 1", Valid: true},
			LastCheckDate: sql.NullTime{Time: utils.MustParseDate("2023-10-15"), Valid: true},
			LastChecker:   sql.NullString{String: "Yuka Hayase", Valid: true},
			Notes:         sql.NullString{Valid: false},
		},
	}
	assetsAllManagement = model.AllDataOfAsset{
		Master: model.AssetsMaster{
			Name:                 "Test Equipment",
			ManagementCategoryID: 2,                                    // 全体管理の管理区分ID
			GenreID:              sql.NullInt64{Int64: 4, Valid: true}, //1: 個人 2: 事務 3: ファシリティ 4:組込みシステム 5: 高度情報演習
			Manufacturer:         sql.NullString{String: "秋月電子通商", Valid: true},
			ModelNumber:          sql.NullString{String: "赤色LED 5mm", Valid: true},
		},
		Asset: model.Asset{
			ItemMasterID:  0, // 初期値は0、登録後に更新される
			Quantity:      2173,
			SerialNumber:  sql.NullString{Valid: false}, // 全体管理では個別のシリアル番号は不要
			StatusID:      1,
			PurchaseDate:  sql.NullTime{Time: utils.MustParseDate("2023-10-01"), Valid: true},
			Owner:         sql.NullString{String: "Noa Ushio", Valid: true},
			Location:      sql.NullString{String: "seminar room 1", Valid: true},
			LastCheckDate: sql.NullTime{Time: utils.MustParseDate("2023-10-15"), Valid: true},
			LastChecker:   sql.NullString{String: "Yuka Hayase", Valid: true},
			Notes:         sql.NullString{Valid: false},
		},
	}
}

func initDisposalAssets() {
	disposalAssets = model.AssetsDisposal{
		AssetID:      6,                                                // 資産ID
		Quantity:     1,                                                // 廃棄数量
		DisposalDate: utils.MustParseDate("2024-08-25"),                // 廃棄日
		Reason:       sql.NullString{String: "故障のため", Valid: true},     // 廃棄理由
		ProcessedBy:  sql.NullString{String: "Noa Ushio", Valid: true}, // 処理担当者
		IsIndividual: true,                                             // 個別管理かどうか
	}

	disposalAssetsAllManagement = model.AssetsDisposal{
		AssetID:      7,                                                // 資産ID
		Quantity:     100,                                              // 廃棄数量
		DisposalDate: utils.MustParseDate("2024-08-25"),                // 廃棄日
		Reason:       sql.NullString{String: "不要になったため", Valid: true},  // 廃棄理由
		ProcessedBy:  sql.NullString{String: "Noa Ushio", Valid: true}, // 処理担当者
		IsIndividual: false,                                            // 全体管理
	}
}

func initLendAssets() {
	lendAssetsAllManagement = model.AssetsLend{
		AssetID:            4,                                                                  // 全体管理のassets.idを想定
		Borrower:           "Nonomi Izayoi",                                                    // 借用者名
		Quantity:           2,                                                                  // 貸出数
		LendDate:           utils.MustParseDate("2024-12-31"),                                  // 貸出日
		ExpectedReturnDate: sql.NullTime{Time: utils.MustParseDate("2025-03-31"), Valid: true}, // 返却予定日
		ActualReturnDate:   sql.NullTime{Valid: false},                                         // まだ返却されていない
		Notes:              sql.NullString{String: "LED実験用に一時貸出", Valid: true},                 // 備考あり
	}
}

func AssetTest(db *sql.DB) {
	// assets.RegistrationEquipmentForIndivisual(db, assetsIndivisual) //OK
	// assets.RegistrationEquipmentForGeneral(db,assetsAllManagement) //OK

	//資産取得テスト
	// list, err := assets.FetchAssetsAll(db)
	// log.Printf("全資産情報一覧: %v,  エラー: %v", list, err)
}

func DisporsalTest(db *sql.DB) {
	// 個別管理の資産廃棄登録テスト //OK
	// disposal.RegisterDisposal(db,disposalAssets)

	// 全体管理の資産廃棄登録テスト //OK
	// disposal.RegisterDisposal(db,disposalAssetsAllManagement)

	// 資産廃棄登録の取り消しテスト // OK
	// disposal.UndoRegisterDisposal(db, 7)
}

func LendTest(db *sql.DB) {

	//貸出登録テスト //OK
	// lends.RegisterLend(db, lendAssetsAllManagement)

	// 貸出情報の取得テスト // OK
	// list, err := lends.FetchLendsAll(db)
	// log.Printf("全貸出情報一覧: %v,  エラー: %v", list, err)

	// 資産IDによる貸出情報の取得テスト // OK
	// list, err := lends.FetchLendsByAssetID(db, 4)
	// log.Printf("資産IDによる貸出情報一覧: %v,  エラー: %v", list, err)

	// lendAssetsAllManagement=list[0] // 取得した貸出情報の最初の要素を更新テスト用に使用

	// 貸出情報の更新テスト // OK
	// lendAssetsAllManagement.Notes=sql.NullString{String:"9月末までの貸出",Valid: true}
	// lends.UpdateLend(db, lendAssetsAllManagement)
}

func DeveloperTest(db *sql.DB) {
	// var user model.Developer
	// user.StudentNumber = "12345678" // テスト用学籍番号
	//user.PasswordHash, _ = utils.HashPassword("12345") // テスト用パスワードハッシュ
	// inputPassword := "12345" //ユーザーからの入力値を想定（平文）

	// 開発者登録のテスト //OK
	// developer.RegisterDeveloper(db, user)

	// 開発者認証のテスト // OK
	// developer.CheckDeveloper(db, user.StudentNumber, inputPassword)
}
