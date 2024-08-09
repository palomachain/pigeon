// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package compass

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

// Struct2 is an auto generated low-level Go binding around an user-defined struct.
type Struct2 struct {
	Valset     Struct0
	Signatures []Struct1
}

// Struct3 is an auto generated low-level Go binding around an user-defined struct.
type Struct3 struct {
	LogicContractAddress common.Address
	Payload              []byte
}

// Struct5 is an auto generated low-level Go binding around an user-defined struct.
type Struct5 struct {
	Receiver []common.Address
	Amount   []*big.Int
}

// Struct0 is an auto generated low-level Go binding around an user-defined struct.
type Struct0 struct {
	Validators []common.Address
	Powers     []*big.Int
	ValsetId   *big.Int
}

// Struct1 is an auto generated low-level Go binding around an user-defined struct.
type Struct1 struct {
	V *big.Int
	R *big.Int
	S *big.Int
}

// Struct4 is an auto generated low-level Go binding around an user-defined struct.
type Struct4 struct {
	RelayerFee            *big.Int
	CommunityFee          *big.Int
	SecurityFee           *big.Int
	FeePayerPalomaAddress [32]byte
}

// CompassMetaData contains all meta data concerning the Compass contract.
var CompassMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"checkpoint\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"valset_id\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"ValsetUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"logic_contract_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"message_id\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"LogicCallEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"receiver\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"SendToPalomaEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"batch_id\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"BatchSendEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"paloma_denom\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"token_contract\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"symbol\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"decimals\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"ERC20DeployedEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"depositor_paloma_address\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsDepositedEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsWithdrawnEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"new_compass\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"UpdateCompassAddressInFeeManager\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"contract_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"paloma\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"node_count\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"grain_amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"NodeSaleEvent\",\"type\":\"event\"},{\"inputs\":[{\"name\":\"_compass_id\",\"type\":\"bytes32\"},{\"name\":\"_event_id\",\"type\":\"uint256\"},{\"name\":\"_gravity_nonce\",\"type\":\"uint256\"},{\"components\":[{\"name\":\"validators\",\"type\":\"address[]\"},{\"name\":\"powers\",\"type\":\"uint256[]\"},{\"name\":\"valset_id\",\"type\":\"uint256\"}],\"name\":\"valset\",\"type\":\"tuple\"},{\"name\":\"fee_manager\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"components\":[{\"components\":[{\"name\":\"validators\",\"type\":\"address[]\"},{\"name\":\"powers\",\"type\":\"uint256[]\"},{\"name\":\"valset_id\",\"type\":\"uint256\"}],\"name\":\"valset\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"v\",\"type\":\"uint256\"},{\"name\":\"r\",\"type\":\"uint256\"},{\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"signatures\",\"type\":\"tuple[]\"}],\"name\":\"consensus\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"validators\",\"type\":\"address[]\"},{\"name\":\"powers\",\"type\":\"uint256[]\"},{\"name\":\"valset_id\",\"type\":\"uint256\"}],\"name\":\"new_valset\",\"type\":\"tuple\"},{\"name\":\"relayer\",\"type\":\"address\"},{\"name\":\"gas_estimate\",\"type\":\"uint256\"}],\"name\":\"update_valset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"name\":\"validators\",\"type\":\"address[]\"},{\"name\":\"powers\",\"type\":\"uint256[]\"},{\"name\":\"valset_id\",\"type\":\"uint256\"}],\"name\":\"valset\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"v\",\"type\":\"uint256\"},{\"name\":\"r\",\"type\":\"uint256\"},{\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"signatures\",\"type\":\"tuple[]\"}],\"name\":\"consensus\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"logic_contract_address\",\"type\":\"address\"},{\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"args\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"relayer_fee\",\"type\":\"uint256\"},{\"name\":\"community_fee\",\"type\":\"uint256\"},{\"name\":\"security_fee\",\"type\":\"uint256\"},{\"name\":\"fee_payer_paloma_address\",\"type\":\"bytes32\"}],\"name\":\"fee_args\",\"type\":\"tuple\"},{\"name\":\"message_id\",\"type\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\"},{\"name\":\"relayer\",\"type\":\"address\"}],\"name\":\"submit_logic_call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"receiver\",\"type\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"send_token_to_paloma\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"name\":\"validators\",\"type\":\"address[]\"},{\"name\":\"powers\",\"type\":\"uint256[]\"},{\"name\":\"valset_id\",\"type\":\"uint256\"}],\"name\":\"valset\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"v\",\"type\":\"uint256\"},{\"name\":\"r\",\"type\":\"uint256\"},{\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"signatures\",\"type\":\"tuple[]\"}],\"name\":\"consensus\",\"type\":\"tuple\"},{\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"name\":\"receiver\",\"type\":\"address[]\"},{\"name\":\"amount\",\"type\":\"uint256[]\"}],\"name\":\"args\",\"type\":\"tuple\"},{\"name\":\"batch_id\",\"type\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\"},{\"name\":\"relayer\",\"type\":\"address\"},{\"name\":\"gas_estimate\",\"type\":\"uint256\"}],\"name\":\"submit_batch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"buyer\",\"type\":\"address\"},{\"name\":\"paloma\",\"type\":\"bytes32\"},{\"name\":\"node_count\",\"type\":\"uint256\"},{\"name\":\"grain_amount\",\"type\":\"uint256\"}],\"name\":\"emit_nodesale_event\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_paloma_denom\",\"type\":\"string\"},{\"name\":\"_name\",\"type\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\"},{\"name\":\"_decimals\",\"type\":\"uint8\"},{\"name\":\"_blueprint\",\"type\":\"address\"}],\"name\":\"deploy_erc20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"contract_address\",\"type\":\"address\"},{\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"arbitrary_view\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"depositor_paloma_address\",\"type\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"dex\",\"type\":\"address\"},{\"name\":\"payload\",\"type\":\"bytes\"},{\"name\":\"min_grain\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"security_fee_topup\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"name\":\"validators\",\"type\":\"address[]\"},{\"name\":\"powers\",\"type\":\"uint256[]\"},{\"name\":\"valset_id\",\"type\":\"uint256\"}],\"name\":\"valset\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"v\",\"type\":\"uint256\"},{\"name\":\"r\",\"type\":\"uint256\"},{\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"signatures\",\"type\":\"tuple[]\"}],\"name\":\"consensus\",\"type\":\"tuple\"},{\"name\":\"message_id\",\"type\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes32\"},{\"name\":\"relayer\",\"type\":\"address\"},{\"name\":\"gas_estimate\",\"type\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"dex\",\"type\":\"address\"},{\"name\":\"payload\",\"type\":\"bytes\"},{\"name\":\"min_grain\",\"type\":\"uint256\"}],\"name\":\"bridge_community_tax_to_paloma\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"name\":\"validators\",\"type\":\"address[]\"},{\"name\":\"powers\",\"type\":\"uint256[]\"},{\"name\":\"valset_id\",\"type\":\"uint256\"}],\"name\":\"valset\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"v\",\"type\":\"uint256\"},{\"name\":\"r\",\"type\":\"uint256\"},{\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"signatures\",\"type\":\"tuple[]\"}],\"name\":\"consensus\",\"type\":\"tuple\"},{\"name\":\"deadline\",\"type\":\"uint256\"},{\"name\":\"gas_estimate\",\"type\":\"uint256\"},{\"name\":\"_new_compass\",\"type\":\"address\"},{\"name\":\"relayer\",\"type\":\"address\"}],\"name\":\"update_compass_address_in_fee_manager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"compass_id\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"last_checkpoint\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"last_valset_id\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"last_event_id\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"last_gravity_nonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"arg0\",\"type\":\"address\"}],\"name\":\"last_batch_id\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"arg0\",\"type\":\"uint256\"}],\"name\":\"message_id_used\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"slc_switch\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_MANAGER\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// CompassABI is the input ABI used to generate the binding from.
// Deprecated: Use CompassMetaData.ABI instead.
var CompassABI = CompassMetaData.ABI

