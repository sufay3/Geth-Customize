/**
 * @package hook
 * @file hook.go
 * @author sufay
 *
 * abstract a hook object to interact with the blockchain
 */

package hook

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
)

var (
	HookEnabled      bool                   = false // indicate if hook is enabled. Disabled by default
	GlobalHook       *Hook                          // the global hook
	hookLogSeparator = "------------------"         // log separator
)

// EnableHook enables or disables hook according to the given bool value
func EnableHook(enabled bool) {
	HookEnabled = enabled

	if HookEnabled {
		fmt.Println(hookLogSeparator, "[Hook]: hook is enabled]", hookLogSeparator)
	}
}

// InitContext initializes the context
func InitContext(db ethdb.Database) {
	GlobalHook = New(db)
	fmt.Println(hookLogSeparator, "[Hook]: hook context is set]", hookLogSeparator)
}

// Hook for interacting with blockchain
type Hook struct {
	Db    ethdb.Database // database
	Batch ethdb.Batch    // database batch operation

	Block *types.Block       // the current block
	Tx    *types.Transaction // the current tx

	InternalTxs        []*InternalTx // the block internal txs
	CurInternalTxCount uint          // the internal tx count in the current tx
}

// New creates a hook instance
func New(db ethdb.Database) *Hook {
	return &Hook{
		Db:          db,
		Batch:       db.NewBatch(),
		InternalTxs: make([]*InternalTx, 0),
	}
}

// PrepareBlock sets the current block
func (h *Hook) PrepareBlock(block *types.Block) {
	h.Reset()
	h.Block = block
}

// PrepareTx sets the current tx
func (h *Hook) PrepareTx(tx *types.Transaction) {
	h.Tx = tx
	h.CurInternalTxCount = 0
}

// HandleTxError handles a tx error
func (h *Hook) HandleTxError(errorMessage string) error {
	txError := NewTxError(h.Tx.Hash(), errorMessage)
	return txError.Write(h.Batch)
}

// HanldeInternalTx handles an internal tx
func (h *Hook) HandleInternalTx(callType string, from common.Address, to common.Address, gas uint64, value *big.Int) {
	internalTx := NewInternalTx(h.Block.NumberU64(), h.Block.Hash(), h.Tx.Hash(), callType, from, to, gas, value)
	h.InternalTxs = append(h.InternalTxs, internalTx)
	h.CurInternalTxCount++
}

// DiscardCurInternalTxs discards the internal txs of the current tx
func (h *Hook) DiscardCurInternalTxs() {
	h.InternalTxs = h.InternalTxs[:len(h.InternalTxs)-int(h.CurInternalTxCount)]
}

// FinalizeBlock commits the batch operation of the current block
func (h *Hook) FinalizeBlock() error {
	err := h.Batch.Write()
	if err != nil {
		return err
	}

	blockInternalTxs := NewBlockInternalTxs(h.Block.Hash(), h.InternalTxs)
	return blockInternalTxs.Write(h.Batch)
}

// Reset discards the current block
func (h *Hook) Reset() {
	h.Batch.Reset()
	h.InternalTxs = make([]*InternalTx, 0)
}

// GetTxError gets the tx error by the specified tx hash
func (h *Hook) GetTxError(txHash common.Hash) (*TxError, error) {
	bytes, err := h.Db.Get(MakeTxErrorKey(txHash))
	if err != nil {
		return nil, err
	}

	txError := new(TxError)

	if err = RlpDecode(bytes, txError); err != nil {
		return nil, err
	}

	return txError, nil
}

// GetBlockInternalTxCount gets the block internal tx count by the specified block hash
func (h *Hook) GetBlockInternalTxCount(blockHash common.Hash) (uint, error) {
	bytes, err := h.Db.Get(MakeBlockInternalTxCountKey(blockHash))
	if err != nil {
		return 0, err
	}

	txCount := new(uint)

	if err = RlpDecode(bytes, txCount); err != nil {
		return 0, err
	}

	return *txCount, nil
}

// GetBlockInternalTxs gets the block internal txs by the specified block hash
func (h *Hook) GetBlockInternalTxs(blockHash common.Hash) ([]*InternalTx, error) {
	bytes, err := h.Db.Get(MakeBlockInternalTxsKey(blockHash))
	if err != nil {
		return nil, err
	}

	internalTxs := make([]*InternalTx, 0)

	if err = RlpDecode(bytes, &internalTxs); err != nil {
		return nil, err
	}

	return internalTxs, nil
}
