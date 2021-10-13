// Copyright 2020 ChainSafe Systems (ON) Corp.
// This file is part of gossamer.
//
// The gossamer library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The gossamer library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the gossamer library. If not, see <http://www.gnu.org/licenses/>.

package modules

import (
	"github.com/ChainSafe/gossamer/dot/state"
	"github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/genesis"
	"github.com/ChainSafe/gossamer/lib/runtime"
	"github.com/ChainSafe/gossamer/lib/trie"
	log "github.com/ChainSafe/log15"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/big"
	"testing"
)

// test data
var (
	sampleBodyBytes = *types.NewBody([]types.Extrinsic{[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}})
	// sampleBodyString is string conversion of sampleBodyBytes
	sampleBodyString = []string{"0x2800010203040506070809"}
)

func newTestStateService(t *testing.T) *state.Service {
	testDatadirPath, err := ioutil.TempDir("/tmp", "test-datadir-*")
	require.NoError(t, err)

	config := state.Config{
		Path:     testDatadirPath,
		LogLevel: log.LvlInfo,
	}

	// Mock this
	stateSrvc := state.NewService(config)
	stateSrvc.UseMemDB()

	gen, genTrie, genesisHeader := genesis.NewTestGenesisWithTrieAndHeader(t)

	err = stateSrvc.Initialise(gen, genesisHeader, genTrie)
	require.NoError(t, err)

	err = stateSrvc.Start()
	require.NoError(t, err)

	rt, err := stateSrvc.CreateGenesisRuntime(genTrie, gen)
	require.NoError(t, err)

	err = loadTestBlocks(t, genesisHeader.Hash(), stateSrvc.Block, rt)
	require.NoError(t, err)

	t.Cleanup(func() {
		stateSrvc.Stop()
	})
	return stateSrvc
}

func loadTestBlocks(t *testing.T, gh common.Hash, bs *state.BlockState, rt runtime.Instance) error {
	// Create header
	header0 := &types.Header{
		Number:     big.NewInt(0),
		Digest:     types.NewDigest(),
		ParentHash: gh,
		StateRoot:  trie.EmptyHash,
	}
	// Create blockHash
	blockHash0 := header0.Hash()
	block0 := &types.Block{
		Header: *header0,
		Body:   sampleBodyBytes,
	}

	err := bs.AddBlock(block0)
	if err != nil {
		return err
	}

	bs.StoreRuntime(block0.Header.Hash(), rt)

	// Create header & blockData for block 1
	digest := types.NewDigest()
	err = digest.Add(*types.NewBabeSecondaryPlainPreDigest(0, 1).ToPreRuntimeDigest())
	require.NoError(t, err)
	header1 := &types.Header{
		Number:     big.NewInt(1),
		Digest:     digest,
		ParentHash: blockHash0,
		StateRoot:  trie.EmptyHash,
	}

	block1 := &types.Block{
		Header: *header1,
		Body:   sampleBodyBytes,
	}

	// Add the block1 to the DB
	err = bs.AddBlock(block1)
	if err != nil {
		return err
	}

	bs.StoreRuntime(block1.Header.Hash(), rt)

	return nil
}
