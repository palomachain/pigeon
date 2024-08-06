// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	evmtypes "github.com/palomachain/paloma/x/evm/types"

	mock "github.com/stretchr/testify/mock"

	paloma "github.com/palomachain/pigeon/chain/paloma"

	proto "github.com/cosmos/gogoproto/proto"

	types "github.com/palomachain/paloma/x/skyway/types"

	valsettypes "github.com/palomachain/paloma/x/valset/types"
)

// PalomaClienter is an autogenerated mock type for the PalomaClienter type
type PalomaClienter struct {
	mock.Mock
}

// AddMessageEvidence provides a mock function with given fields: ctx, queueTypeName, messageID, proof
func (_m *PalomaClienter) AddMessageEvidence(ctx context.Context, queueTypeName string, messageID uint64, proof proto.Message) error {
	ret := _m.Called(ctx, queueTypeName, messageID, proof)

	if len(ret) == 0 {
		panic("no return value specified for AddMessageEvidence")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uint64, proto.Message) error); ok {
		r0 = rf(ctx, queueTypeName, messageID, proof)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStatus provides a mock function with given fields:
func (_m *PalomaClienter) NewStatus() paloma.StatusUpdater {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NewStatus")
	}

	var r0 paloma.StatusUpdater
	if rf, ok := ret.Get(0).(func() paloma.StatusUpdater); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(paloma.StatusUpdater)
		}
	}

	return r0
}

// QueryBatchRequestByNonce provides a mock function with given fields: ctx, nonce, contract
func (_m *PalomaClienter) QueryBatchRequestByNonce(ctx context.Context, nonce uint64, contract string) (types.OutgoingTxBatch, error) {
	ret := _m.Called(ctx, nonce, contract)

	if len(ret) == 0 {
		panic("no return value specified for QueryBatchRequestByNonce")
	}

	var r0 types.OutgoingTxBatch
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, string) (types.OutgoingTxBatch, error)); ok {
		return rf(ctx, nonce, contract)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, string) types.OutgoingTxBatch); ok {
		r0 = rf(ctx, nonce, contract)
	} else {
		r0 = ret.Get(0).(types.OutgoingTxBatch)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, string) error); ok {
		r1 = rf(ctx, nonce, contract)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryGetEVMValsetByID provides a mock function with given fields: ctx, id, chainID
func (_m *PalomaClienter) QueryGetEVMValsetByID(ctx context.Context, id uint64, chainID string) (*evmtypes.Valset, error) {
	ret := _m.Called(ctx, id, chainID)

	if len(ret) == 0 {
		panic("no return value specified for QueryGetEVMValsetByID")
	}

	var r0 *evmtypes.Valset
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, string) (*evmtypes.Valset, error)); ok {
		return rf(ctx, id, chainID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, string) *evmtypes.Valset); ok {
		r0 = rf(ctx, id, chainID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*evmtypes.Valset)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, string) error); ok {
		r1 = rf(ctx, id, chainID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryGetLatestPublishedSnapshot provides a mock function with given fields: ctx, chainReferenceID
func (_m *PalomaClienter) QueryGetLatestPublishedSnapshot(ctx context.Context, chainReferenceID string) (*valsettypes.Snapshot, error) {
	ret := _m.Called(ctx, chainReferenceID)

	if len(ret) == 0 {
		panic("no return value specified for QueryGetLatestPublishedSnapshot")
	}

	var r0 *valsettypes.Snapshot
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*valsettypes.Snapshot, error)); ok {
		return rf(ctx, chainReferenceID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *valsettypes.Snapshot); ok {
		r0 = rf(ctx, chainReferenceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*valsettypes.Snapshot)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, chainReferenceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryLastObservedSkywayNonceByAddr provides a mock function with given fields: ctx, chainReferenceID, orchestrator
func (_m *PalomaClienter) QueryLastObservedSkywayNonceByAddr(ctx context.Context, chainReferenceID string, orchestrator string) (uint64, error) {
	ret := _m.Called(ctx, chainReferenceID, orchestrator)

	if len(ret) == 0 {
		panic("no return value specified for QueryLastObservedSkywayNonceByAddr")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (uint64, error)); ok {
		return rf(ctx, chainReferenceID, orchestrator)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) uint64); ok {
		r0 = rf(ctx, chainReferenceID, orchestrator)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, chainReferenceID, orchestrator)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendBatchSendToEVMClaim provides a mock function with given fields: ctx, claim
func (_m *PalomaClienter) SendBatchSendToEVMClaim(ctx context.Context, claim types.MsgBatchSendToRemoteClaim) error {
	ret := _m.Called(ctx, claim)

	if len(ret) == 0 {
		panic("no return value specified for SendBatchSendToEVMClaim")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.MsgBatchSendToRemoteClaim) error); ok {
		r0 = rf(ctx, claim)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendLightNodeSaleClaim provides a mock function with given fields: ctx, claim
func (_m *PalomaClienter) SendLightNodeSaleClaim(ctx context.Context, claim types.MsgLightNodeSaleClaim) error {
	ret := _m.Called(ctx, claim)

	if len(ret) == 0 {
		panic("no return value specified for SendLightNodeSaleClaim")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.MsgLightNodeSaleClaim) error); ok {
		r0 = rf(ctx, claim)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendSendToPalomaClaim provides a mock function with given fields: ctx, claim
func (_m *PalomaClienter) SendSendToPalomaClaim(ctx context.Context, claim types.MsgSendToPalomaClaim) error {
	ret := _m.Called(ctx, claim)

	if len(ret) == 0 {
		panic("no return value specified for SendSendToPalomaClaim")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.MsgSendToPalomaClaim) error); ok {
		r0 = rf(ctx, claim)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetErrorData provides a mock function with given fields: ctx, queueTypeName, messageID, data
func (_m *PalomaClienter) SetErrorData(ctx context.Context, queueTypeName string, messageID uint64, data []byte) error {
	ret := _m.Called(ctx, queueTypeName, messageID, data)

	if len(ret) == 0 {
		panic("no return value specified for SetErrorData")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uint64, []byte) error); ok {
		r0 = rf(ctx, queueTypeName, messageID, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetPublicAccessData provides a mock function with given fields: ctx, queueTypeName, messageID, valsetID, data
func (_m *PalomaClienter) SetPublicAccessData(ctx context.Context, queueTypeName string, messageID uint64, valsetID uint64, data []byte) error {
	ret := _m.Called(ctx, queueTypeName, messageID, valsetID, data)

	if len(ret) == 0 {
		panic("no return value specified for SetPublicAccessData")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uint64, uint64, []byte) error); ok {
		r0 = rf(ctx, queueTypeName, messageID, valsetID, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPalomaClienter creates a new instance of PalomaClienter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPalomaClienter(t interface {
	mock.TestingT
	Cleanup(func())
}) *PalomaClienter {
	mock := &PalomaClienter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
