// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	common "github.com/ChainSafe/gossamer/lib/common"
	mock "github.com/stretchr/testify/mock"

	runtime "github.com/ChainSafe/gossamer/lib/runtime"

	types "github.com/ChainSafe/gossamer/dot/types"
)

// BlockAPI is an autogenerated mock type for the BlockAPI type
type BlockAPI struct {
	mock.Mock
}

// BestBlockHash provides a mock function with given fields:
func (_m *BlockAPI) BestBlockHash() common.Hash {
	ret := _m.Called()

	var r0 common.Hash
	if rf, ok := ret.Get(0).(func() common.Hash); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Hash)
		}
	}

	return r0
}

// FreeFinalisedNotifierChannel provides a mock function with given fields: ch
func (_m *BlockAPI) FreeFinalisedNotifierChannel(ch chan *types.FinalisationInfo) {
	_m.Called(ch)
}

// FreeImportedBlockNotifierChannel provides a mock function with given fields: ch
func (_m *BlockAPI) FreeImportedBlockNotifierChannel(ch chan *types.Block) {
	_m.Called(ch)
}

// GetBlockByHash provides a mock function with given fields: hash
func (_m *BlockAPI) GetBlockByHash(hash common.Hash) (*types.Block, error) {
	ret := _m.Called(hash)

	var r0 *types.Block
	if rf, ok := ret.Get(0).(func(common.Hash) *types.Block); ok {
		r0 = rf(hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Block)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Hash) error); ok {
		r1 = rf(hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFinalisedHash provides a mock function with given fields: _a0, _a1
func (_m *BlockAPI) GetFinalisedHash(_a0 uint64, _a1 uint64) (common.Hash, error) {
	ret := _m.Called(_a0, _a1)

	var r0 common.Hash
	if rf, ok := ret.Get(0).(func(uint64, uint64) common.Hash); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, uint64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFinalisedNotifierChannel provides a mock function with given fields:
func (_m *BlockAPI) GetFinalisedNotifierChannel() chan *types.FinalisationInfo {
	ret := _m.Called()

	var r0 chan *types.FinalisationInfo
	if rf, ok := ret.Get(0).(func() chan *types.FinalisationInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan *types.FinalisationInfo)
		}
	}

	return r0
}

// GetHashByNumber provides a mock function with given fields: blockNumber
func (_m *BlockAPI) GetHashByNumber(blockNumber uint) (common.Hash, error) {
	ret := _m.Called(blockNumber)

	var r0 common.Hash
	if rf, ok := ret.Get(0).(func(uint) common.Hash); ok {
		r0 = rf(blockNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHeader provides a mock function with given fields: hash
func (_m *BlockAPI) GetHeader(hash common.Hash) (*types.Header, error) {
	ret := _m.Called(hash)

	var r0 *types.Header
	if rf, ok := ret.Get(0).(func(common.Hash) *types.Header); ok {
		r0 = rf(hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Header)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Hash) error); ok {
		r1 = rf(hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHighestFinalisedHash provides a mock function with given fields:
func (_m *BlockAPI) GetHighestFinalisedHash() (common.Hash, error) {
	ret := _m.Called()

	var r0 common.Hash
	if rf, ok := ret.Get(0).(func() common.Hash); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetImportedBlockNotifierChannel provides a mock function with given fields:
func (_m *BlockAPI) GetImportedBlockNotifierChannel() chan *types.Block {
	ret := _m.Called()

	var r0 chan *types.Block
	if rf, ok := ret.Get(0).(func() chan *types.Block); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan *types.Block)
		}
	}

	return r0
}

// GetJustification provides a mock function with given fields: hash
func (_m *BlockAPI) GetJustification(hash common.Hash) ([]byte, error) {
	ret := _m.Called(hash)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(common.Hash) []byte); ok {
		r0 = rf(hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Hash) error); ok {
		r1 = rf(hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRuntime provides a mock function with given fields: hash
func (_m *BlockAPI) GetRuntime(hash *common.Hash) (runtime.Instance, error) {
	ret := _m.Called(hash)

	var r0 runtime.Instance
	if rf, ok := ret.Get(0).(func(*common.Hash) runtime.Instance); ok {
		r0 = rf(hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(runtime.Instance)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*common.Hash) error); ok {
		r1 = rf(hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HasJustification provides a mock function with given fields: hash
func (_m *BlockAPI) HasJustification(hash common.Hash) (bool, error) {
	ret := _m.Called(hash)

	var r0 bool
	if rf, ok := ret.Get(0).(func(common.Hash) bool); ok {
		r0 = rf(hash)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Hash) error); ok {
		r1 = rf(hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterRuntimeUpdatedChannel provides a mock function with given fields: ch
func (_m *BlockAPI) RegisterRuntimeUpdatedChannel(ch chan<- runtime.Version) (uint32, error) {
	ret := _m.Called(ch)

	var r0 uint32
	if rf, ok := ret.Get(0).(func(chan<- runtime.Version) uint32); ok {
		r0 = rf(ch)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(chan<- runtime.Version) error); ok {
		r1 = rf(ch)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubChain provides a mock function with given fields: start, end
func (_m *BlockAPI) SubChain(start common.Hash, end common.Hash) ([]common.Hash, error) {
	ret := _m.Called(start, end)

	var r0 []common.Hash
	if rf, ok := ret.Get(0).(func(common.Hash, common.Hash) []common.Hash); ok {
		r0 = rf(start, end)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]common.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Hash, common.Hash) error); ok {
		r1 = rf(start, end)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnregisterRuntimeUpdatedChannel provides a mock function with given fields: id
func (_m *BlockAPI) UnregisterRuntimeUpdatedChannel(id uint32) bool {
	ret := _m.Called(id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uint32) bool); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
