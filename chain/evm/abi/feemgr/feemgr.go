// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package feemgr

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// Struct0 is an auto generated low-level Go binding around an user-defined struct.
type Struct0 struct {
	RelayerFee            *big.Int
	CommunityFee          *big.Int
	SecurityFee           *big.Int
	FeePayerPalomaAddress [32]byte
}

// FeemgrMetaData contains all meta data concerning the Feemgr contract.
var FeemgrMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"name\":\"grain\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"name\":\"depositor_paloma_address\",\"type\":\"bytes32\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"dex\",\"type\":\"address\"},{\"name\":\"payload\",\"type\":\"bytes\"},{\"name\":\"min_grain\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"security_fee_topup\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"name\":\"relayer_fee\",\"type\":\"uint256\"},{\"name\":\"community_fee\",\"type\":\"uint256\"},{\"name\":\"security_fee\",\"type\":\"uint256\"},{\"name\":\"fee_payer_paloma_address\",\"type\":\"bytes32\"}],\"name\":\"fee_args\",\"type\":\"tuple\"},{\"name\":\"relayer\",\"type\":\"address\"}],\"name\":\"transfer_fees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"gas_fee_amount\",\"type\":\"uint256\"}],\"name\":\"reserve_security_fee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"dex\",\"type\":\"address\"},{\"name\":\"payload\",\"type\":\"bytes\"},{\"name\":\"min_grain\",\"type\":\"uint256\"}],\"name\":\"bridge_community_fee_to_paloma\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_new_compass\",\"type\":\"address\"}],\"name\":\"update_compass\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_compass\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"compass\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"grain\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewards_community_balance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewards_security_balance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"arg0\",\"type\":\"bytes32\"}],\"name\":\"funds\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"arg0\",\"type\":\"address\"}],\"name\":\"claimable_rewards\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"total_funds\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"total_claims\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// FeemgrABI is the input ABI used to generate the binding from.
// Deprecated: Use FeemgrMetaData.ABI instead.
var FeemgrABI = FeemgrMetaData.ABI

// Feemgr is an auto generated Go binding around an Ethereum contract.
type Feemgr struct {
	FeemgrCaller     // Read-only binding to the contract
	FeemgrTransactor // Write-only binding to the contract
	FeemgrFilterer   // Log filterer for contract events
}

// FeemgrCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeemgrCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeemgrTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeemgrTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeemgrFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeemgrFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeemgrSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeemgrSession struct {
	Contract     *Feemgr           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeemgrCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeemgrCallerSession struct {
	Contract *FeemgrCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FeemgrTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeemgrTransactorSession struct {
	Contract     *FeemgrTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeemgrRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeemgrRaw struct {
	Contract *Feemgr // Generic contract binding to access the raw methods on
}

// FeemgrCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeemgrCallerRaw struct {
	Contract *FeemgrCaller // Generic read-only contract binding to access the raw methods on
}

// FeemgrTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeemgrTransactorRaw struct {
	Contract *FeemgrTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeemgr creates a new instance of Feemgr, bound to a specific deployed contract.
func NewFeemgr(address common.Address, backend bind.ContractBackend) (*Feemgr, error) {
	contract, err := bindFeemgr(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Feemgr{FeemgrCaller: FeemgrCaller{contract: contract}, FeemgrTransactor: FeemgrTransactor{contract: contract}, FeemgrFilterer: FeemgrFilterer{contract: contract}}, nil
}

// NewFeemgrCaller creates a new read-only instance of Feemgr, bound to a specific deployed contract.
func NewFeemgrCaller(address common.Address, caller bind.ContractCaller) (*FeemgrCaller, error) {
	contract, err := bindFeemgr(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeemgrCaller{contract: contract}, nil
}

// NewFeemgrTransactor creates a new write-only instance of Feemgr, bound to a specific deployed contract.
func NewFeemgrTransactor(address common.Address, transactor bind.ContractTransactor) (*FeemgrTransactor, error) {
	contract, err := bindFeemgr(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeemgrTransactor{contract: contract}, nil
}

// NewFeemgrFilterer creates a new log filterer instance of Feemgr, bound to a specific deployed contract.
func NewFeemgrFilterer(address common.Address, filterer bind.ContractFilterer) (*FeemgrFilterer, error) {
	contract, err := bindFeemgr(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeemgrFilterer{contract: contract}, nil
}

// bindFeemgr binds a generic wrapper to an already deployed contract.
func bindFeemgr(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeemgrMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Feemgr *FeemgrRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Feemgr.Contract.FeemgrCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Feemgr *FeemgrRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Feemgr.Contract.FeemgrTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Feemgr *FeemgrRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Feemgr.Contract.FeemgrTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Feemgr *FeemgrCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Feemgr.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Feemgr *FeemgrTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Feemgr.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Feemgr *FeemgrTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Feemgr.Contract.contract.Transact(opts, method, params...)
}

// ClaimableRewards is a free data retrieval call binding the contract method 0x6d84534a.
//
// Solidity: function claimable_rewards(address arg0) view returns(uint256)
func (_Feemgr *FeemgrCaller) ClaimableRewards(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Feemgr.contract.Call(opts, &out, "claimable_rewards", arg0)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// ClaimableRewards is a free data retrieval call binding the contract method 0x6d84534a.
//
// Solidity: function claimable_rewards(address arg0) view returns(uint256)
func (_Feemgr *FeemgrSession) ClaimableRewards(arg0 common.Address) (*big.Int, error) {
	return _Feemgr.Contract.ClaimableRewards(&_Feemgr.CallOpts, arg0)
}

// ClaimableRewards is a free data retrieval call binding the contract method 0x6d84534a.
//
// Solidity: function claimable_rewards(address arg0) view returns(uint256)
func (_Feemgr *FeemgrCallerSession) ClaimableRewards(arg0 common.Address) (*big.Int, error) {
	return _Feemgr.Contract.ClaimableRewards(&_Feemgr.CallOpts, arg0)
}

// Compass is a free data retrieval call binding the contract method 0xeb8acce6.
//
// Solidity: function compass() view returns(address)
func (_Feemgr *FeemgrCaller) Compass(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feemgr.contract.Call(opts, &out, "compass")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Compass is a free data retrieval call binding the contract method 0xeb8acce6.
//
// Solidity: function compass() view returns(address)
func (_Feemgr *FeemgrSession) Compass() (common.Address, error) {
	return _Feemgr.Contract.Compass(&_Feemgr.CallOpts)
}

// Compass is a free data retrieval call binding the contract method 0xeb8acce6.
//
// Solidity: function compass() view returns(address)
func (_Feemgr *FeemgrCallerSession) Compass() (common.Address, error) {
	return _Feemgr.Contract.Compass(&_Feemgr.CallOpts)
}

// Funds is a free data retrieval call binding the contract method 0x6a29a620.
//
// Solidity: function funds(bytes32 arg0) view returns(uint256)
func (_Feemgr *FeemgrCaller) Funds(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Feemgr.contract.Call(opts, &out, "funds", arg0)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Funds is a free data retrieval call binding the contract method 0x6a29a620.
//
// Solidity: function funds(bytes32 arg0) view returns(uint256)
func (_Feemgr *FeemgrSession) Funds(arg0 [32]byte) (*big.Int, error) {
	return _Feemgr.Contract.Funds(&_Feemgr.CallOpts, arg0)
}

// Funds is a free data retrieval call binding the contract method 0x6a29a620.
//
// Solidity: function funds(bytes32 arg0) view returns(uint256)
func (_Feemgr *FeemgrCallerSession) Funds(arg0 [32]byte) (*big.Int, error) {
	return _Feemgr.Contract.Funds(&_Feemgr.CallOpts, arg0)
}

// Grain is a free data retrieval call binding the contract method 0x7f67b0ff.
//
// Solidity: function grain() view returns(address)
func (_Feemgr *FeemgrCaller) Grain(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feemgr.contract.Call(opts, &out, "grain")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Grain is a free data retrieval call binding the contract method 0x7f67b0ff.
//
// Solidity: function grain() view returns(address)
func (_Feemgr *FeemgrSession) Grain() (common.Address, error) {
	return _Feemgr.Contract.Grain(&_Feemgr.CallOpts)
}

// Grain is a free data retrieval call binding the contract method 0x7f67b0ff.
//
// Solidity: function grain() view returns(address)
func (_Feemgr *FeemgrCallerSession) Grain() (common.Address, error) {
	return _Feemgr.Contract.Grain(&_Feemgr.CallOpts)
}

// RewardsCommunityBalance is a free data retrieval call binding the contract method 0x1bf23890.
//
// Solidity: function rewards_community_balance() view returns(uint256)
func (_Feemgr *FeemgrCaller) RewardsCommunityBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Feemgr.contract.Call(opts, &out, "rewards_community_balance")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// RewardsCommunityBalance is a free data retrieval call binding the contract method 0x1bf23890.
//
// Solidity: function rewards_community_balance() view returns(uint256)
func (_Feemgr *FeemgrSession) RewardsCommunityBalance() (*big.Int, error) {
	return _Feemgr.Contract.RewardsCommunityBalance(&_Feemgr.CallOpts)
}

// RewardsCommunityBalance is a free data retrieval call binding the contract method 0x1bf23890.
//
// Solidity: function rewards_community_balance() view returns(uint256)
func (_Feemgr *FeemgrCallerSession) RewardsCommunityBalance() (*big.Int, error) {
	return _Feemgr.Contract.RewardsCommunityBalance(&_Feemgr.CallOpts)
}

// RewardsSecurityBalance is a free data retrieval call binding the contract method 0x346c573e.
//
// Solidity: function rewards_security_balance() view returns(uint256)
func (_Feemgr *FeemgrCaller) RewardsSecurityBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Feemgr.contract.Call(opts, &out, "rewards_security_balance")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// RewardsSecurityBalance is a free data retrieval call binding the contract method 0x346c573e.
//
// Solidity: function rewards_security_balance() view returns(uint256)
func (_Feemgr *FeemgrSession) RewardsSecurityBalance() (*big.Int, error) {
	return _Feemgr.Contract.RewardsSecurityBalance(&_Feemgr.CallOpts)
}

// RewardsSecurityBalance is a free data retrieval call binding the contract method 0x346c573e.
//
// Solidity: function rewards_security_balance() view returns(uint256)
func (_Feemgr *FeemgrCallerSession) RewardsSecurityBalance() (*big.Int, error) {
	return _Feemgr.Contract.RewardsSecurityBalance(&_Feemgr.CallOpts)
}

// TotalClaims is a free data retrieval call binding the contract method 0xc22416b0.
//
// Solidity: function total_claims() view returns(uint256)
func (_Feemgr *FeemgrCaller) TotalClaims(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Feemgr.contract.Call(opts, &out, "total_claims")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// TotalClaims is a free data retrieval call binding the contract method 0xc22416b0.
//
// Solidity: function total_claims() view returns(uint256)
func (_Feemgr *FeemgrSession) TotalClaims() (*big.Int, error) {
	return _Feemgr.Contract.TotalClaims(&_Feemgr.CallOpts)
}

// TotalClaims is a free data retrieval call binding the contract method 0xc22416b0.
//
// Solidity: function total_claims() view returns(uint256)
func (_Feemgr *FeemgrCallerSession) TotalClaims() (*big.Int, error) {
	return _Feemgr.Contract.TotalClaims(&_Feemgr.CallOpts)
}

// TotalFunds is a free data retrieval call binding the contract method 0x34138814.
//
// Solidity: function total_funds() view returns(uint256)
func (_Feemgr *FeemgrCaller) TotalFunds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Feemgr.contract.Call(opts, &out, "total_funds")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// TotalFunds is a free data retrieval call binding the contract method 0x34138814.
//
// Solidity: function total_funds() view returns(uint256)
func (_Feemgr *FeemgrSession) TotalFunds() (*big.Int, error) {
	return _Feemgr.Contract.TotalFunds(&_Feemgr.CallOpts)
}

// TotalFunds is a free data retrieval call binding the contract method 0x34138814.
//
// Solidity: function total_funds() view returns(uint256)
func (_Feemgr *FeemgrCallerSession) TotalFunds() (*big.Int, error) {
	return _Feemgr.Contract.TotalFunds(&_Feemgr.CallOpts)
}

// BridgeCommunityFeeToPaloma is a paid mutator transaction binding the contract method 0x06c9624d.
//
// Solidity: function bridge_community_fee_to_paloma(uint256 amount, address dex, bytes payload, uint256 min_grain) returns(uint256)
func (_Feemgr *FeemgrTransactor) BridgeCommunityFeeToPaloma(opts *bind.TransactOpts, amount *big.Int, dex common.Address, payload []byte, min_grain *big.Int) (*types.Transaction, error) {
	return _Feemgr.contract.Transact(opts, "bridge_community_fee_to_paloma", amount, dex, payload, min_grain)
}

// BridgeCommunityFeeToPaloma is a paid mutator transaction binding the contract method 0x06c9624d.
//
// Solidity: function bridge_community_fee_to_paloma(uint256 amount, address dex, bytes payload, uint256 min_grain) returns(uint256)
func (_Feemgr *FeemgrSession) BridgeCommunityFeeToPaloma(amount *big.Int, dex common.Address, payload []byte, min_grain *big.Int) (*types.Transaction, error) {
	return _Feemgr.Contract.BridgeCommunityFeeToPaloma(&_Feemgr.TransactOpts, amount, dex, payload, min_grain)
}

// BridgeCommunityFeeToPaloma is a paid mutator transaction binding the contract method 0x06c9624d.
//
// Solidity: function bridge_community_fee_to_paloma(uint256 amount, address dex, bytes payload, uint256 min_grain) returns(uint256)
func (_Feemgr *FeemgrTransactorSession) BridgeCommunityFeeToPaloma(amount *big.Int, dex common.Address, payload []byte, min_grain *big.Int) (*types.Transaction, error) {
	return _Feemgr.Contract.BridgeCommunityFeeToPaloma(&_Feemgr.TransactOpts, amount, dex, payload, min_grain)
}

// Deposit is a paid mutator transaction binding the contract method 0xb214faa5.
//
// Solidity: function deposit(bytes32 depositor_paloma_address) payable returns()
func (_Feemgr *FeemgrTransactor) Deposit(opts *bind.TransactOpts, depositor_paloma_address [32]byte) (*types.Transaction, error) {
	return _Feemgr.contract.Transact(opts, "deposit", depositor_paloma_address)
}

// Deposit is a paid mutator transaction binding the contract method 0xb214faa5.
//
// Solidity: function deposit(bytes32 depositor_paloma_address) payable returns()
func (_Feemgr *FeemgrSession) Deposit(depositor_paloma_address [32]byte) (*types.Transaction, error) {
	return _Feemgr.Contract.Deposit(&_Feemgr.TransactOpts, depositor_paloma_address)
}

// Deposit is a paid mutator transaction binding the contract method 0xb214faa5.
//
// Solidity: function deposit(bytes32 depositor_paloma_address) payable returns()
func (_Feemgr *FeemgrTransactorSession) Deposit(depositor_paloma_address [32]byte) (*types.Transaction, error) {
	return _Feemgr.Contract.Deposit(&_Feemgr.TransactOpts, depositor_paloma_address)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _compass) returns()
func (_Feemgr *FeemgrTransactor) Initialize(opts *bind.TransactOpts, _compass common.Address) (*types.Transaction, error) {
	return _Feemgr.contract.Transact(opts, "initialize", _compass)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _compass) returns()
func (_Feemgr *FeemgrSession) Initialize(_compass common.Address) (*types.Transaction, error) {
	return _Feemgr.Contract.Initialize(&_Feemgr.TransactOpts, _compass)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _compass) returns()
func (_Feemgr *FeemgrTransactorSession) Initialize(_compass common.Address) (*types.Transaction, error) {
	return _Feemgr.Contract.Initialize(&_Feemgr.TransactOpts, _compass)
}

// ReserveSecurityFee is a paid mutator transaction binding the contract method 0xd4bbab4f.
//
// Solidity: function reserve_security_fee(address sender, uint256 gas_fee_amount) returns()
func (_Feemgr *FeemgrTransactor) ReserveSecurityFee(opts *bind.TransactOpts, sender common.Address, gas_fee_amount *big.Int) (*types.Transaction, error) {
	return _Feemgr.contract.Transact(opts, "reserve_security_fee", sender, gas_fee_amount)
}

// ReserveSecurityFee is a paid mutator transaction binding the contract method 0xd4bbab4f.
//
// Solidity: function reserve_security_fee(address sender, uint256 gas_fee_amount) returns()
func (_Feemgr *FeemgrSession) ReserveSecurityFee(sender common.Address, gas_fee_amount *big.Int) (*types.Transaction, error) {
	return _Feemgr.Contract.ReserveSecurityFee(&_Feemgr.TransactOpts, sender, gas_fee_amount)
}

// ReserveSecurityFee is a paid mutator transaction binding the contract method 0xd4bbab4f.
//
// Solidity: function reserve_security_fee(address sender, uint256 gas_fee_amount) returns()
func (_Feemgr *FeemgrTransactorSession) ReserveSecurityFee(sender common.Address, gas_fee_amount *big.Int) (*types.Transaction, error) {
	return _Feemgr.Contract.ReserveSecurityFee(&_Feemgr.TransactOpts, sender, gas_fee_amount)
}

// SecurityFeeTopup is a paid mutator transaction binding the contract method 0x1d046325.
//
// Solidity: function security_fee_topup() payable returns()
func (_Feemgr *FeemgrTransactor) SecurityFeeTopup(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Feemgr.contract.Transact(opts, "security_fee_topup")
}

// SecurityFeeTopup is a paid mutator transaction binding the contract method 0x1d046325.
//
// Solidity: function security_fee_topup() payable returns()
func (_Feemgr *FeemgrSession) SecurityFeeTopup() (*types.Transaction, error) {
	return _Feemgr.Contract.SecurityFeeTopup(&_Feemgr.TransactOpts)
}

// SecurityFeeTopup is a paid mutator transaction binding the contract method 0x1d046325.
//
// Solidity: function security_fee_topup() payable returns()
func (_Feemgr *FeemgrTransactorSession) SecurityFeeTopup() (*types.Transaction, error) {
	return _Feemgr.Contract.SecurityFeeTopup(&_Feemgr.TransactOpts)
}

// TransferFees is a paid mutator transaction binding the contract method 0xffd8d4aa.
//
// Solidity: function transfer_fees((uint256,uint256,uint256,bytes32) fee_args, address relayer) returns()
func (_Feemgr *FeemgrTransactor) TransferFees(opts *bind.TransactOpts, fee_args Struct0, relayer common.Address) (*types.Transaction, error) {
	return _Feemgr.contract.Transact(opts, "transfer_fees", fee_args, relayer)
}

// TransferFees is a paid mutator transaction binding the contract method 0xffd8d4aa.
//
// Solidity: function transfer_fees((uint256,uint256,uint256,bytes32) fee_args, address relayer) returns()
func (_Feemgr *FeemgrSession) TransferFees(fee_args Struct0, relayer common.Address) (*types.Transaction, error) {
	return _Feemgr.Contract.TransferFees(&_Feemgr.TransactOpts, fee_args, relayer)
}

// TransferFees is a paid mutator transaction binding the contract method 0xffd8d4aa.
//
// Solidity: function transfer_fees((uint256,uint256,uint256,bytes32) fee_args, address relayer) returns()
func (_Feemgr *FeemgrTransactorSession) TransferFees(fee_args Struct0, relayer common.Address) (*types.Transaction, error) {
	return _Feemgr.Contract.TransferFees(&_Feemgr.TransactOpts, fee_args, relayer)
}

// UpdateCompass is a paid mutator transaction binding the contract method 0x6974af69.
//
// Solidity: function update_compass(address _new_compass) returns()
func (_Feemgr *FeemgrTransactor) UpdateCompass(opts *bind.TransactOpts, _new_compass common.Address) (*types.Transaction, error) {
	return _Feemgr.contract.Transact(opts, "update_compass", _new_compass)
}

// UpdateCompass is a paid mutator transaction binding the contract method 0x6974af69.
//
// Solidity: function update_compass(address _new_compass) returns()
func (_Feemgr *FeemgrSession) UpdateCompass(_new_compass common.Address) (*types.Transaction, error) {
	return _Feemgr.Contract.UpdateCompass(&_Feemgr.TransactOpts, _new_compass)
}

// UpdateCompass is a paid mutator transaction binding the contract method 0x6974af69.
//
// Solidity: function update_compass(address _new_compass) returns()
func (_Feemgr *FeemgrTransactorSession) UpdateCompass(_new_compass common.Address) (*types.Transaction, error) {
	return _Feemgr.Contract.UpdateCompass(&_Feemgr.TransactOpts, _new_compass)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd07e9fa0.
//
// Solidity: function withdraw(address receiver, uint256 amount, address dex, bytes payload, uint256 min_grain) returns()
func (_Feemgr *FeemgrTransactor) Withdraw(opts *bind.TransactOpts, receiver common.Address, amount *big.Int, dex common.Address, payload []byte, min_grain *big.Int) (*types.Transaction, error) {
	return _Feemgr.contract.Transact(opts, "withdraw", receiver, amount, dex, payload, min_grain)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd07e9fa0.
//
// Solidity: function withdraw(address receiver, uint256 amount, address dex, bytes payload, uint256 min_grain) returns()
func (_Feemgr *FeemgrSession) Withdraw(receiver common.Address, amount *big.Int, dex common.Address, payload []byte, min_grain *big.Int) (*types.Transaction, error) {
	return _Feemgr.Contract.Withdraw(&_Feemgr.TransactOpts, receiver, amount, dex, payload, min_grain)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd07e9fa0.
//
// Solidity: function withdraw(address receiver, uint256 amount, address dex, bytes payload, uint256 min_grain) returns()
func (_Feemgr *FeemgrTransactorSession) Withdraw(receiver common.Address, amount *big.Int, dex common.Address, payload []byte, min_grain *big.Int) (*types.Transaction, error) {
	return _Feemgr.Contract.Withdraw(&_Feemgr.TransactOpts, receiver, amount, dex, payload, min_grain)
}
