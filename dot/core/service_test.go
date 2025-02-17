// Copyright 2022 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package core

import (
	"context"
	"errors"
	"testing"

	"github.com/ChainSafe/gossamer/dot/network"
	"github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/blocktree"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/keystore"
	"github.com/ChainSafe/gossamer/lib/runtime"
	mocksruntime "github.com/ChainSafe/gossamer/lib/runtime/mocks"
	rtstorage "github.com/ChainSafe/gossamer/lib/runtime/storage"
	"github.com/ChainSafe/gossamer/lib/runtime/wasmer"
	"github.com/ChainSafe/gossamer/lib/transaction"
	"github.com/ChainSafe/gossamer/lib/trie"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errTestDummyError = errors.New("test dummy error")

func Test_Service_StorageRoot(t *testing.T) {
	t.Parallel()
	emptyTrie := trie.NewEmptyTrie()
	ts, err := rtstorage.NewTrieState(emptyTrie)
	require.NoError(t, err)
	tests := []struct {
		name          string
		service       *Service
		exp           common.Hash
		retTrieState  *rtstorage.TrieState
		trieStateCall bool
		retErr        error
		expErr        error
		expErrMsg     string
	}{
		{
			name:      "nil storage state",
			service:   &Service{},
			expErr:    ErrNilStorageState,
			expErrMsg: ErrNilStorageState.Error(),
		},
		{
			name:          "storage trie state error",
			service:       &Service{},
			retErr:        errTestDummyError,
			expErr:        errTestDummyError,
			expErrMsg:     errTestDummyError.Error(),
			trieStateCall: true,
		},
		{
			name:    "storage trie state ok",
			service: &Service{},
			exp: common.Hash{0x3, 0x17, 0xa, 0x2e, 0x75, 0x97, 0xb7, 0xb7, 0xe3, 0xd8, 0x4c, 0x5, 0x39, 0x1d, 0x13, 0x9a,
				0x62, 0xb1, 0x57, 0xe7, 0x87, 0x86, 0xd8, 0xc0, 0x82, 0xf2, 0x9d, 0xcf, 0x4c, 0x11, 0x13, 0x14},
			retTrieState:  ts,
			trieStateCall: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := tt.service
			if tt.trieStateCall {
				ctrl := gomock.NewController(t)
				mockStorageState := NewMockStorageState(ctrl)
				mockStorageState.EXPECT().TrieState(nil).Return(tt.retTrieState, tt.retErr)
				s.storageState = mockStorageState
			}

			res, err := s.StorageRoot()
			assert.ErrorIs(t, err, tt.expErr)
			if tt.expErr != nil {
				assert.EqualError(t, err, tt.expErrMsg)
			}
			assert.Equal(t, tt.exp, res)
		})
	}
}

