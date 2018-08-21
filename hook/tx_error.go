/**
 * @package hook
 * @file tx_error.go
 * @author sufay
 *
 * used to handle the tx error which is produced by evm
 * during a transaction executing.
 */

package hook

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

// TxError for the error abstraction of a tx
type TxError struct {
	TxHash  common.Hash // the hash of the tx where the error happened
	Message string      // the error message
}

// NewTxError creates a new TxError instance
func NewTxError(txHash common.Hash, errorMessage string) *TxError {
	return &TxError{
		TxHash:  txHash,
		Message: errorMessage,
	}
}

// SetBlock sets the block
func (te *TxError) SetTxHash(txHash common.Hash) {
	te.TxHash = txHash
}

// SetMessage sets the error message
func (te *TxError) SetMessage(message string) {
	te.Message = message
}

// Serialize serializes the TxError instance using rlp encoding.
func (te *TxError) Serialize() []byte {
	return MustRlpEncode(te)
}

// Write writes the TxError object encoded to database
func (te *TxError) Write(batch ethdb.Batch) error {
	return batch.Put(MakeTxErrorKey(te.TxHash), te.Serialize())
}
