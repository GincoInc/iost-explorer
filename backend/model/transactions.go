package model

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/gogo/protobuf/proto"
	"github.com/GincoInc/iost-explorer/backend/model/blockchain/rpcpb"
	"github.com/GincoInc/iost-explorer/backend/model/db"
	contract "github.com/GincoInc/iost-explorer/backend/model/pb"
	"github.com/GincoInc/iost-explorer/backend/util"
)

/// this struct is used as json to return
type TxnDetail struct {
	Hash          string  `json:"txHash"`
	BlockNumber   int64   `json:"blockHeight"`
	From          string  `json:"from"`
	To            string  `json:"to"`
	Amount        float64 `json:"amount"`
	GasLimit      float64 `json:"gasLimit"`
	GasPrice      float64 `json:"price"`
	Age           string  `json:"age"`
	UTCTime       string  `json:"utcTime"`
	Code          string  `json:"code"`
	StatusCode    int32   `json:"statusCode"`
	StatusMessage string  `json:"statusMessage"`
	Contract      string  `json:"contract"`
	ActionName    string  `json:"actionName"`
	Data          string  `json:"data"`
	Memo          string  `json:"memo"`
}

type TxJson struct {
	Hash        string  `json:"hash"`
	BlockNumber int64   `json:"blockNumber"`
	From        string  `json:"from"`
	To          string  `json:"to"`
	Amount      float64 `json:"amount"`
	GasLimit    float64 `json:"gasLimit"`
	GasPrice    float64 `json:"gasPrice"`
	Age         string  `json:"age"`
	UTCTime     string  `json:"utcTime"`
}

func ConvertTxJsons(txs []*db.TxStore) []*TxJson {
	var ret []*TxJson
	for _, tx := range txs {
		txnOut := &TxJson{
			Hash:        tx.Tx.Hash,
			BlockNumber: tx.BlockNumber,
			From:        tx.Tx.Publisher,
			To:          tx.Tx.Actions[0].Contract,
			GasLimit:    tx.Tx.GasLimit,
			GasPrice:    tx.Tx.GasRatio,
			Age:         util.ModifyIntToTimeStr(tx.Tx.Time / 1e9),
			UTCTime:     util.FormatUTCTime(tx.Tx.Time),
		}

		if tx.Tx.Actions[0].Contract == "token.iost" && tx.Tx.Actions[0].ActionName == "transfer" &&
			tx.Tx.TxReceipt.StatusCode == rpcpb.TxReceipt_SUCCESS {
			var params []string
			err := json.Unmarshal([]byte(tx.Tx.Actions[0].Data), &params)
			if err == nil && len(params) == 5 && params[0] == "iost" {
				txnOut.From = params[1]
				txnOut.To = params[2]
				f := getIOSTAmount(params[3])
				txnOut.Amount = f
			}
		} else if tx.Tx.Actions[0].Contract == "exchange.iost" && tx.Tx.Actions[0].ActionName == "transfer" &&
			tx.Tx.TxReceipt.StatusCode == rpcpb.TxReceipt_SUCCESS {
			var params []string
			err := json.Unmarshal([]byte(tx.Tx.Actions[0].Data), &params)
			if err == nil && len(params) == 4 && params[0] == "iost" {
				if params[1] != "" {
					txnOut.From = tx.Tx.Publisher
					txnOut.To = params[1]
					f := getIOSTAmount(params[2])
					txnOut.Amount = f
				}
			}
		}
		ret = append(ret, txnOut)
	}
	return ret
}

func ConvertTxsOutputs(tx []*db.TxStore) []*TxnDetail {
	var ret []*TxnDetail
	for _, t := range tx {
		ret = append(ret, ConvertTxOutput(t))
	}
	return ret
}

