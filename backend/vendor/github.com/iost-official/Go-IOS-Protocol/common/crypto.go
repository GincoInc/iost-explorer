package common

import (
	"crypto/sha256"
	"encoding/hex"

	"encoding/binary"
	"hash/crc32"

	"github.com/btcsuite/btcutil/base58"
	"github.com/iost-official/Go-IOS-Protocol/common/hash"
	"golang.org/x/crypto/ripemd160"
)

func Sha256(raw []byte) []byte {
	var data = sha256.Sum256(raw)
	return data[:]
}

func Sha3(raw []byte) []byte {
	var data = hash.Sha3_256(raw)
	return data[:]
}

func Hash160(raw []byte) []byte {
	var data = sha256.Sum256(raw)
	return ripemd160.New().Sum(data[len(data):])
}

func Base58Encode(raw []byte) string {
	return base58.Encode(raw)
}

func Base58Decode(s string) []byte {
	return base58.Decode(s)
}

func Parity(bit []byte) []byte {
	crc32q := crc32.MakeTable(crc32.Koopman)
	crc := crc32.Checksum(bit, crc32q)
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, crc)
	return bs
}

func ToHex(data []byte) string {
	return hex.EncodeToString(data)
}

func ParseHex(s string) []byte {
	d, err := hex.DecodeString(s)
	if err != nil {
		println(err)
		return nil
	}
	return d
}