// Compass is an auto generated Go binding around an Ethereum contract.
type Compass struct {
	CompassCaller     // Read-only binding to the contract
	CompassTransactor // Write-only binding to the contract
	CompassFilterer   // Log filterer for contract events
}

// CompassCaller is an auto generated read-only Go binding around an Ethereum contract.
type CompassCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompassTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CompassTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompassFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CompassFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompassSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CompassSession struct {
	Contract     *Compass          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CompassCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CompassCallerSession struct {
	Contract *CompassCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// CompassTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CompassTransactorSession struct {
	Contract     *CompassTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// CompassRaw is an auto generated low-level Go binding around an Ethereum contract.
type CompassRaw struct {
	Contract *Compass // Generic contract binding to access the raw methods on
}

// CompassCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CompassCallerRaw struct {
	Contract *CompassCaller // Generic read-only contract binding to access the raw methods on
}

// CompassTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CompassTransactorRaw struct {
	Contract *CompassTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCompass creates a new instance of Compass, bound to a specific deployed contract.
func NewCompass(address common.Address, backend bind.ContractBackend) (*Compass, error) {
	contract, err := bindCompass(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Compass{CompassCaller: CompassCaller{contract: contract}, CompassTransactor: CompassTransactor{contract: contract}, CompassFilterer: CompassFilterer{contract: contract}}, nil
}

// NewCompassCaller creates a new read-only instance of Compass, bound to a specific deployed contract.
func NewCompassCaller(address common.Address, caller bind.ContractCaller) (*CompassCaller, error) {
	contract, err := bindCompass(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CompassCaller{contract: contract}, nil
}

// NewCompassTransactor creates a new write-only instance of Compass, bound to a specific deployed contract.
func NewCompassTransactor(address common.Address, transactor bind.ContractTransactor) (*CompassTransactor, error) {
	contract, err := bindCompass(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CompassTransactor{contract: contract}, nil
}

// NewCompassFilterer creates a new log filterer instance of Compass, bound to a specific deployed contract.
func NewCompassFilterer(address common.Address, filterer bind.ContractFilterer) (*CompassFilterer, error) {
	contract, err := bindCompass(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CompassFilterer{contract: contract}, nil
}

// bindCompass binds a generic wrapper to an already deployed contract.
func bindCompass(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CompassMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Compass *CompassRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Compass.Contract.CompassCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Compass *CompassRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Compass.Contract.CompassTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Compass *CompassRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Compass.Contract.CompassTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Compass *CompassCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Compass.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Compass *CompassTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Compass.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Compass *CompassTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Compass.Contract.contract.Transact(opts, method, params...)
}

// FEEMANAGER is a free data retrieval call binding the contract method 0xea26266c.
//
// Solidity: function FEE_MANAGER() view returns(address)
func (_Compass *CompassCaller) FEEMANAGER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Compass.contract.Call(opts, &out, "FEE_MANAGER")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// FEEMANAGER is a free data retrieval call binding the contract method 0xea26266c.
//
// Solidity: function FEE_MANAGER() view returns(address)
func (_Compass *CompassSession) FEEMANAGER() (common.Address, error) {
	return _Compass.Contract.FEEMANAGER(&_Compass.CallOpts)
}

// FEEMANAGER is a free data retrieval call binding the contract method 0xea26266c.
//
// Solidity: function FEE_MANAGER() view returns(address)
func (_Compass *CompassCallerSession) FEEMANAGER() (common.Address, error) {
	return _Compass.Contract.FEEMANAGER(&_Compass.CallOpts)
}

// ArbitraryView is a free data retrieval call binding the contract method 0x0b14c545.
//
// Solidity: function arbitrary_view(address contract_address, bytes payload) view returns(bytes)
func (_Compass *CompassCaller) ArbitraryView(opts *bind.CallOpts, contract_address common.Address, payload []byte) ([]byte, error) {
	var out []interface{}
	err := _Compass.contract.Call(opts, &out, "arbitrary_view", contract_address, payload)
	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err
}

// ArbitraryView is a free data retrieval call binding the contract method 0x0b14c545.
//
// Solidity: function arbitrary_view(address contract_address, bytes payload) view returns(bytes)
func (_Compass *CompassSession) ArbitraryView(contract_address common.Address, payload []byte) ([]byte, error) {
	return _Compass.Contract.ArbitraryView(&_Compass.CallOpts, contract_address, payload)
}

// ArbitraryView is a free data retrieval call binding the contract method 0x0b14c545.
//
// Solidity: function arbitrary_view(address contract_address, bytes payload) view returns(bytes)
func (_Compass *CompassCallerSession) ArbitraryView(contract_address common.Address, payload []byte) ([]byte, error) {
	return _Compass.Contract.ArbitraryView(&_Compass.CallOpts, contract_address, payload)
}

// CompassId is a free data retrieval call binding the contract method 0x84d38f2c.
//
// Solidity: function compass_id() view returns(bytes32)
func (_Compass *CompassCaller) CompassId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Compass.contract.Call(opts, &out, "compass_id")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// CompassId is a free data retrieval call binding the contract method 0x84d38f2c.
//
// Solidity: function compass_id() view returns(bytes32)
func (_Compass *CompassSession) CompassId() ([32]byte, error) {
	return _Compass.Contract.CompassId(&_Compass.CallOpts)
}

// CompassId is a free data retrieval call binding the contract method 0x84d38f2c.
//
// Solidity: function compass_id() view returns(bytes32)
func (_Compass *CompassCallerSession) CompassId() ([32]byte, error) {
	return _Compass.Contract.CompassId(&_Compass.CallOpts)
}

// LastBatchId is a free data retrieval call binding the contract method 0xfa822567.
//
// Solidity: function last_batch_id(address arg0) view returns(uint256)
func (_Compass *CompassCaller) LastBatchId(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Compass.contract.Call(opts, &out, "last_batch_id", arg0)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// LastBatchId is a free data retrieval call binding the contract method 0xfa822567.
//
// Solidity: function last_batch_id(address arg0) view returns(uint256)
func (_Compass *CompassSession) LastBatchId(arg0 common.Address) (*big.Int, error) {
	return _Compass.Contract.LastBatchId(&_Compass.CallOpts, arg0)
}

// LastBatchId is a free data retrieval call binding the contract method 0xfa822567.
//
// Solidity: function last_batch_id(address arg0) view returns(uint256)
func (_Compass *CompassCallerSession) LastBatchId(arg0 common.Address) (*big.Int, error) {
	return _Compass.Contract.LastBatchId(&_Compass.CallOpts, arg0)
}

// LastCheckpoint is a free data retrieval call binding the contract method 0xa9a4a983.
//
// Solidity: function last_checkpoint() view returns(bytes32)
func (_Compass *CompassCaller) LastCheckpoint(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Compass.contract.Call(opts, &out, "last_checkpoint")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// LastCheckpoint is a free data retrieval call binding the contract method 0xa9a4a983.
//
// Solidity: function last_checkpoint() view returns(bytes32)
func (_Compass *CompassSession) LastCheckpoint() ([32]byte, error) {
	return _Compass.Contract.LastCheckpoint(&_Compass.CallOpts)
}

// LastCheckpoint is a free data retrieval call binding the contract method 0xa9a4a983.
//
// Solidity: function last_checkpoint() view returns(bytes32)
func (_Compass *CompassCallerSession) LastCheckpoint() ([32]byte, error) {
	return _Compass.Contract.LastCheckpoint(&_Compass.CallOpts)
}

// LastEventId is a free data retrieval call binding the contract method 0x19429b8d.
//
// Solidity: function last_event_id() view returns(uint256)
func (_Compass *CompassCaller) LastEventId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Compass.contract.Call(opts, &out, "last_event_id")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// LastEventId is a free data retrieval call binding the contract method 0x19429b8d.
//
// Solidity: function last_event_id() view returns(uint256)
func (_Compass *CompassSession) LastEventId() (*big.Int, error) {
	return _Compass.Contract.LastEventId(&_Compass.CallOpts)
}

// LastEventId is a free data retrieval call binding the contract method 0x19429b8d.
//
// Solidity: function last_event_id() view returns(uint256)
func (_Compass *CompassCallerSession) LastEventId() (*big.Int, error) {
	return _Compass.Contract.LastEventId(&_Compass.CallOpts)
}

// LastGravityNonce is a free data retrieval call binding the contract method 0x0cb39e96.
//
// Solidity: function last_gravity_nonce() view returns(uint256)
func (_Compass *CompassCaller) LastGravityNonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Compass.contract.Call(opts, &out, "last_gravity_nonce")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// LastGravityNonce is a free data retrieval call binding the contract method 0x0cb39e96.
//
// Solidity: function last_gravity_nonce() view returns(uint256)
func (_Compass *CompassSession) LastGravityNonce() (*big.Int, error) {
	return _Compass.Contract.LastGravityNonce(&_Compass.CallOpts)
}

// LastGravityNonce is a free data retrieval call binding the contract method 0x0cb39e96.
//
// Solidity: function last_gravity_nonce() view returns(uint256)
func (_Compass *CompassCallerSession) LastGravityNonce() (*big.Int, error) {
	return _Compass.Contract.LastGravityNonce(&_Compass.CallOpts)
}

// LastValsetId is a free data retrieval call binding the contract method 0x4da6ecc9.
//
// Solidity: function last_valset_id() view returns(uint256)
func (_Compass *CompassCaller) LastValsetId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Compass.contract.Call(opts, &out, "last_valset_id")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// LastValsetId is a free data retrieval call binding the contract method 0x4da6ecc9.
//
// Solidity: function last_valset_id() view returns(uint256)
func (_Compass *CompassSession) LastValsetId() (*big.Int, error) {
	return _Compass.Contract.LastValsetId(&_Compass.CallOpts)
}

// LastValsetId is a free data retrieval call binding the contract method 0x4da6ecc9.
//
// Solidity: function last_valset_id() view returns(uint256)
func (_Compass *CompassCallerSession) LastValsetId() (*big.Int, error) {
	return _Compass.Contract.LastValsetId(&_Compass.CallOpts)
}

// MessageIdUsed is a free data retrieval call binding the contract method 0x38d6172d.
//
// Solidity: function message_id_used(uint256 arg0) view returns(bool)
func (_Compass *CompassCaller) MessageIdUsed(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var out []interface{}
	err := _Compass.contract.Call(opts, &out, "message_id_used", arg0)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// MessageIdUsed is a free data retrieval call binding the contract method 0x38d6172d.
//
// Solidity: function message_id_used(uint256 arg0) view returns(bool)
func (_Compass *CompassSession) MessageIdUsed(arg0 *big.Int) (bool, error) {
	return _Compass.Contract.MessageIdUsed(&_Compass.CallOpts, arg0)
}

// MessageIdUsed is a free data retrieval call binding the contract method 0x38d6172d.
//
// Solidity: function message_id_used(uint256 arg0) view returns(bool)
func (_Compass *CompassCallerSession) MessageIdUsed(arg0 *big.Int) (bool, error) {
	return _Compass.Contract.MessageIdUsed(&_Compass.CallOpts, arg0)
}

// SlcSwitch is a free data retrieval call binding the contract method 0x844105e1.
//
// Solidity: function slc_switch() view returns(bool)
func (_Compass *CompassCaller) SlcSwitch(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Compass.contract.Call(opts, &out, "slc_switch")
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// SlcSwitch is a free data retrieval call binding the contract method 0x844105e1.
//
// Solidity: function slc_switch() view returns(bool)
func (_Compass *CompassSession) SlcSwitch() (bool, error) {
	return _Compass.Contract.SlcSwitch(&_Compass.CallOpts)
}

// SlcSwitch is a free data retrieval call binding the contract method 0x844105e1.
//
// Solidity: function slc_switch() view returns(bool)
func (_Compass *CompassCallerSession) SlcSwitch() (bool, error) {
	return _Compass.Contract.SlcSwitch(&_Compass.CallOpts)
}

// BridgeCommunityTaxToPaloma is a paid mutator transaction binding the contract method 0xa0dd2e49.
//
// Solidity: function bridge_community_tax_to_paloma(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, uint256 message_id, uint256 deadline, bytes32 receiver, address relayer, uint256 gas_estimate, uint256 amount, address dex, bytes payload, uint256 min_grain) returns()
func (_Compass *CompassTransactor) BridgeCommunityTaxToPaloma(opts *bind.TransactOpts, consensus Struct2, message_id *big.Int, deadline *big.Int, receiver [32]byte, relayer common.Address, gas_estimate *big.Int, amount *big.Int, dex common.Address, payload []byte, min_grain *big.Int) (*types.Transaction, error) {
	return _Compass.contract.Transact(opts, "bridge_community_tax_to_paloma", consensus, message_id, deadline, receiver, relayer, gas_estimate, amount, dex, payload, min_grain)
}

// BridgeCommunityTaxToPaloma is a paid mutator transaction binding the contract method 0xa0dd2e49.
//
// Solidity: function bridge_community_tax_to_paloma(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, uint256 message_id, uint256 deadline, bytes32 receiver, address relayer, uint256 gas_estimate, uint256 amount, address dex, bytes payload, uint256 min_grain) returns()
func (_Compass *CompassSession) BridgeCommunityTaxToPaloma(consensus Struct2, message_id *big.Int, deadline *big.Int, receiver [32]byte, relayer common.Address, gas_estimate *big.Int, amount *big.Int, dex common.Address, payload []byte, min_grain *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.BridgeCommunityTaxToPaloma(&_Compass.TransactOpts, consensus, message_id, deadline, receiver, relayer, gas_estimate, amount, dex, payload, min_grain)
}

// BridgeCommunityTaxToPaloma is a paid mutator transaction binding the contract method 0xa0dd2e49.
//
// Solidity: function bridge_community_tax_to_paloma(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, uint256 message_id, uint256 deadline, bytes32 receiver, address relayer, uint256 gas_estimate, uint256 amount, address dex, bytes payload, uint256 min_grain) returns()
func (_Compass *CompassTransactorSession) BridgeCommunityTaxToPaloma(consensus Struct2, message_id *big.Int, deadline *big.Int, receiver [32]byte, relayer common.Address, gas_estimate *big.Int, amount *big.Int, dex common.Address, payload []byte, min_grain *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.BridgeCommunityTaxToPaloma(&_Compass.TransactOpts, consensus, message_id, deadline, receiver, relayer, gas_estimate, amount, dex, payload, min_grain)
}

// DeployErc20 is a paid mutator transaction binding the contract method 0x08a92ad7.
//
// Solidity: function deploy_erc20(string _paloma_denom, string _name, string _symbol, uint8 _decimals, address _blueprint) returns()
func (_Compass *CompassTransactor) DeployErc20(opts *bind.TransactOpts, _paloma_denom string, _name string, _symbol string, _decimals uint8, _blueprint common.Address) (*types.Transaction, error) {
	return _Compass.contract.Transact(opts, "deploy_erc20", _paloma_denom, _name, _symbol, _decimals, _blueprint)
}

// DeployErc20 is a paid mutator transaction binding the contract method 0x08a92ad7.
//
// Solidity: function deploy_erc20(string _paloma_denom, string _name, string _symbol, uint8 _decimals, address _blueprint) returns()
func (_Compass *CompassSession) DeployErc20(_paloma_denom string, _name string, _symbol string, _decimals uint8, _blueprint common.Address) (*types.Transaction, error) {
	return _Compass.Contract.DeployErc20(&_Compass.TransactOpts, _paloma_denom, _name, _symbol, _decimals, _blueprint)
}

// DeployErc20 is a paid mutator transaction binding the contract method 0x08a92ad7.
//
// Solidity: function deploy_erc20(string _paloma_denom, string _name, string _symbol, uint8 _decimals, address _blueprint) returns()
func (_Compass *CompassTransactorSession) DeployErc20(_paloma_denom string, _name string, _symbol string, _decimals uint8, _blueprint common.Address) (*types.Transaction, error) {
	return _Compass.Contract.DeployErc20(&_Compass.TransactOpts, _paloma_denom, _name, _symbol, _decimals, _blueprint)
}

// Deposit is a paid mutator transaction binding the contract method 0x1de26e16.
//
// Solidity: function deposit(bytes32 depositor_paloma_address, uint256 amount) payable returns()
func (_Compass *CompassTransactor) Deposit(opts *bind.TransactOpts, depositor_paloma_address [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Compass.contract.Transact(opts, "deposit", depositor_paloma_address, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x1de26e16.
//
// Solidity: function deposit(bytes32 depositor_paloma_address, uint256 amount) payable returns()
func (_Compass *CompassSession) Deposit(depositor_paloma_address [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.Deposit(&_Compass.TransactOpts, depositor_paloma_address, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x1de26e16.
//
// Solidity: function deposit(bytes32 depositor_paloma_address, uint256 amount) payable returns()
func (_Compass *CompassTransactorSession) Deposit(depositor_paloma_address [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.Deposit(&_Compass.TransactOpts, depositor_paloma_address, amount)
}

// EmitNodesaleEvent is a paid mutator transaction binding the contract method 0xa3299771.
//
// Solidity: function emit_nodesale_event(address buyer, bytes32 paloma, uint256 node_count, uint256 grain_amount) returns()
func (_Compass *CompassTransactor) EmitNodesaleEvent(opts *bind.TransactOpts, buyer common.Address, paloma [32]byte, node_count *big.Int, grain_amount *big.Int) (*types.Transaction, error) {
	return _Compass.contract.Transact(opts, "emit_nodesale_event", buyer, paloma, node_count, grain_amount)
}

// EmitNodesaleEvent is a paid mutator transaction binding the contract method 0xa3299771.
//
// Solidity: function emit_nodesale_event(address buyer, bytes32 paloma, uint256 node_count, uint256 grain_amount) returns()
func (_Compass *CompassSession) EmitNodesaleEvent(buyer common.Address, paloma [32]byte, node_count *big.Int, grain_amount *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.EmitNodesaleEvent(&_Compass.TransactOpts, buyer, paloma, node_count, grain_amount)
}

// EmitNodesaleEvent is a paid mutator transaction binding the contract method 0xa3299771.
//
// Solidity: function emit_nodesale_event(address buyer, bytes32 paloma, uint256 node_count, uint256 grain_amount) returns()
func (_Compass *CompassTransactorSession) EmitNodesaleEvent(buyer common.Address, paloma [32]byte, node_count *big.Int, grain_amount *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.EmitNodesaleEvent(&_Compass.TransactOpts, buyer, paloma, node_count, grain_amount)
}

// SecurityFeeTopup is a paid mutator transaction binding the contract method 0xda29a8c6.
//
// Solidity: function security_fee_topup(uint256 amount) payable returns()
func (_Compass *CompassTransactor) SecurityFeeTopup(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Compass.contract.Transact(opts, "security_fee_topup", amount)
}

// SecurityFeeTopup is a paid mutator transaction binding the contract method 0xda29a8c6.
//
// Solidity: function security_fee_topup(uint256 amount) payable returns()
func (_Compass *CompassSession) SecurityFeeTopup(amount *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.SecurityFeeTopup(&_Compass.TransactOpts, amount)
}

// SecurityFeeTopup is a paid mutator transaction binding the contract method 0xda29a8c6.
//
// Solidity: function security_fee_topup(uint256 amount) payable returns()
func (_Compass *CompassTransactorSession) SecurityFeeTopup(amount *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.SecurityFeeTopup(&_Compass.TransactOpts, amount)
}

// SendTokenToPaloma is a paid mutator transaction binding the contract method 0xf650f6e9.
//
// Solidity: function send_token_to_paloma(address token, bytes32 receiver, uint256 amount) returns()
func (_Compass *CompassTransactor) SendTokenToPaloma(opts *bind.TransactOpts, token common.Address, receiver [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Compass.contract.Transact(opts, "send_token_to_paloma", token, receiver, amount)
}

// SendTokenToPaloma is a paid mutator transaction binding the contract method 0xf650f6e9.
//
// Solidity: function send_token_to_paloma(address token, bytes32 receiver, uint256 amount) returns()
func (_Compass *CompassSession) SendTokenToPaloma(token common.Address, receiver [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.SendTokenToPaloma(&_Compass.TransactOpts, token, receiver, amount)
}

// SendTokenToPaloma is a paid mutator transaction binding the contract method 0xf650f6e9.
//
// Solidity: function send_token_to_paloma(address token, bytes32 receiver, uint256 amount) returns()
func (_Compass *CompassTransactorSession) SendTokenToPaloma(token common.Address, receiver [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.SendTokenToPaloma(&_Compass.TransactOpts, token, receiver, amount)
}

// SubmitBatch is a paid mutator transaction binding the contract method 0x14fe978f.
//
// Solidity: function submit_batch(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, address token, (address[],uint256[]) args, uint256 batch_id, uint256 deadline, address relayer, uint256 gas_estimate) returns()
func (_Compass *CompassTransactor) SubmitBatch(opts *bind.TransactOpts, consensus Struct2, token common.Address, args Struct5, batch_id *big.Int, deadline *big.Int, relayer common.Address, gas_estimate *big.Int) (*types.Transaction, error) {
	return _Compass.contract.Transact(opts, "submit_batch", consensus, token, args, batch_id, deadline, relayer, gas_estimate)
}

// SubmitBatch is a paid mutator transaction binding the contract method 0x14fe978f.
//
// Solidity: function submit_batch(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, address token, (address[],uint256[]) args, uint256 batch_id, uint256 deadline, address relayer, uint256 gas_estimate) returns()
func (_Compass *CompassSession) SubmitBatch(consensus Struct2, token common.Address, args Struct5, batch_id *big.Int, deadline *big.Int, relayer common.Address, gas_estimate *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.SubmitBatch(&_Compass.TransactOpts, consensus, token, args, batch_id, deadline, relayer, gas_estimate)
}

// SubmitBatch is a paid mutator transaction binding the contract method 0x14fe978f.
//
// Solidity: function submit_batch(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, address token, (address[],uint256[]) args, uint256 batch_id, uint256 deadline, address relayer, uint256 gas_estimate) returns()
func (_Compass *CompassTransactorSession) SubmitBatch(consensus Struct2, token common.Address, args Struct5, batch_id *big.Int, deadline *big.Int, relayer common.Address, gas_estimate *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.SubmitBatch(&_Compass.TransactOpts, consensus, token, args, batch_id, deadline, relayer, gas_estimate)
}

// SubmitLogicCall is a paid mutator transaction binding the contract method 0xa930e8dc.
//
// Solidity: function submit_logic_call(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, (address,bytes) args, (uint256,uint256,uint256,bytes32) fee_args, uint256 message_id, uint256 deadline, address relayer) returns()
func (_Compass *CompassTransactor) SubmitLogicCall(opts *bind.TransactOpts, consensus Struct2, args Struct3, fee_args Struct4, message_id *big.Int, deadline *big.Int, relayer common.Address) (*types.Transaction, error) {
	return _Compass.contract.Transact(opts, "submit_logic_call", consensus, args, fee_args, message_id, deadline, relayer)
}

// SubmitLogicCall is a paid mutator transaction binding the contract method 0xa930e8dc.
//
// Solidity: function submit_logic_call(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, (address,bytes) args, (uint256,uint256,uint256,bytes32) fee_args, uint256 message_id, uint256 deadline, address relayer) returns()
func (_Compass *CompassSession) SubmitLogicCall(consensus Struct2, args Struct3, fee_args Struct4, message_id *big.Int, deadline *big.Int, relayer common.Address) (*types.Transaction, error) {
	return _Compass.Contract.SubmitLogicCall(&_Compass.TransactOpts, consensus, args, fee_args, message_id, deadline, relayer)
}

// SubmitLogicCall is a paid mutator transaction binding the contract method 0xa930e8dc.
//
// Solidity: function submit_logic_call(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, (address,bytes) args, (uint256,uint256,uint256,bytes32) fee_args, uint256 message_id, uint256 deadline, address relayer) returns()
func (_Compass *CompassTransactorSession) SubmitLogicCall(consensus Struct2, args Struct3, fee_args Struct4, message_id *big.Int, deadline *big.Int, relayer common.Address) (*types.Transaction, error) {
	return _Compass.Contract.SubmitLogicCall(&_Compass.TransactOpts, consensus, args, fee_args, message_id, deadline, relayer)
}

// UpdateCompassAddressInFeeManager is a paid mutator transaction binding the contract method 0x120e8ddd.
//
// Solidity: function update_compass_address_in_fee_manager(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, uint256 deadline, uint256 gas_estimate, address _new_compass, address relayer) returns()
func (_Compass *CompassTransactor) UpdateCompassAddressInFeeManager(opts *bind.TransactOpts, consensus Struct2, deadline *big.Int, gas_estimate *big.Int, _new_compass common.Address, relayer common.Address) (*types.Transaction, error) {
	return _Compass.contract.Transact(opts, "update_compass_address_in_fee_manager", consensus, deadline, gas_estimate, _new_compass, relayer)
}

// UpdateCompassAddressInFeeManager is a paid mutator transaction binding the contract method 0x120e8ddd.
//
// Solidity: function update_compass_address_in_fee_manager(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, uint256 deadline, uint256 gas_estimate, address _new_compass, address relayer) returns()
func (_Compass *CompassSession) UpdateCompassAddressInFeeManager(consensus Struct2, deadline *big.Int, gas_estimate *big.Int, _new_compass common.Address, relayer common.Address) (*types.Transaction, error) {
	return _Compass.Contract.UpdateCompassAddressInFeeManager(&_Compass.TransactOpts, consensus, deadline, gas_estimate, _new_compass, relayer)
}

// UpdateCompassAddressInFeeManager is a paid mutator transaction binding the contract method 0x120e8ddd.
//
// Solidity: function update_compass_address_in_fee_manager(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, uint256 deadline, uint256 gas_estimate, address _new_compass, address relayer) returns()
func (_Compass *CompassTransactorSession) UpdateCompassAddressInFeeManager(consensus Struct2, deadline *big.Int, gas_estimate *big.Int, _new_compass common.Address, relayer common.Address) (*types.Transaction, error) {
	return _Compass.Contract.UpdateCompassAddressInFeeManager(&_Compass.TransactOpts, consensus, deadline, gas_estimate, _new_compass, relayer)
}

// UpdateValset is a paid mutator transaction binding the contract method 0xf064acb2.
//
// Solidity: function update_valset(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, (address[],uint256[],uint256) new_valset, address relayer, uint256 gas_estimate) returns()
func (_Compass *CompassTransactor) UpdateValset(opts *bind.TransactOpts, consensus Struct2, new_valset Struct0, relayer common.Address, gas_estimate *big.Int) (*types.Transaction, error) {
	return _Compass.contract.Transact(opts, "update_valset", consensus, new_valset, relayer, gas_estimate)
}

// UpdateValset is a paid mutator transaction binding the contract method 0xf064acb2.
//
// Solidity: function update_valset(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, (address[],uint256[],uint256) new_valset, address relayer, uint256 gas_estimate) returns()
func (_Compass *CompassSession) UpdateValset(consensus Struct2, new_valset Struct0, relayer common.Address, gas_estimate *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.UpdateValset(&_Compass.TransactOpts, consensus, new_valset, relayer, gas_estimate)
}

// UpdateValset is a paid mutator transaction binding the contract method 0xf064acb2.
//
// Solidity: function update_valset(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, (address[],uint256[],uint256) new_valset, address relayer, uint256 gas_estimate) returns()
func (_Compass *CompassTransactorSession) UpdateValset(consensus Struct2, new_valset Struct0, relayer common.Address, gas_estimate *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.UpdateValset(&_Compass.TransactOpts, consensus, new_valset, relayer, gas_estimate)
}

// Withdraw is a paid mutator transaction binding the contract method 0x048a245d.
//
// Solidity: function withdraw(uint256 amount, address dex, bytes payload, uint256 min_grain) returns()
func (_Compass *CompassTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int, dex common.Address, payload []byte, min_grain *big.Int) (*types.Transaction, error) {
	return _Compass.contract.Transact(opts, "withdraw", amount, dex, payload, min_grain)
}

// Withdraw is a paid mutator transaction binding the contract method 0x048a245d.
//
// Solidity: function withdraw(uint256 amount, address dex, bytes payload, uint256 min_grain) returns()
func (_Compass *CompassSession) Withdraw(amount *big.Int, dex common.Address, payload []byte, min_grain *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.Withdraw(&_Compass.TransactOpts, amount, dex, payload, min_grain)
}

// Withdraw is a paid mutator transaction binding the contract method 0x048a245d.
//
// Solidity: function withdraw(uint256 amount, address dex, bytes payload, uint256 min_grain) returns()
func (_Compass *CompassTransactorSession) Withdraw(amount *big.Int, dex common.Address, payload []byte, min_grain *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.Withdraw(&_Compass.TransactOpts, amount, dex, payload, min_grain)
}

// CompassBatchSendEventIterator is returned from FilterBatchSendEvent and is used to iterate over the raw logs and unpacked data for BatchSendEvent events raised by the Compass contract.
type CompassBatchSendEventIterator struct {
	Event *CompassBatchSendEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CompassBatchSendEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompassBatchSendEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CompassBatchSendEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CompassBatchSendEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompassBatchSendEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompassBatchSendEvent represents a BatchSendEvent event raised by the Compass contract.
type CompassBatchSendEvent struct {
	Token   common.Address
	BatchId *big.Int
	Nonce   *big.Int
	EventId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBatchSendEvent is a free log retrieval operation binding the contract event 0x0ba40544a53f11e70bd7e03a4cfeec841fc3566e81dfbef26f669358a705ad2c.
//
// Solidity: event BatchSendEvent(address token, uint256 batch_id, uint256 nonce, uint256 event_id)
func (_Compass *CompassFilterer) FilterBatchSendEvent(opts *bind.FilterOpts) (*CompassBatchSendEventIterator, error) {
	logs, sub, err := _Compass.contract.FilterLogs(opts, "BatchSendEvent")
	if err != nil {
		return nil, err
	}
	return &CompassBatchSendEventIterator{contract: _Compass.contract, event: "BatchSendEvent", logs: logs, sub: sub}, nil
}

// WatchBatchSendEvent is a free log subscription operation binding the contract event 0x0ba40544a53f11e70bd7e03a4cfeec841fc3566e81dfbef26f669358a705ad2c.
//
// Solidity: event BatchSendEvent(address token, uint256 batch_id, uint256 nonce, uint256 event_id)
func (_Compass *CompassFilterer) WatchBatchSendEvent(opts *bind.WatchOpts, sink chan<- *CompassBatchSendEvent) (event.Subscription, error) {
	logs, sub, err := _Compass.contract.WatchLogs(opts, "BatchSendEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompassBatchSendEvent)
				if err := _Compass.contract.UnpackLog(event, "BatchSendEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBatchSendEvent is a log parse operation binding the contract event 0x0ba40544a53f11e70bd7e03a4cfeec841fc3566e81dfbef26f669358a705ad2c.
//
// Solidity: event BatchSendEvent(address token, uint256 batch_id, uint256 nonce, uint256 event_id)
func (_Compass *CompassFilterer) ParseBatchSendEvent(log types.Log) (*CompassBatchSendEvent, error) {
	event := new(CompassBatchSendEvent)
	if err := _Compass.contract.UnpackLog(event, "BatchSendEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompassERC20DeployedEventIterator is returned from FilterERC20DeployedEvent and is used to iterate over the raw logs and unpacked data for ERC20DeployedEvent events raised by the Compass contract.
type CompassERC20DeployedEventIterator struct {
	Event *CompassERC20DeployedEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CompassERC20DeployedEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompassERC20DeployedEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CompassERC20DeployedEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CompassERC20DeployedEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompassERC20DeployedEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompassERC20DeployedEvent represents a ERC20DeployedEvent event raised by the Compass contract.
type CompassERC20DeployedEvent struct {
	PalomaDenom   string
	TokenContract common.Address
	Name          string
	Symbol        string
	Decimals      uint8
	EventId       *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterERC20DeployedEvent is a free log retrieval operation binding the contract event 0x82fe3a4fa49c6382d0c085746698ddbbafe6c2bf61285b19410644b5b26287c7.
//
// Solidity: event ERC20DeployedEvent(string paloma_denom, address token_contract, string name, string symbol, uint8 decimals, uint256 event_id)
func (_Compass *CompassFilterer) FilterERC20DeployedEvent(opts *bind.FilterOpts) (*CompassERC20DeployedEventIterator, error) {
	logs, sub, err := _Compass.contract.FilterLogs(opts, "ERC20DeployedEvent")
	if err != nil {
		return nil, err
	}
	return &CompassERC20DeployedEventIterator{contract: _Compass.contract, event: "ERC20DeployedEvent", logs: logs, sub: sub}, nil
}

// WatchERC20DeployedEvent is a free log subscription operation binding the contract event 0x82fe3a4fa49c6382d0c085746698ddbbafe6c2bf61285b19410644b5b26287c7.
//
// Solidity: event ERC20DeployedEvent(string paloma_denom, address token_contract, string name, string symbol, uint8 decimals, uint256 event_id)
func (_Compass *CompassFilterer) WatchERC20DeployedEvent(opts *bind.WatchOpts, sink chan<- *CompassERC20DeployedEvent) (event.Subscription, error) {
	logs, sub, err := _Compass.contract.WatchLogs(opts, "ERC20DeployedEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompassERC20DeployedEvent)
				if err := _Compass.contract.UnpackLog(event, "ERC20DeployedEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseERC20DeployedEvent is a log parse operation binding the contract event 0x82fe3a4fa49c6382d0c085746698ddbbafe6c2bf61285b19410644b5b26287c7.
//
// Solidity: event ERC20DeployedEvent(string paloma_denom, address token_contract, string name, string symbol, uint8 decimals, uint256 event_id)
func (_Compass *CompassFilterer) ParseERC20DeployedEvent(log types.Log) (*CompassERC20DeployedEvent, error) {
	event := new(CompassERC20DeployedEvent)
	if err := _Compass.contract.UnpackLog(event, "ERC20DeployedEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompassFundsDepositedEventIterator is returned from FilterFundsDepositedEvent and is used to iterate over the raw logs and unpacked data for FundsDepositedEvent events raised by the Compass contract.
type CompassFundsDepositedEventIterator struct {
	Event *CompassFundsDepositedEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CompassFundsDepositedEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompassFundsDepositedEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CompassFundsDepositedEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CompassFundsDepositedEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompassFundsDepositedEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompassFundsDepositedEvent represents a FundsDepositedEvent event raised by the Compass contract.
type CompassFundsDepositedEvent struct {
	DepositorPalomaAddress [32]byte
	Sender                 common.Address
	Amount                 *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterFundsDepositedEvent is a free log retrieval operation binding the contract event 0x4f78bbd9a86543dc57c484da46f56d43190ac1148b43565fa8d522b1d4fe5298.
//
// Solidity: event FundsDepositedEvent(bytes32 depositor_paloma_address, address sender, uint256 amount)
func (_Compass *CompassFilterer) FilterFundsDepositedEvent(opts *bind.FilterOpts) (*CompassFundsDepositedEventIterator, error) {
	logs, sub, err := _Compass.contract.FilterLogs(opts, "FundsDepositedEvent")
	if err != nil {
		return nil, err
	}
	return &CompassFundsDepositedEventIterator{contract: _Compass.contract, event: "FundsDepositedEvent", logs: logs, sub: sub}, nil
}

// WatchFundsDepositedEvent is a free log subscription operation binding the contract event 0x4f78bbd9a86543dc57c484da46f56d43190ac1148b43565fa8d522b1d4fe5298.
//
// Solidity: event FundsDepositedEvent(bytes32 depositor_paloma_address, address sender, uint256 amount)
func (_Compass *CompassFilterer) WatchFundsDepositedEvent(opts *bind.WatchOpts, sink chan<- *CompassFundsDepositedEvent) (event.Subscription, error) {
	logs, sub, err := _Compass.contract.WatchLogs(opts, "FundsDepositedEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompassFundsDepositedEvent)
				if err := _Compass.contract.UnpackLog(event, "FundsDepositedEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFundsDepositedEvent is a log parse operation binding the contract event 0x4f78bbd9a86543dc57c484da46f56d43190ac1148b43565fa8d522b1d4fe5298.
//
// Solidity: event FundsDepositedEvent(bytes32 depositor_paloma_address, address sender, uint256 amount)
func (_Compass *CompassFilterer) ParseFundsDepositedEvent(log types.Log) (*CompassFundsDepositedEvent, error) {
	event := new(CompassFundsDepositedEvent)
	if err := _Compass.contract.UnpackLog(event, "FundsDepositedEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompassFundsWithdrawnEventIterator is returned from FilterFundsWithdrawnEvent and is used to iterate over the raw logs and unpacked data for FundsWithdrawnEvent events raised by the Compass contract.
type CompassFundsWithdrawnEventIterator struct {
	Event *CompassFundsWithdrawnEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CompassFundsWithdrawnEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompassFundsWithdrawnEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CompassFundsWithdrawnEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CompassFundsWithdrawnEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompassFundsWithdrawnEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompassFundsWithdrawnEvent represents a FundsWithdrawnEvent event raised by the Compass contract.
type CompassFundsWithdrawnEvent struct {
	Receiver common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterFundsWithdrawnEvent is a free log retrieval operation binding the contract event 0xbcd7c5f94d828115734ea3d51400c2e1ad93894d1a5099a1808794a924f71f24.
//
// Solidity: event FundsWithdrawnEvent(address receiver, uint256 amount)
func (_Compass *CompassFilterer) FilterFundsWithdrawnEvent(opts *bind.FilterOpts) (*CompassFundsWithdrawnEventIterator, error) {
	logs, sub, err := _Compass.contract.FilterLogs(opts, "FundsWithdrawnEvent")
	if err != nil {
		return nil, err
	}
	return &CompassFundsWithdrawnEventIterator{contract: _Compass.contract, event: "FundsWithdrawnEvent", logs: logs, sub: sub}, nil
}

// WatchFundsWithdrawnEvent is a free log subscription operation binding the contract event 0xbcd7c5f94d828115734ea3d51400c2e1ad93894d1a5099a1808794a924f71f24.
//
// Solidity: event FundsWithdrawnEvent(address receiver, uint256 amount)
func (_Compass *CompassFilterer) WatchFundsWithdrawnEvent(opts *bind.WatchOpts, sink chan<- *CompassFundsWithdrawnEvent) (event.Subscription, error) {
	logs, sub, err := _Compass.contract.WatchLogs(opts, "FundsWithdrawnEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompassFundsWithdrawnEvent)
				if err := _Compass.contract.UnpackLog(event, "FundsWithdrawnEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFundsWithdrawnEvent is a log parse operation binding the contract event 0xbcd7c5f94d828115734ea3d51400c2e1ad93894d1a5099a1808794a924f71f24.
//
// Solidity: event FundsWithdrawnEvent(address receiver, uint256 amount)
func (_Compass *CompassFilterer) ParseFundsWithdrawnEvent(log types.Log) (*CompassFundsWithdrawnEvent, error) {
	event := new(CompassFundsWithdrawnEvent)
	if err := _Compass.contract.UnpackLog(event, "FundsWithdrawnEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompassLogicCallEventIterator is returned from FilterLogicCallEvent and is used to iterate over the raw logs and unpacked data for LogicCallEvent events raised by the Compass contract.
type CompassLogicCallEventIterator struct {
	Event *CompassLogicCallEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CompassLogicCallEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompassLogicCallEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CompassLogicCallEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CompassLogicCallEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompassLogicCallEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompassLogicCallEvent represents a LogicCallEvent event raised by the Compass contract.
type CompassLogicCallEvent struct {
	LogicContractAddress common.Address
	Payload              []byte
	MessageId            *big.Int
	EventId              *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterLogicCallEvent is a free log retrieval operation binding the contract event 0x0594b174e11e17c2cb4d0d303c2125060bea4f4da113a4e79edce87465592d00.
//
// Solidity: event LogicCallEvent(address logic_contract_address, bytes payload, uint256 message_id, uint256 event_id)
func (_Compass *CompassFilterer) FilterLogicCallEvent(opts *bind.FilterOpts) (*CompassLogicCallEventIterator, error) {
	logs, sub, err := _Compass.contract.FilterLogs(opts, "LogicCallEvent")
	if err != nil {
		return nil, err
	}
	return &CompassLogicCallEventIterator{contract: _Compass.contract, event: "LogicCallEvent", logs: logs, sub: sub}, nil
}

// WatchLogicCallEvent is a free log subscription operation binding the contract event 0x0594b174e11e17c2cb4d0d303c2125060bea4f4da113a4e79edce87465592d00.
//
// Solidity: event LogicCallEvent(address logic_contract_address, bytes payload, uint256 message_id, uint256 event_id)
func (_Compass *CompassFilterer) WatchLogicCallEvent(opts *bind.WatchOpts, sink chan<- *CompassLogicCallEvent) (event.Subscription, error) {
	logs, sub, err := _Compass.contract.WatchLogs(opts, "LogicCallEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompassLogicCallEvent)
				if err := _Compass.contract.UnpackLog(event, "LogicCallEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLogicCallEvent is a log parse operation binding the contract event 0x0594b174e11e17c2cb4d0d303c2125060bea4f4da113a4e79edce87465592d00.
//
// Solidity: event LogicCallEvent(address logic_contract_address, bytes payload, uint256 message_id, uint256 event_id)
func (_Compass *CompassFilterer) ParseLogicCallEvent(log types.Log) (*CompassLogicCallEvent, error) {
	event := new(CompassLogicCallEvent)
	if err := _Compass.contract.UnpackLog(event, "LogicCallEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompassNodeSaleEventIterator is returned from FilterNodeSaleEvent and is used to iterate over the raw logs and unpacked data for NodeSaleEvent events raised by the Compass contract.
type CompassNodeSaleEventIterator struct {
	Event *CompassNodeSaleEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CompassNodeSaleEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompassNodeSaleEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CompassNodeSaleEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CompassNodeSaleEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompassNodeSaleEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompassNodeSaleEvent represents a NodeSaleEvent event raised by the Compass contract.
type CompassNodeSaleEvent struct {
	ContractAddress common.Address
	Buyer           common.Address
	Paloma          [32]byte
	NodeCount       *big.Int
	GrainAmount     *big.Int
	Nonce           *big.Int
	EventId         *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterNodeSaleEvent is a free log retrieval operation binding the contract event 0xc72b917679ae2dea3062a0a355d542c92296f3e5c39cfcb0db7c30e28c816349.
//
// Solidity: event NodeSaleEvent(address contract_address, address buyer, bytes32 paloma, uint256 node_count, uint256 grain_amount, uint256 nonce, uint256 event_id)
func (_Compass *CompassFilterer) FilterNodeSaleEvent(opts *bind.FilterOpts) (*CompassNodeSaleEventIterator, error) {
	logs, sub, err := _Compass.contract.FilterLogs(opts, "NodeSaleEvent")
	if err != nil {
		return nil, err
	}
	return &CompassNodeSaleEventIterator{contract: _Compass.contract, event: "NodeSaleEvent", logs: logs, sub: sub}, nil
}

// WatchNodeSaleEvent is a free log subscription operation binding the contract event 0xc72b917679ae2dea3062a0a355d542c92296f3e5c39cfcb0db7c30e28c816349.
//
// Solidity: event NodeSaleEvent(address contract_address, address buyer, bytes32 paloma, uint256 node_count, uint256 grain_amount, uint256 nonce, uint256 event_id)
func (_Compass *CompassFilterer) WatchNodeSaleEvent(opts *bind.WatchOpts, sink chan<- *CompassNodeSaleEvent) (event.Subscription, error) {
	logs, sub, err := _Compass.contract.WatchLogs(opts, "NodeSaleEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompassNodeSaleEvent)
				if err := _Compass.contract.UnpackLog(event, "NodeSaleEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNodeSaleEvent is a log parse operation binding the contract event 0xc72b917679ae2dea3062a0a355d542c92296f3e5c39cfcb0db7c30e28c816349.
//
// Solidity: event NodeSaleEvent(address contract_address, address buyer, bytes32 paloma, uint256 node_count, uint256 grain_amount, uint256 nonce, uint256 event_id)
func (_Compass *CompassFilterer) ParseNodeSaleEvent(log types.Log) (*CompassNodeSaleEvent, error) {
	event := new(CompassNodeSaleEvent)
	if err := _Compass.contract.UnpackLog(event, "NodeSaleEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompassSendToPalomaEventIterator is returned from FilterSendToPalomaEvent and is used to iterate over the raw logs and unpacked data for SendToPalomaEvent events raised by the Compass contract.
type CompassSendToPalomaEventIterator struct {
	Event *CompassSendToPalomaEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CompassSendToPalomaEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompassSendToPalomaEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CompassSendToPalomaEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CompassSendToPalomaEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompassSendToPalomaEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompassSendToPalomaEvent represents a SendToPalomaEvent event raised by the Compass contract.
type CompassSendToPalomaEvent struct {
	Token    common.Address
	Sender   common.Address
	Receiver [32]byte
	Amount   *big.Int
	Nonce    *big.Int
	EventId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSendToPalomaEvent is a free log retrieval operation binding the contract event 0xc5bdbcfcfae5c5b84c56bbf0860c4286d627aefaf28ce4011ba4fcb9b5aadf08.
//
// Solidity: event SendToPalomaEvent(address token, address sender, bytes32 receiver, uint256 amount, uint256 nonce, uint256 event_id)
func (_Compass *CompassFilterer) FilterSendToPalomaEvent(opts *bind.FilterOpts) (*CompassSendToPalomaEventIterator, error) {
	logs, sub, err := _Compass.contract.FilterLogs(opts, "SendToPalomaEvent")
	if err != nil {
		return nil, err
	}
	return &CompassSendToPalomaEventIterator{contract: _Compass.contract, event: "SendToPalomaEvent", logs: logs, sub: sub}, nil
}

// WatchSendToPalomaEvent is a free log subscription operation binding the contract event 0xc5bdbcfcfae5c5b84c56bbf0860c4286d627aefaf28ce4011ba4fcb9b5aadf08.
//
// Solidity: event SendToPalomaEvent(address token, address sender, bytes32 receiver, uint256 amount, uint256 nonce, uint256 event_id)
func (_Compass *CompassFilterer) WatchSendToPalomaEvent(opts *bind.WatchOpts, sink chan<- *CompassSendToPalomaEvent) (event.Subscription, error) {
	logs, sub, err := _Compass.contract.WatchLogs(opts, "SendToPalomaEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompassSendToPalomaEvent)
				if err := _Compass.contract.UnpackLog(event, "SendToPalomaEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSendToPalomaEvent is a log parse operation binding the contract event 0xc5bdbcfcfae5c5b84c56bbf0860c4286d627aefaf28ce4011ba4fcb9b5aadf08.
//
// Solidity: event SendToPalomaEvent(address token, address sender, bytes32 receiver, uint256 amount, uint256 nonce, uint256 event_id)
func (_Compass *CompassFilterer) ParseSendToPalomaEvent(log types.Log) (*CompassSendToPalomaEvent, error) {
	event := new(CompassSendToPalomaEvent)
	if err := _Compass.contract.UnpackLog(event, "SendToPalomaEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompassUpdateCompassAddressInFeeManagerIterator is returned from FilterUpdateCompassAddressInFeeManager and is used to iterate over the raw logs and unpacked data for UpdateCompassAddressInFeeManager events raised by the Compass contract.
type CompassUpdateCompassAddressInFeeManagerIterator struct {
	Event *CompassUpdateCompassAddressInFeeManager // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CompassUpdateCompassAddressInFeeManagerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompassUpdateCompassAddressInFeeManager)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CompassUpdateCompassAddressInFeeManager)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CompassUpdateCompassAddressInFeeManagerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompassUpdateCompassAddressInFeeManagerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompassUpdateCompassAddressInFeeManager represents a UpdateCompassAddressInFeeManager event raised by the Compass contract.
type CompassUpdateCompassAddressInFeeManager struct {
	NewCompass common.Address
	EventId    *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpdateCompassAddressInFeeManager is a free log retrieval operation binding the contract event 0x0662fa994447aee93c78fd46a6801f664461ecbba2af6c5a0a6aecbc955835fc.
//
// Solidity: event UpdateCompassAddressInFeeManager(address new_compass, uint256 event_id)
func (_Compass *CompassFilterer) FilterUpdateCompassAddressInFeeManager(opts *bind.FilterOpts) (*CompassUpdateCompassAddressInFeeManagerIterator, error) {
	logs, sub, err := _Compass.contract.FilterLogs(opts, "UpdateCompassAddressInFeeManager")
	if err != nil {
		return nil, err
	}
	return &CompassUpdateCompassAddressInFeeManagerIterator{contract: _Compass.contract, event: "UpdateCompassAddressInFeeManager", logs: logs, sub: sub}, nil
}

// WatchUpdateCompassAddressInFeeManager is a free log subscription operation binding the contract event 0x0662fa994447aee93c78fd46a6801f664461ecbba2af6c5a0a6aecbc955835fc.
//
// Solidity: event UpdateCompassAddressInFeeManager(address new_compass, uint256 event_id)
func (_Compass *CompassFilterer) WatchUpdateCompassAddressInFeeManager(opts *bind.WatchOpts, sink chan<- *CompassUpdateCompassAddressInFeeManager) (event.Subscription, error) {
	logs, sub, err := _Compass.contract.WatchLogs(opts, "UpdateCompassAddressInFeeManager")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompassUpdateCompassAddressInFeeManager)
				if err := _Compass.contract.UnpackLog(event, "UpdateCompassAddressInFeeManager", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdateCompassAddressInFeeManager is a log parse operation binding the contract event 0x0662fa994447aee93c78fd46a6801f664461ecbba2af6c5a0a6aecbc955835fc.
//
// Solidity: event UpdateCompassAddressInFeeManager(address new_compass, uint256 event_id)
func (_Compass *CompassFilterer) ParseUpdateCompassAddressInFeeManager(log types.Log) (*CompassUpdateCompassAddressInFeeManager, error) {
	event := new(CompassUpdateCompassAddressInFeeManager)
	if err := _Compass.contract.UnpackLog(event, "UpdateCompassAddressInFeeManager", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompassValsetUpdatedIterator is returned from FilterValsetUpdated and is used to iterate over the raw logs and unpacked data for ValsetUpdated events raised by the Compass contract.
type CompassValsetUpdatedIterator struct {
	Event *CompassValsetUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CompassValsetUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompassValsetUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CompassValsetUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CompassValsetUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompassValsetUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompassValsetUpdated represents a ValsetUpdated event raised by the Compass contract.
type CompassValsetUpdated struct {
	Checkpoint [32]byte
	ValsetId   *big.Int
	EventId    *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterValsetUpdated is a free log retrieval operation binding the contract event 0xb7ca5e46e360950244488bf096bf742a1f63183cf1ee5b3b0c53045b6247bf5b.
//
// Solidity: event ValsetUpdated(bytes32 checkpoint, uint256 valset_id, uint256 event_id)
func (_Compass *CompassFilterer) FilterValsetUpdated(opts *bind.FilterOpts) (*CompassValsetUpdatedIterator, error) {
	logs, sub, err := _Compass.contract.FilterLogs(opts, "ValsetUpdated")
	if err != nil {
		return nil, err
	}
	return &CompassValsetUpdatedIterator{contract: _Compass.contract, event: "ValsetUpdated", logs: logs, sub: sub}, nil
}

// WatchValsetUpdated is a free log subscription operation binding the contract event 0xb7ca5e46e360950244488bf096bf742a1f63183cf1ee5b3b0c53045b6247bf5b.
//
// Solidity: event ValsetUpdated(bytes32 checkpoint, uint256 valset_id, uint256 event_id)
func (_Compass *CompassFilterer) WatchValsetUpdated(opts *bind.WatchOpts, sink chan<- *CompassValsetUpdated) (event.Subscription, error) {
	logs, sub, err := _Compass.contract.WatchLogs(opts, "ValsetUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompassValsetUpdated)
				if err := _Compass.contract.UnpackLog(event, "ValsetUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValsetUpdated is a log parse operation binding the contract event 0xb7ca5e46e360950244488bf096bf742a1f63183cf1ee5b3b0c53045b6247bf5b.
//
// Solidity: event ValsetUpdated(bytes32 checkpoint, uint256 valset_id, uint256 event_id)
func (_Compass *CompassFilterer) ParseValsetUpdated(log types.Log) (*CompassValsetUpdated, error) {
	event := new(CompassValsetUpdated)
	if err := _Compass.contract.UnpackLog(event, "ValsetUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
