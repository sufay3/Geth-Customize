/**
 * @package hook
 * @file utils.go
 * @author sufay
 *
 * used to provide some utility functions
 */

package hook

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	keyPrefixTxError              = []byte("e")  // the key prefix of the tx error
	keyPrefixBlockInternalTxCount = []byte("ic") // the key prefix of the block internal tx count
	keyPrefixBlockInternalTxs     = []byte("i")  // the key prefix of the block internal txs
)

// MakeTxErrorKey makes a tx error key from the specified tx hash
func MakeTxErrorKey(txHash common.Hash) []byte {
	return append(keyPrefixTxError, txHash.Bytes()...)
}

// MakeBlockInternalTxCountKey makes a key of the block internal tx count from the specified block hash
func MakeBlockInternalTxCountKey(blockHash common.Hash) []byte {
	return append(keyPrefixBlockInternalTxCount, blockHash.Bytes()...)
}

// MakeBlockInternalTxsKey makes a key of the block internal txs from the specified block hash
func MakeBlockInternalTxsKey(blockHash common.Hash) []byte {
	return append(keyPrefixBlockInternalTxs, blockHash.Bytes()...)
}

// MustRlpEncode encodes the specified input using rlp
// panic when an error occurs
func MustRlpEncode(in interface{}) []byte {
	bytes, err := rlp.EncodeToBytes(in)

	if err != nil {
		panic(err)
	}

	return bytes
}

// RlpDecode decodes the specified data to value
func RlpDecode(data []byte, value interface{}) error {
	return rlp.DecodeBytes(data, value)
}