func Test_Service_handleCodeSubstitution(t *testing.T) {
	t.Parallel()
	newTestInstance := func(code []byte, cfg *wasmer.Config) (*wasmer.Instance, error) {
		return &wasmer.Instance{}, nil
	}

	execTest := func(t *testing.T, s *Service, blockHash common.Hash, expErr error) {
		err := s.handleCodeSubstitution(blockHash, nil, newTestInstance)
		assert.ErrorIs(t, err, expErr)
		if expErr != nil {
			assert.EqualError(t, err, errTestDummyError.Error())
		}
	}
	testRuntime := []byte{21}
	t.Run("nil value", func(t *testing.T) {
		t.Parallel()
		s := &Service{codeSubstitute: map[common.Hash]string{}}
		err := s.handleCodeSubstitution(common.Hash{}, nil, newTestInstance)
		assert.NoError(t, err)
	})

	t.Run("getRuntime error", func(t *testing.T) {
		t.Parallel()
		// hash for known test code substitution
		blockHash := common.MustHexToHash("0x86aa36a140dfc449c30dbce16ce0fea33d5c3786766baa764e33f336841b9e29")
		testCodeSubstitute := map[common.Hash]string{
			blockHash: common.BytesToHex(testRuntime),
		}

		ctrl := gomock.NewController(t)
		mockBlockState := NewMockBlockState(ctrl)
		mockBlockState.EXPECT().GetRuntime(&blockHash).Return(nil, errTestDummyError)
		s := &Service{
			codeSubstitute: testCodeSubstitute,
			blockState:     mockBlockState,
		}
		execTest(t, s, blockHash, errTestDummyError)
	})

	t.Run("code substitute error", func(t *testing.T) {
		t.Parallel()
		// hash for known test code substitution
		blockHash := common.MustHexToHash("0x86aa36a140dfc449c30dbce16ce0fea33d5c3786766baa764e33f336841b9e29")
		testCodeSubstitute := map[common.Hash]string{
			blockHash: common.BytesToHex(testRuntime),
		}

		runtimeMock := new(mocksruntime.Instance)
		runtimeMock.On("Keystore").Return(&keystore.GlobalKeystore{})
		runtimeMock.On("NodeStorage").Return(runtime.NodeStorage{})
		runtimeMock.On("NetworkService").Return(new(runtime.TestRuntimeNetwork))
		runtimeMock.On("Validator").Return(true)

		ctrl := gomock.NewController(t)
		mockBlockState := NewMockBlockState(ctrl)
		mockBlockState.EXPECT().GetRuntime(&blockHash).Return(runtimeMock, nil)
		mockCodeSubState := NewMockCodeSubstitutedState(ctrl)
		mockCodeSubState.EXPECT().StoreCodeSubstitutedBlockHash(blockHash).Return(errTestDummyError)
		s := &Service{
			codeSubstitute:       testCodeSubstitute,
			blockState:           mockBlockState,
			codeSubstitutedState: mockCodeSubState,
		}
		execTest(t, s, blockHash, errTestDummyError)
	})

	t.Run("happyPath", func(t *testing.T) {
		t.Parallel()
		// hash for known test code substitution
		blockHash := common.MustHexToHash("0x86aa36a140dfc449c30dbce16ce0fea33d5c3786766baa764e33f336841b9e29")
		testCodeSubstitute := map[common.Hash]string{
			blockHash: common.BytesToHex(testRuntime),
		}

		runtimeMock := new(mocksruntime.Instance)
		runtimeMock.On("Keystore").Return(&keystore.GlobalKeystore{})
		runtimeMock.On("NodeStorage").Return(runtime.NodeStorage{})
		runtimeMock.On("NetworkService").Return(new(runtime.TestRuntimeNetwork))
		runtimeMock.On("Validator").Return(true)

		ctrl := gomock.NewController(t)
		mockBlockState := NewMockBlockState(ctrl)
		mockBlockState.EXPECT().GetRuntime(&blockHash).Return(runtimeMock, nil)
		mockBlockState.EXPECT().StoreRuntime(blockHash, gomock.Any())
		mockCodeSubState := NewMockCodeSubstitutedState(ctrl)
		mockCodeSubState.EXPECT().StoreCodeSubstitutedBlockHash(blockHash).Return(nil)
		s := &Service{
			codeSubstitute:       testCodeSubstitute,
			blockState:           mockBlockState,
			codeSubstitutedState: mockCodeSubState,
		}
		err := s.handleCodeSubstitution(blockHash, nil, newTestInstance)
		assert.NoError(t, err)
	})
}

