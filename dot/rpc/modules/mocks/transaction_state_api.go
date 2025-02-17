// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	common "github.com/ChainSafe/gossamer/lib/common"
	mock "github.com/stretchr/testify/mock"

	transaction "github.com/ChainSafe/gossamer/lib/transaction"

	types "github.com/ChainSafe/gossamer/dot/types"
)

// TransactionStateAPI is an autogenerated mock type for the TransactionStateAPI type
type TransactionStateAPI struct {
	mock.Mock
}

// AddToPool provides a mock function with given fields: _a0
func (_m *TransactionStateAPI) AddToPool(_a0 *transaction.ValidTransaction) common.Hash {
	ret := _m.Called(_a0)

	var r0 common.Hash
	if rf, ok := ret.Get(0).(func(*transaction.ValidTransaction) common.Hash); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Hash)
		}
	}

	return r0
}

// FreeStatusNotifierChannel provides a mock function with given fields: ch
func (_m *TransactionStateAPI) FreeStatusNotifierChannel(ch chan transaction.Status) {
	_m.Called(ch)
}

// GetStatusNotifierChannel provides a mock function with given fields: ext
func (_m *TransactionStateAPI) GetStatusNotifierChannel(ext types.Extrinsic) chan transaction.Status {
	ret := _m.Called(ext)

	var r0 chan transaction.Status
	if rf, ok := ret.Get(0).(func(types.Extrinsic) chan transaction.Status); ok {
		r0 = rf(ext)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan transaction.Status)
		}
	}

	return r0
}

// Peek provides a mock function with given fields:
func (_m *TransactionStateAPI) Peek() *transaction.ValidTransaction {
	ret := _m.Called()

	var r0 *transaction.ValidTransaction
	if rf, ok := ret.Get(0).(func() *transaction.ValidTransaction); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transaction.ValidTransaction)
		}
	}

	return r0
}

// Pending provides a mock function with given fields:
func (_m *TransactionStateAPI) Pending() []*transaction.ValidTransaction {
	ret := _m.Called()

	var r0 []*transaction.ValidTransaction
	if rf, ok := ret.Get(0).(func() []*transaction.ValidTransaction); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*transaction.ValidTransaction)
		}
	}

	return r0
}

// Pop provides a mock function with given fields:
func (_m *TransactionStateAPI) Pop() *transaction.ValidTransaction {
	ret := _m.Called()

	var r0 *transaction.ValidTransaction
	if rf, ok := ret.Get(0).(func() *transaction.ValidTransaction); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transaction.ValidTransaction)
		}
	}

	return r0
}
