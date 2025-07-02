package service

import (
	"bytes"
	// "encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync/atomic"
	"time"

	"equipmentManager/domain"

	"github.com/ebfe/scard"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type NfcService struct {
	Dummy           bool         // ダミー用フラグ
	latestStudentID atomic.Value // スレッドセーフ
}

func NewNfcService() *NfcService {
	return &NfcService{
		Dummy: true,
	}
}

func (n *NfcService) RunNFCReader() {
	if n.Dummy {
		// ダミーデータ
		n.latestStudentID.Store("MA12345")
		log.Println("NFC Service running in dummy mode with student ID: MA12345")
		return
	}
	for {
		id, err := n.GetNFCData()
		if err == nil && id != "" {
			n.latestStudentID.Store(id)
		}
	}
}

func (n *NfcService) GetLatestStudentID() string {
	val := n.latestStudentID.Load()
	if val == nil {
		return ""
	}
	return val.(string)
}

func (n *NfcService) Close() {
	// ここでは特にリソース解放は不要だが、必要なら実装
	log.Println("NFC Service closed")
}

func (n *NfcService) GetNFCData() (string, error) {
	// --- PC/SC初期化 ---
	ctx, err := scard.EstablishContext()
	if err != nil {
		log.Fatalf("PC/SC context error: %v", err)
	}
	defer ctx.Release()

	readers, err := ctx.ListReaders()
	if err != nil || len(readers) == 0 {
		log.Fatalf("No NFC reader found")
	}

	// --- カード挿入待ち ---
	var card *scard.Card
	for {
		card, err = ctx.Connect(readers[0], scard.ShareShared, scard.ProtocolAny)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	defer card.Disconnect(scard.LeaveCard)

	// --- Polling (IDm取得) ---
	_, err = pollingFeliCa(card, []byte{0x82, 0x77}) // system_code: 0x8277
	if err != nil {
		log.Fatalf("Polling failed: %v", err)
	}
	// fmt.Printf("IDm: %s\n", hex.EncodeToString(idm))

	// --- 学生番号ブロック読出し ---
	blockData, err := readBlock(card, 0x010b, 0)
	if err != nil {
		log.Fatalf("Read block failed: %v", err)
	}

	// --- shift_jis→utf8変換 ---
	utf8data, err := shiftJISToUTF8(blockData)
	if err != nil {
		return "", err
	}
	// 学籍番号が格納されている部分だけ抜き出し
	studentNum := ""
	if len(utf8data) >= 10 {
		studentNum = string(utf8data[3:10])
	} else {
		studentNum = string(utf8data)
	}

	// JSONエンコードしてstringで返す
	jsonData, err := json.Marshal(&domain.StudentCardInfo{
		StudentNumber: studentNum,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func shiftJISToUTF8(data []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(data), japanese.ShiftJIS.NewDecoder())
	return ioutil.ReadAll(reader)
}

func toHexString(b []byte) string {
	result := ""
	for i, v := range b {
		result += fmt.Sprintf("%02X", v)
		if i != len(b)-1 {
			result += " "
		}
	}
	return result
}

func transmitAndPrint(card *scard.Card, apdu []byte, label string) {
	fmt.Printf("> %s: %s\n", label, toHexString(apdu))
	response, err := card.Transmit(apdu)
	if err != nil {
		log.Fatalf("Transmit error: %v", err)
	}
	// 最後の2バイトはSW1, SW2
	if len(response) < 2 {
		fmt.Printf("< (invalid response)\n")
		return
	}
	data := response[:len(response)-2]
	sw1 := response[len(response)-2]
	sw2 := response[len(response)-1]
	fmt.Printf("< %s %02X %02X\n", toHexString(data), sw1, sw2)
}

// Pollingコマンド
func pollingFeliCa(card *scard.Card, systemCode []byte) ([]byte, error) {
	// FeliCa Polling: 0x00 FF FF system_code_1 system_code_2 01 00
	command := []byte{0xFF, 0xCA, 0x00, 0x00, 0x00} // NOTE: 多くのPC/SCリーダではFeliCaのPolling直接APDUは不可。SDK仕様次第
	// 多くはIDm取得用のAPDU (Get Data)で代用できる場合がある
	// 本格的にPollingしたい場合はPaSoRiなどリーダ独自APIが必要なことも
	// ここでは一般的なIDm取得APDUを使う
	resp, err := card.Transmit(command)
	if err != nil {
		return nil, err
	}
	// 最後2バイトはSW1, SW2
	if len(resp) < 2 {
		return nil, fmt.Errorf("too short")
	}
	idm := resp[:len(resp)-2]
	return idm, nil
}

// FeliCa read without encryption
func readBlock(card *scard.Card, serviceCode uint16, blockNum byte) ([]byte, error) {
	// ここではパナソニックRC-S380等で一般的なFeliCa READコマンドをAPDUでラップして送る形式
	// 一般的なICカードリーダー用APDUの例
	// [FF, B0, 00, blockNum, blockLen]
	apdu := []byte{0xFF, 0xB0, 0x00, blockNum, 0x10} // 16バイト読み出し
	resp, err := card.Transmit(apdu)
	if err != nil {
		return nil, err
	}
	if len(resp) < 2 {
		return nil, fmt.Errorf("too short")
	}
	return resp[:len(resp)-2], nil
}