func Test_Service_handleBlock(t *testing.T) {
	t.Parallel()
	execTest := func(t *testing.T, s *Service, block *types.Block, trieState *rtstorage.TrieState, expErr error) {
		err := s.handleBlock(block, trieState)
		assert.ErrorIs(t, err, expErr)
		if expErr != nil {
			assert.EqualError(t, err, expErr.Error())
		}
	}
	t.Run("nil input", func(t *testing.T) {
		t.Parallel()
		s := &Service{}
		execTest(t, s, nil, nil, ErrNilBlockHandlerParameter)
	})

	t.Run("storeTrie error", func(t *testing.T) {
		t.Parallel()
		emptyTrie := trie.NewEmptyTrie()
		trieState, err := rtstorage.NewTrieState(emptyTrie)
		require.NoError(t, err)

		testHeader := types.NewEmptyHeader()
		block := types.NewBlock(*testHeader, *types.NewBody([]types.Extrinsic{[]byte{21}}))
		block.Header.Number = 21

		ctrl := gomock.NewController(t)
		mockStorageState := NewMockStorageState(ctrl)
		mockStorageState.EXPECT().StoreTrie(trieState, &block.Header).Return(errTestDummyError)

		s := &Service{storageState: mockStorageState}
		execTest(t, s, &block, trieState, errTestDummyError)
	})

	t.Run("addBlock quit error", func(t *testing.T) {
		t.Parallel()
		emptyTrie := trie.NewEmptyTrie()
		trieState, err := rtstorage.NewTrieState(emptyTrie)
		require.NoError(t, err)

		testHeader := types.NewEmptyHeader()
		block := types.NewBlock(*testHeader, *types.NewBody([]types.Extrinsic{[]byte{21}}))
		block.Header.Number = 21

		ctrl := gomock.NewController(t)
		mockStorageState := NewMockStorageState(ctrl)
		mockStorageState.EXPECT().StoreTrie(trieState, &block.Header).Return(nil)
		mockBlockState := NewMockBlockState(ctrl)
		mockBlockState.EXPECT().AddBlock(&block).Return(errTestDummyError)

		s := &Service{
			storageState: mockStorageState,
			blockState:   mockBlockState,
		}
		execTest(t, s, &block, trieState, errTestDummyError)
	})

	t.Run("addBlock parent not found error", func(t *testing.T) {
		t.Parallel()
		emptyTrie := trie.NewEmptyTrie()
		trieState, err := rtstorage.NewTrieState(emptyTrie)
		require.NoError(t, err)

		testHeader := types.NewEmptyHeader()
		block := types.NewBlock(*testHeader, *types.NewBody([]types.Extrinsic{[]byte{21}}))
		block.Header.Number = 21

		ctrl := gomock.NewController(t)
		mockStorageState := NewMockStorageState(ctrl)
		mockStorageState.EXPECT().StoreTrie(trieState, &block.Header).Return(nil)
		mockBlockState := NewMockBlockState(ctrl)
		mockBlockState.EXPECT().AddBlock(&block).Return(blocktree.ErrParentNotFound)

		s := &Service{
			storageState: mockStorageState,
			blockState:   mockBlockState,
		}
		execTest(t, s, &block, trieState, blocktree.ErrParentNotFound)
	})

	t.Run("addBlock error continue", func(t *testing.T) {
		t.Parallel()
		emptyTrie := trie.NewEmptyTrie()
		trieState, err := rtstorage.NewTrieState(emptyTrie)
		require.NoError(t, err)

		testHeader := types.NewEmptyHeader()
		block := types.NewBlock(*testHeader, *types.NewBody([]types.Extrinsic{[]byte{21}}))
		block.Header.Number = 21

		ctrl := gomock.NewController(t)
		mockStorageState := NewMockStorageState(ctrl)
		mockStorageState.EXPECT().StoreTrie(trieState, &block.Header).Return(nil)
		mockBlockState := NewMockBlockState(ctrl)
		mockBlockState.EXPECT().AddBlock(&block).Return(blocktree.ErrBlockExists)
		mockBlockState.EXPECT().GetRuntime(&block.Header.ParentHash).Return(nil, errTestDummyError)
		mockDigestHandler := NewMockDigestHandler(ctrl)
		mockDigestHandler.EXPECT().HandleDigests(&block.Header)

		s := &Service{
			storageState:  mockStorageState,
			blockState:    mockBlockState,
			digestHandler: mockDigestHandler,
		}
		execTest(t, s, &block, trieState, errTestDummyError)
	})

	t.Run("handle runtime changes error", func(t *testing.T) {
		t.Parallel()
		emptyTrie := trie.NewEmptyTrie()
		trieState, err := rtstorage.NewTrieState(emptyTrie)
		require.NoError(t, err)

		testHeader := types.NewEmptyHeader()
		block := types.NewBlock(*testHeader, *types.NewBody([]types.Extrinsic{[]byte{21}}))
		block.Header.Number = 21

		ctrl := gomock.NewController(t)
		runtimeMock := new(mocksruntime.Instance)
		mockStorageState := NewMockStorageState(ctrl)
		mockStorageState.EXPECT().StoreTrie(trieState, &block.Header).Return(nil)
		mockBlockState := NewMockBlockState(ctrl)
		mockBlockState.EXPECT().AddBlock(&block).Return(blocktree.ErrBlockExists)
		mockBlockState.EXPECT().GetRuntime(&block.Header.ParentHash).Return(runtimeMock, nil)
		mockBlockState.EXPECT().HandleRuntimeChanges(trieState, runtimeMock, block.Header.Hash()).
			Return(errTestDummyError)
		mockDigestHandler := NewMockDigestHandler(ctrl)
		mockDigestHandler.EXPECT().HandleDigests(&block.Header)

		s := &Service{
			storageState:  mockStorageState,
			blockState:    mockBlockState,
			digestHandler: mockDigestHandler,
		}
		execTest(t, s, &block, trieState, errTestDummyError)
	})

	t.Run("code substitution ok", func(t *testing.T) {
		t.Parallel()
		emptyTrie := trie.NewEmptyTrie()
		trieState, err := rtstorage.NewTrieState(emptyTrie)
		require.NoError(t, err)

		testHeader := types.NewEmptyHeader()
		block := types.NewBlock(*testHeader, *types.NewBody([]types.Extrinsic{[]byte{21}}))
		block.Header.Number = 21

		ctrl := gomock.NewController(t)
		runtimeMock := new(mocksruntime.Instance)
		mockStorageState := NewMockStorageState(ctrl)
		mockStorageState.EXPECT().StoreTrie(trieState, &block.Header).Return(nil)
		mockBlockState := NewMockBlockState(ctrl)
		mockBlockState.EXPECT().AddBlock(&block).Return(blocktree.ErrBlockExists)
		mockBlockState.EXPECT().GetRuntime(&block.Header.ParentHash).Return(runtimeMock, nil)
		mockBlockState.EXPECT().HandleRuntimeChanges(trieState, runtimeMock, block.Header.Hash()).Return(nil)
		mockDigestHandler := NewMockDigestHandler(ctrl)
		mockDigestHandler.EXPECT().HandleDigests(&block.Header)

		s := &Service{
			storageState:  mockStorageState,
			blockState:    mockBlockState,
			digestHandler: mockDigestHandler,
			ctx:           context.Background(),
		}
		execTest(t, s, &block, trieState, nil)
	})
}

