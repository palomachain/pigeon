// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	big "math/big"

	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	compass "github.com/palomachain/pigeon/chain/evm/abi/compass"

	event "github.com/ethereum/go-ethereum/event"

	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// CompassBinding is an autogenerated mock type for the CompassBinding type
type CompassBinding struct {
	mock.Mock
}

// FilterLogicCallEvent provides a mock function with given fields: opts
func (_m *CompassBinding) FilterLogicCallEvent(opts *bind.FilterOpts) (*compass.CompassLogicCallEventIterator, error) {
	ret := _m.Called(opts)

	var r0 *compass.CompassLogicCallEventIterator
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.FilterOpts) (*compass.CompassLogicCallEventIterator, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.FilterOpts) *compass.CompassLogicCallEventIterator); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*compass.CompassLogicCallEventIterator)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.FilterOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FilterValsetUpdated provides a mock function with given fields: opts
func (_m *CompassBinding) FilterValsetUpdated(opts *bind.FilterOpts) (*compass.CompassValsetUpdatedIterator, error) {
	ret := _m.Called(opts)

	var r0 *compass.CompassValsetUpdatedIterator
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.FilterOpts) (*compass.CompassValsetUpdatedIterator, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.FilterOpts) *compass.CompassValsetUpdatedIterator); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*compass.CompassValsetUpdatedIterator)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.FilterOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LastCheckpoint provides a mock function with given fields: opts
func (_m *CompassBinding) LastCheckpoint(opts *bind.CallOpts) ([32]byte, error) {
	ret := _m.Called(opts)

	var r0 [32]byte
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) ([32]byte, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) [32]byte); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([32]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LastValsetId provides a mock function with given fields: opts
func (_m *CompassBinding) LastValsetId(opts *bind.CallOpts) (*big.Int, error) {
	ret := _m.Called(opts)

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) (*big.Int, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) *big.Int); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MessageIdUsed provides a mock function with given fields: opts, arg0
func (_m *CompassBinding) MessageIdUsed(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	ret := _m.Called(opts, arg0)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, *big.Int) (bool, error)); ok {
		return rf(opts, arg0)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, *big.Int) bool); ok {
		r0 = rf(opts, arg0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts, *big.Int) error); ok {
		r1 = rf(opts, arg0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseLogicCallEvent provides a mock function with given fields: log
func (_m *CompassBinding) ParseLogicCallEvent(log types.Log) (*compass.CompassLogicCallEvent, error) {
	ret := _m.Called(log)

	var r0 *compass.CompassLogicCallEvent
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Log) (*compass.CompassLogicCallEvent, error)); ok {
		return rf(log)
	}
	if rf, ok := ret.Get(0).(func(types.Log) *compass.CompassLogicCallEvent); ok {
		r0 = rf(log)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*compass.CompassLogicCallEvent)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Log) error); ok {
		r1 = rf(log)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseValsetUpdated provides a mock function with given fields: log
func (_m *CompassBinding) ParseValsetUpdated(log types.Log) (*compass.CompassValsetUpdated, error) {
	ret := _m.Called(log)

	var r0 *compass.CompassValsetUpdated
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Log) (*compass.CompassValsetUpdated, error)); ok {
		return rf(log)
	}
	if rf, ok := ret.Get(0).(func(types.Log) *compass.CompassValsetUpdated); ok {
		r0 = rf(log)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*compass.CompassValsetUpdated)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Log) error); ok {
		r1 = rf(log)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubmitLogicCall provides a mock function with given fields: opts, consensus, args, messageId, deadline
func (_m *CompassBinding) SubmitLogicCall(opts *bind.TransactOpts, consensus compass.Struct2, args compass.Struct3, messageId *big.Int, deadline *big.Int) (*types.Transaction, error) {
	ret := _m.Called(opts, consensus, args, messageId, deadline)

	var r0 *types.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, compass.Struct2, compass.Struct3, *big.Int, *big.Int) (*types.Transaction, error)); ok {
		return rf(opts, consensus, args, messageId, deadline)
	}
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, compass.Struct2, compass.Struct3, *big.Int, *big.Int) *types.Transaction); ok {
		r0 = rf(opts, consensus, args, messageId, deadline)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.TransactOpts, compass.Struct2, compass.Struct3, *big.Int, *big.Int) error); ok {
		r1 = rf(opts, consensus, args, messageId, deadline)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TurnstoneId provides a mock function with given fields: opts
func (_m *CompassBinding) TurnstoneId(opts *bind.CallOpts) ([32]byte, error) {
	ret := _m.Called(opts)

	var r0 [32]byte
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) ([32]byte, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) [32]byte); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([32]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateValset provides a mock function with given fields: opts, consensus, newValset
func (_m *CompassBinding) UpdateValset(opts *bind.TransactOpts, consensus compass.Struct2, newValset compass.Struct0) (*types.Transaction, error) {
	ret := _m.Called(opts, consensus, newValset)

	var r0 *types.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, compass.Struct2, compass.Struct0) (*types.Transaction, error)); ok {
		return rf(opts, consensus, newValset)
	}
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, compass.Struct2, compass.Struct0) *types.Transaction); ok {
		r0 = rf(opts, consensus, newValset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.TransactOpts, compass.Struct2, compass.Struct0) error); ok {
		r1 = rf(opts, consensus, newValset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WatchLogicCallEvent provides a mock function with given fields: opts, sink
func (_m *CompassBinding) WatchLogicCallEvent(opts *bind.WatchOpts, sink chan<- *compass.CompassLogicCallEvent) (event.Subscription, error) {
	ret := _m.Called(opts, sink)

	var r0 event.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.WatchOpts, chan<- *compass.CompassLogicCallEvent) (event.Subscription, error)); ok {
		return rf(opts, sink)
	}
	if rf, ok := ret.Get(0).(func(*bind.WatchOpts, chan<- *compass.CompassLogicCallEvent) event.Subscription); ok {
		r0 = rf(opts, sink)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(event.Subscription)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.WatchOpts, chan<- *compass.CompassLogicCallEvent) error); ok {
		r1 = rf(opts, sink)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WatchValsetUpdated provides a mock function with given fields: opts, sink
func (_m *CompassBinding) WatchValsetUpdated(opts *bind.WatchOpts, sink chan<- *compass.CompassValsetUpdated) (event.Subscription, error) {
	ret := _m.Called(opts, sink)

	var r0 event.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.WatchOpts, chan<- *compass.CompassValsetUpdated) (event.Subscription, error)); ok {
		return rf(opts, sink)
	}
	if rf, ok := ret.Get(0).(func(*bind.WatchOpts, chan<- *compass.CompassValsetUpdated) event.Subscription); ok {
		r0 = rf(opts, sink)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(event.Subscription)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.WatchOpts, chan<- *compass.CompassValsetUpdated) error); ok {
		r1 = rf(opts, sink)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCompassBinding creates a new instance of CompassBinding. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCompassBinding(t interface {
	mock.TestingT
	Cleanup(func())
}) *CompassBinding {
	mock := &CompassBinding{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
