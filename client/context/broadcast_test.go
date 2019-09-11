package context

import (
	"fmt"
	"testing"

	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/mempool"
	"github.com/tendermint/tendermint/rpc/client/mock"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/types"
)

type MockClient struct {
	mock.Client
	err error
}

func (c MockClient) BroadcastTxCommit(tx tmtypes.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	return nil, c.err
}

func (c MockClient) BroadcastTxAsync(tx tmtypes.Tx) (*ctypes.ResultBroadcastTx, error) {
	return nil, c.err
}

func (c MockClient) BroadcastTxSync(tx tmtypes.Tx) (*ctypes.ResultBroadcastTx, error) {
	return nil, c.err
}

func CreateContextWithErrorAndMode(err error, mode string) CLIContext {
	return CLIContext{
		Client:        MockClient{err: err},
		BroadcastMode: mode,
	}
}

// Test the correct code is returned when
func TestBroadcastError(t *testing.T) {
	errors := map[error]uint32{
		mempool.ErrTxInCache:       uint32(types.CodeTxInMempoolCache),
		mempool.ErrTxTooLarge{}:    uint32(types.CodeTxTooLarge),
		mempool.ErrMempoolIsFull{}: uint32(types.CodeMempoolIsFull),
	}

	modes := []string{
		flags.BroadcastAsync,
		flags.BroadcastBlock,
		flags.BroadcastSync,
	}

	txBytes := []byte{0xA, 0xB}
	txHash := fmt.Sprintf("%X", tmhash.Sum(txBytes))

	for _, mode := range modes {
		for err, code := range errors {
			ctx := CreateContextWithErrorAndMode(err, mode)
			resp, returnedErr := ctx.BroadcastTx(txBytes)
			require.NoError(t, returnedErr)
			require.Equal(t, code, resp.Code)
			require.Equal(t, txHash, resp.TxHash)
		}
	}

}