func Test_Service_HandleBlockProduced(t *testing.T) {
	t.Parallel()
	execTest := func(t *testing.T, s *Service, block *types.Block, trieState *rtstorage.TrieState, expErr error) {
		err := s.HandleBlockProduced(block, trieState)
		assert.ErrorIs(t, err, expErr)
		if expErr != nil {
			assert.EqualError(t, err, expErr.Error())
		}
	}
	t.Run("nil input", func(t *testing.T) {
		t.Parallel()
		s := &Service{}
		execTest(t, s, nil, nil, ErrNilBlockHandlerParameter)
	})

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		emptyTrie := trie.NewEmptyTrie()
		trieState, err := rtstorage.NewTrieState(emptyTrie)
		require.NoError(t, err)

		digest := types.NewDigest()
		err = digest.Add(
			types.PreRuntimeDigest{
				ConsensusEngineID: types.BabeEngineID,
				Data:              common.MustHexToBytes("0x0201000000ef55a50f00000000"),
			})
		require.NoError(t, err)

		testHeader := types.NewEmptyHeader()
		block := types.NewBlock(*testHeader, *types.NewBody([]types.Extrinsic{[]byte{21}}))
		block.Header.Number = 21
		block.Header.Digest = digest
		msg := &network.BlockAnnounceMessage{
			ParentHash:     block.Header.ParentHash,
			Number:         block.Header.Number,
			StateRoot:      block.Header.StateRoot,
			ExtrinsicsRoot: block.Header.ExtrinsicsRoot,
			Digest:         digest,
			BestBlock:      true,
		}

		ctrl := gomock.NewController(t)
		runtimeMock := new(mocksruntime.Instance)
		mockStorageState := NewMockStorageState(ctrl)
		mockStorageState.EXPECT().StoreTrie(trieState, &block.Header).Return(nil)
		mockBlockState := NewMockBlockState(ctrl)
		mockBlockState.EXPECT().AddBlock(&block).Return(blocktree.ErrBlockExists)
		mockBlockState.EXPECT().GetRuntime(&block.Header.ParentHash).Return(runtimeMock, nil)
		mockBlockState.EXPECT().HandleRuntimeChanges(trieState, runtimeMock, block.Header.Hash()).Return(nil)
		mockDigestHandler := NewMockDigestHandler(ctrl)
		mockDigestHandler.EXPECT().HandleDigests(&block.Header)
		mockNetwork := NewMockNetwork(ctrl)
		mockNetwork.EXPECT().GossipMessage(msg)

		s := &Service{
			storageState:  mockStorageState,
			blockState:    mockBlockState,
			digestHandler: mockDigestHandler,
			net:           mockNetwork,
			ctx:           context.Background(),
		}
		execTest(t, s, &block, trieState, nil)
	})
}

