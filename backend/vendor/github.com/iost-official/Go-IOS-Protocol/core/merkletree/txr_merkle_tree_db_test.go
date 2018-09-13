package merkletree

import (
	"bytes"
	"log"
	"math/rand"
	"os/exec"
	"reflect"
	"testing"
	"time"

	"github.com/iost-official/Go-IOS-Protocol/core/tx"
	"github.com/smartystreets/goconvey/convey"
)

func TestTXRMerkleTreeDB(t *testing.T) {
	convey.Convey("Test of TXRMTDB", t, func() {
		m := TXRMerkleTree{}
		txrs := []*tx.TxReceipt{
			{TxHash: []byte("node1")},
			{TxHash: []byte("node2")},
			{TxHash: []byte("node3")},
			{TxHash: []byte("node4")},
			{TxHash: []byte("node5")},
		}
		m.Build(txrs)
		Init("./")
		err := TXRMTDB.Put(&m, 32342)
		if err != nil {
			log.Panic(err)
		}
		var m_read *TXRMerkleTree
		m_read, err = TXRMTDB.Get(32342)
		if err != nil {
			log.Panic(err)
		}
		convey.So(reflect.DeepEqual(m.Tx2Txr, m_read.Tx2Txr), convey.ShouldBeTrue)
		for i := 0; i < 16; i++ {
			convey.So(bytes.Equal(m.Mt.HashList[i], m_read.Mt.HashList[i]), convey.ShouldBeTrue)
		}
		cmd := exec.Command("rm", "-r", "./TXRMerkleTreeDB")
		cmd.Run()
	})
}

func BenchmarkTXRMerkleTreeDB(b *testing.B) { //Put: 1544788ns = 1.5ms, Get: 621922ns = 0.6ms
	rand.Seed(time.Now().UnixNano())
	Init("./")
	var txrs []*tx.TxReceipt
	for i := 0; i < 3000; i++ {
		txrs = append(txrs, &tx.TxReceipt{TxHash: []byte("node1")})
	}
	m := TXRMerkleTree{}
	m.Build(txrs)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TXRMTDB.Put(&m, 33)
		TXRMTDB.Get(33)
	}
}
