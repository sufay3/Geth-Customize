/**
 * @package hook
 * @file internal_tx.go
 * @author sufay
 *
 * used to handle the internal tx which is produced by evm call
 * during a contract transaction processing.
 */

package hook

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

// InternalTx for the internal tx of a contract transaction
type InternalTx struct {
	BlockNumber  uint64         // the number of the block where the internal tx lies
	BlockHash    common.Hash    // the hash of the block where the internal tx lies
	ParentTxHash common.Hash    // the paranet hash of the internal tx
	CallType     string         // the call type of the internal tx
	From         common.Address // the from address of the internal tx
	To           common.Address // the to address of the internal tx
	GasLimit     uint64         // the gas limit of the internal tx
	Value        *big.Int       // the value of the internal tx
}

// NewInternalTx creates an internal tx
func NewInternalTx(blockNumber uint64, blockHash common.Hash, parent common.Hash, callType string, from common.Address, to common.Address, gas uint64, value *big.Int) *InternalTx {
	return &InternalTx{
		BlockNumber:  blockNumber,
		BlockHash:    blockHash,
		ParentTxHash: parent,
		CallType:     callType,
		From:         from,
		To:           to,
		GasLimit:     gas,
		Value:        new(big.Int).Set(value),
	}
}

// BlockInternalTxs represents all internal txs in a block
type BlockInternalTxs struct {
	BlockHash   common.Hash   // the block hash
	InternalTxs []*InternalTx // the internal tx collection in the block
}

// NewBlockInernalTxs creates a new BlockInernalTxs instance
func NewBlockInternalTxs(blockHash common.Hash, internalTxs []*InternalTx) *BlockInternalTxs {
	return &BlockInternalTxs{
		BlockHash:   blockHash,
		InternalTxs: internalTxs,
	}
}

// AddInternalTx adds an internalTx
func (bits *BlockInternalTxs) AddInternalTx(internalTx *InternalTx) {
	bits.InternalTxs = append(bits.InternalTxs, internalTx)
}

// Serialize serializes the internal txs.
func (bits *BlockInternalTxs) Serialize() []byte {
	return MustRlpEncode(bits.InternalTxs)
}

// Write writes the internal tx count and internal txs to database
func (bits *BlockInternalTxs) Write(batch ethdb.Batch) error {
	if err := batch.Put(MakeBlockInternalTxCountKey(bits.BlockHash), MustRlpEncode(uint(len(bits.InternalTxs)))); err != nil {
		return err
	}

	if len(bits.InternalTxs) > 0 {
		if err := batch.Put(MakeBlockInternalTxsKey(bits.BlockHash), bits.Serialize()); err != nil {
			return err
		}
	}

	return batch.Write()
}