func Test_Service_maintainTransactionPool(t *testing.T) {
	t.Parallel()
	t.Run("Validate Transaction err", func(t *testing.T) {
		t.Parallel()
		testHeader := types.NewEmptyHeader()
		block := types.NewBlock(*testHeader, *types.NewBody([]types.Extrinsic{[]byte{21}}))
		block.Header.Number = 21

		validity := &transaction.Validity{
			Priority: 0x3e8,
			Requires: [][]byte{{0xb5, 0x47, 0xb1, 0x90, 0x37, 0x10, 0x7e, 0x1f, 0x79,
				0x4c, 0xa8, 0x69, 0x0, 0xa1, 0xb5, 0x98}},
			Provides: [][]byte{{0xe4, 0x80, 0x7d, 0x1b, 0x67, 0x49, 0x37, 0xbf, 0xc7,
				0x89, 0xbb, 0xdd, 0x88, 0x6a, 0xdd, 0xd6}},
			Longevity: 0x40,
			Propagate: true,
		}

		extrinsic := types.Extrinsic{21}
		vt := transaction.NewValidTransaction(extrinsic, validity)

		ctrl := gomock.NewController(t)
		runtimeMock := new(mocksruntime.Instance)
		runtimeMock.On("ValidateTransaction", types.Extrinsic{21}).Return(nil, errTestDummyError)
		mockTxnState := NewMockTransactionState(ctrl)
		mockTxnState.EXPECT().RemoveExtrinsic(types.Extrinsic{21}).Times(2)
		mockTxnState.EXPECT().PendingInPool().Return([]*transaction.ValidTransaction{vt})
		mockBlockState := NewMockBlockState(ctrl)
		mockBlockState.EXPECT().GetRuntime(nil).Return(runtimeMock, nil)
		s := &Service{
			transactionState: mockTxnState,
			blockState:       mockBlockState,
		}
		s.maintainTransactionPool(&block)
	})

	t.Run("Validate Transaction ok", func(t *testing.T) {
		t.Parallel()
		testHeader := types.NewEmptyHeader()
		block := types.NewBlock(*testHeader, *types.NewBody([]types.Extrinsic{[]byte{21}}))
		block.Header.Number = 21

		validity := &transaction.Validity{
			Priority: 0x3e8,
			Requires: [][]byte{{0xb5, 0x47, 0xb1, 0x90, 0x37, 0x10, 0x7e, 0x1f, 0x79, 0x4c,
				0xa8, 0x69, 0x0, 0xa1, 0xb5, 0x98}},
			Provides: [][]byte{{0xe4, 0x80, 0x7d, 0x1b, 0x67, 0x49, 0x37, 0xbf, 0xc7, 0x89,
				0xbb, 0xdd, 0x88, 0x6a, 0xdd, 0xd6}},
			Longevity: 0x40,
			Propagate: true,
		}

		extrinsic := types.Extrinsic{21}
		vt := transaction.NewValidTransaction(extrinsic, validity)
		tx := transaction.NewValidTransaction(types.Extrinsic{21}, &transaction.Validity{Propagate: true})

		ctrl := gomock.NewController(t)
		runtimeMock := new(mocksruntime.Instance)
		runtimeMock.On("ValidateTransaction", types.Extrinsic{21}).
			Return(&transaction.Validity{Propagate: true}, nil)
		mockTxnState := NewMockTransactionState(ctrl)
		mockTxnState.EXPECT().RemoveExtrinsic(types.Extrinsic{21})
		mockTxnState.EXPECT().PendingInPool().Return([]*transaction.ValidTransaction{vt})
		mockTxnState.EXPECT().Push(tx).Return(common.Hash{}, nil)
		mockTxnState.EXPECT().RemoveExtrinsicFromPool(types.Extrinsic{21})
		mockBlockStateOk := NewMockBlockState(ctrl)
		mockBlockStateOk.EXPECT().GetRuntime(nil).Return(runtimeMock, nil)
		s := &Service{
			transactionState: mockTxnState,
			blockState:       mockBlockStateOk,
		}
		s.maintainTransactionPool(&block)
	})
}