func getIOSTAmount(s string) float64 {
	point := strings.Index(s, ".")
	if point > 0 && point+9 <= len(s) {
		s = s[:point+9]
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

/// convert FlatTx to TxnDetail
func ConvertTxOutput(tx *db.TxStore) *TxnDetail {
	txnOut := &TxnDetail{
		Hash:          tx.Tx.Hash,
		BlockNumber:   tx.BlockNumber,
		From:          tx.Tx.Publisher,
		To:            tx.Tx.Actions[0].Contract,
		GasLimit:      tx.Tx.GasLimit,
		GasPrice:      tx.Tx.GasRatio,
		Age:           util.ModifyIntToTimeStr(tx.Tx.Time / 1e9),
		UTCTime:       util.FormatUTCTime(tx.Tx.Time),
		Code:          "",
		StatusCode:    int32(tx.Tx.TxReceipt.StatusCode),
		StatusMessage: tx.Tx.TxReceipt.Message,
		Contract:      tx.Tx.Actions[0].Contract,
		ActionName:    tx.Tx.Actions[0].ActionName,
		Data:          tx.Tx.Actions[0].Data,
	}

	if tx.Tx.Actions[0].Contract == "token.iost" && tx.Tx.Actions[0].ActionName == "transfer" &&
		tx.Tx.TxReceipt.StatusCode == rpcpb.TxReceipt_SUCCESS {
		var params []string
		err := json.Unmarshal([]byte(tx.Tx.Actions[0].Data), &params)
		if err == nil && len(params) == 5 && params[0] == "iost" {
			txnOut.From = params[1]
			txnOut.To = params[2]
			f := getIOSTAmount(params[3])
			txnOut.Amount = f * 1e8
			txnOut.Memo = params[4]
		}
	} else if tx.Tx.Actions[0].Contract == "exchange.iost" && tx.Tx.Actions[0].ActionName == "transfer" &&
		tx.Tx.TxReceipt.StatusCode == rpcpb.TxReceipt_SUCCESS {
		var params []string
		err := json.Unmarshal([]byte(tx.Tx.Actions[0].Data), &params)
		if err == nil && len(params) == 4 && params[0] == "iost" {
			if params[1] != "" {
				txnOut.From = tx.Tx.Publisher
				txnOut.To = params[1]
				f := getIOSTAmount(params[2])
				txnOut.Amount = f * 1e8
			}
		}
	}

	if tx.Tx.Actions[0].Contract == "system.iost" && tx.Tx.Actions[0].ActionName == "setCode" {
		var params []string
		err := json.Unmarshal([]byte(tx.Tx.Actions[0].Data), &params)
		if err == nil && len(params) > 0 {
			j, e := simplejson.NewJson([]byte(params[0]))
			if e != nil {
				log.Printf("json decode setCode param failed. err=%v", err)
			} else {
				txnOut.Code, _ = j.Get("code").String()
			}

			if txnOut.Code == "" {
				bytes, e := base64.StdEncoding.DecodeString(params[0])
				if e != nil {
					log.Printf("base64 decode setCode param failed. err=%v", err)
				} else {
					var con contract.Contract
					proto.Unmarshal(bytes, &con)
					txnOut.Code = con.Code
				}
			}
		}

		// c.B64Decode(code[0])
		/* txnOut.Code = c.Code */
	}

	return txnOut
}

func GetDetailTxn(txHash string) (*TxnDetail, error) {
	tx, err := db.GetTxByHash(txHash)

	if err != nil {
		log.Printf("transaction not found. txHash:%v, err=%v", txHash, err)
		return nil, err
	}

	txnOut := ConvertTxOutput(tx)
	txnOut.Amount /= 1e8

	return txnOut, nil
}

/// get a list of transactions for a specific page using account and block
func GetFlatTxnSlicePage(page, eachPageNum, block int64) ([]*TxnDetail, error) {
	lastPageNum, err := db.GetTxTotalPageCnt(eachPageNum, block)
	if err != nil {
		return nil, err
	}

	if lastPageNum == 0 {
		return []*TxnDetail{}, nil
	}

	if page > lastPageNum {
		page = lastPageNum
	}

	start := int((page - 1) * eachPageNum)
	txnsFlat, err := db.GetTxs(block, start, int(eachPageNum))

	if err != nil {
		return nil, err
	}

	return ConvertTxsOutputs(txnsFlat), nil
}
