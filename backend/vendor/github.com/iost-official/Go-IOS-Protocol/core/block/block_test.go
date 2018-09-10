package block

import (
	"bytes"
	"testing"

	"github.com/iost-official/Go-IOS-Protocol/account"
	"github.com/iost-official/Go-IOS-Protocol/core/tx"
	"github.com/iost-official/Go-IOS-Protocol/crypto"
	"github.com/smartystreets/goconvey/convey"
)

func TestBlockHeadSerialize(t *testing.T) {
	convey.Convey("Test of block head encode and decode", t, func() {
		head := BlockHead{
			Number:     1,
			ParentHash: []byte("parent"),
		}
		hash, err := head.Encode()
		convey.So(err, convey.ShouldBeNil)
		var headRead BlockHead
		err = headRead.Decode(hash)
		convey.So(err, convey.ShouldBeNil)
		convey.So(bytes.Equal(head.ParentHash, headRead.ParentHash), convey.ShouldBeTrue)
		convey.So(headRead.Number == head.Number, convey.ShouldBeTrue)
	})
}

func TestBlockSerialize(t *testing.T) {
	convey.Convey("test Push", t, func() {
		blk := Block{
			Head: &BlockHead{
				Number:     1,
				ParentHash: []byte("parent"),
			},
			Sign: &crypto.Signature{},
		}
		a1, _ := account.NewAccount(nil)
		tx0 := tx.Tx{
			Time: 1,
			Actions: []*tx.Action{{
				Contract:   "contract1",
				ActionName: "actionname1",
				Data:       "{\"num\": 1, \"message\": \"contract1\"}",
			}},
			Signers: [][]byte{a1.Pubkey},
		}
		blk.Txs = append(blk.Txs, &tx0)
		receipt := tx.TxReceipt{
			TxHash:   tx0.Hash(),
			GasUsage: 10,
			Status: tx.Status{
				Code:    tx.Success,
				Message: "run success",
			},
		}
		blk.Receipts = append(blk.Receipts, &receipt)
		convey.Convey("Test of block encode and decode", func() {
			blkByte, err := blk.Encode()
			blkRead := Block{}
			err = blkRead.Decode(blkByte)
			convey.So(err, convey.ShouldBeNil)
			convey.So(bytes.Equal(blkRead.Head.ParentHash, blk.Head.ParentHash), convey.ShouldBeTrue)
			convey.So(len(blkRead.Txs) == len(blk.Txs), convey.ShouldBeTrue)
			convey.So(len(blkRead.Receipts) == len(blk.Receipts), convey.ShouldBeTrue)
			convey.So(bytes.Equal(blkRead.Receipts[0].TxHash, tx0.Hash()), convey.ShouldBeTrue)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}