func Test_Service_handleBlocksAsync(t *testing.T) {
	t.Parallel()
	t.Run("cancelled context", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockBlockState := NewMockBlockState(ctrl)
		mockBlockState.EXPECT().BestBlockHash().Return(common.Hash{})
		blockAddChan := make(chan *types.Block)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		s := &Service{
			blockState: mockBlockState,
			blockAddCh: blockAddChan,
			ctx:        ctx,
		}
		s.handleBlocksAsync()
	})

	t.Run("channel not ok", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockBlockState := NewMockBlockState(ctrl)
		mockBlockState.EXPECT().BestBlockHash().Return(common.Hash{})
		blockAddChan := make(chan *types.Block)
		close(blockAddChan)
		s := &Service{
			blockState: mockBlockState,
			blockAddCh: blockAddChan,
			ctx:        context.Background(),
		}
		s.handleBlocksAsync()
	})

	t.Run("nil block", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockBlockState := NewMockBlockState(ctrl)
		mockBlockState.EXPECT().BestBlockHash().Return(common.Hash{}).Times(2)
		blockAddChan := make(chan *types.Block)
		go func() {
			blockAddChan <- nil
			close(blockAddChan)
		}()
		s := &Service{
			blockState: mockBlockState,
			blockAddCh: blockAddChan,
			ctx:        context.Background(),
		}
		s.handleBlocksAsync()
	})

	t.Run("handleChainReorg error", func(t *testing.T) {
		t.Parallel()
		validity := &transaction.Validity{
			Priority: 0x3e8,
			Requires: [][]byte{{0xb5, 0x47, 0xb1, 0x90, 0x37, 0x10, 0x7e, 0x1f, 0x79, 0x4c,
				0xa8, 0x69, 0x0, 0xa1, 0xb5, 0x98}},
			Provides: [][]byte{{0xe4, 0x80, 0x7d, 0x1b, 0x67, 0x49, 0x37, 0xbf, 0xc7, 0x89,
				0xbb, 0xdd, 0x88, 0x6a, 0xdd, 0xd6}},
			Longevity: 0x40,
			Propagate: true,
		}

		extrinsic := types.Extrinsic{21}
		vt := transaction.NewValidTransaction(extrinsic, validity)

		testHeader := types.NewEmptyHeader()
		block := types.NewBlock(*testHeader, *types.NewBody([]types.Extrinsic{[]byte{21}}))
		block.Header.Number = 21

		ctrl := gomock.NewController(t)
		runtimeMock := new(mocksruntime.Instance)
		runtimeMock.On("ValidateTransaction", types.Extrinsic{21}).Return(nil, errTestDummyError)
		mockBlockState := NewMockBlockState(ctrl)
		mockBlockState.EXPECT().BestBlockHash().Return(common.Hash{}).Times(2)
		mockBlockState.EXPECT().HighestCommonAncestor(common.Hash{}, block.Header.Hash()).
			Return(common.Hash{}, errTestDummyError)
		mockBlockState.EXPECT().GetRuntime(nil).Return(runtimeMock, nil)
		mockTxnStateErr := NewMockTransactionState(ctrl)
		mockTxnStateErr.EXPECT().RemoveExtrinsic(types.Extrinsic{21}).Times(2)
		mockTxnStateErr.EXPECT().PendingInPool().Return([]*transaction.ValidTransaction{vt})
		blockAddChan := make(chan *types.Block)
		go func() {
			blockAddChan <- &block
			close(blockAddChan)
		}()
		s := &Service{
			blockState:       mockBlockState,
			transactionState: mockTxnStateErr,
			blockAddCh:       blockAddChan,
			ctx:              context.Background(),
		}
		s.handleBlocksAsync()
	})
}
