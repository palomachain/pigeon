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

// CompassMetaData contains all meta data concerning the Compass contract.
var CompassMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"checkpoint\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"valset_id\",\"type\":\"uint256\"}],\"name\":\"ValsetUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"logic_contract_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"message_id\",\"type\":\"uint256\"}],\"name\":\"LogicCallEvent\",\"type\":\"event\"},{\"inputs\":[{\"name\":\"turnstone_id\",\"type\":\"bytes32\"},{\"components\":[{\"name\":\"validators\",\"type\":\"address[]\"},{\"name\":\"powers\",\"type\":\"uint256[]\"},{\"name\":\"valset_id\",\"type\":\"uint256\"}],\"name\":\"valset\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"turnstone_id\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"name\":\"validators\",\"type\":\"address[]\"},{\"name\":\"powers\",\"type\":\"uint256[]\"},{\"name\":\"valset_id\",\"type\":\"uint256\"}],\"name\":\"valset\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"v\",\"type\":\"uint256\"},{\"name\":\"r\",\"type\":\"uint256\"},{\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"signatures\",\"type\":\"tuple[]\"}],\"name\":\"consensus\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"validators\",\"type\":\"address[]\"},{\"name\":\"powers\",\"type\":\"uint256[]\"},{\"name\":\"valset_id\",\"type\":\"uint256\"}],\"name\":\"new_valset\",\"type\":\"tuple\"}],\"name\":\"update_valset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"name\":\"validators\",\"type\":\"address[]\"},{\"name\":\"powers\",\"type\":\"uint256[]\"},{\"name\":\"valset_id\",\"type\":\"uint256\"}],\"name\":\"valset\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"v\",\"type\":\"uint256\"},{\"name\":\"r\",\"type\":\"uint256\"},{\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"signatures\",\"type\":\"tuple[]\"}],\"name\":\"consensus\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"logic_contract_address\",\"type\":\"address\"},{\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"args\",\"type\":\"tuple\"},{\"name\":\"message_id\",\"type\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"submit_logic_call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"last_checkpoint\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"last_valset_id\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"arg0\",\"type\":\"uint256\"}],\"name\":\"message_id_used\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]'0x60206111126000396000516020816110f20160003960005181016101406020826110f201600039600051116110ed576020816110f2016000396000518060405260008161014081116110ed57801561008357905b60208160051b60208601016110f2016000396000518060a01c6110ed578160051b60600152600101818118610053575b505050506020602082016110f20160003960005181016101406020826110f201600039600051116110ed576020816110f2016000396000518061286052602082018160051b80826110f20161288039505050506020604082016110f2016000396000516150805250346110ed5760206110f260003960005161a200526040366150a037600060405161014081116110ed57801561018857905b8060051b606001516150e0526150a0516150c051612860518110156110ed5760051b61288001518082018281106110ed57905090506150a05263aaaaaaaa6150a0511061016857610188565b6150c051600181018181106110ed5790506150c05260010181811861011c575b505063aaaaaaaa6150a05110156101ff5760126150e0527f496e73756666696369656e7420506f7765720000000000000000000000000000615100526150e0506150e0518061510001601f826000031636823750506308c379a06150a05260206150c052601f19601f6150e05101166044016150bcfd5b63299018c261510452600460808061512452806151240160006040518083528060051b60008261014081116110ed57801561025357905b8060051b606001518160051b602088010152600101818118610236575b50508201602001915050905081019050806151445280615124016000612860518083528060051b60008261014081116110ed5780156102ac57905b8060051b61288001518160051b60208801015260010181811861028e575b50508201602001915050905081019050615080516151645260206110f26000396000516151845201615100526151008051602082012090506150e0526150e051600055615080516001557f09d40458cf931745f8d532ef13fa9c74bfb7fe0edcee88e0a677b0cbef88f0f96150e0516151005261508051615120526040615100a1610dad61033f61945339610dcd619453f36003361161000c57610a38565b60003560e01c34610d9b5763f0f40504811861003e5760043610610d9b576020610dad60003960005160405260206040f35b63eadf4af78118610518576101e43610610d9b576004356004018035810180358101610140813511610d9b5780358061cb40526000816101408111610d9b5780156100ab57905b8060051b6020850101358060a01c610d9b578160051b61cb600152600101818118610085575b5050505060208101358101610140813511610d9b5780358061f36052602082018160051b808261f3803750505050604081013562011b80525060208101358101610140813511610d9b5780358062011ba0526020820160608202808262011bc037505050505060243560040180358101610140813511610d9b57803580620193c0526000816101408111610d9b57801561016857905b8060051b6020850101358060a01c610d9b578160051b620193e00152600101818118610141575b5050505060208101358101610140813511610d9b578035806201bbe052602082018160051b80826201bc00375050505060408101356201e400525062011b80516201e40051116102215760116201e420527f496e76616c69642056616c7365742049440000000000000000000000000000006201e440526201e420506201e42051806201e44001601f826000031636823750506308c379a06201e3e05260206201e40052601f19601f6201e4205101166044016201e3fcfd5b6040366201e420376000620193c0516101408111610d9b5780156102b857905b8060051b620193e001516201e460526201e420516201e440516201bbe051811015610d9b5760051b6201bc000151808201828110610d9b57905090506201e4205263aaaaaaaa6201e4205110610296576102b8565b6201e4405160018101818110610d9b5790506201e44052600101818118610241575b505063aaaaaaaa6201e4205110156103395760126201e460527f496e73756666696369656e7420506f77657200000000000000000000000000006201e480526201e460506201e46051806201e48001601f826000031636823750506308c379a06201e4205260206201e44052601f19601f6201e4605101166044016201e43cfd5b61cb4051806040528060051b8060608261cb6060045afa50505061f3605180612860528060051b806128808261f38060045afa50505062011b8051615080526103846201e460610cb1565b6201e4605160005418156104015760146201e480527f496e636f727265637420436865636b706f696e740000000000000000000000006201e4a0526201e480506201e48051806201e4a001601f826000031636823750506308c379a06201e4405260206201e46052601f19601f6201e4805101166044016201e45cfd5b620193c051806040528060051b80606082620193e060045afa5050506201bbe05180612860528060051b80612880826201bc0060045afa5050506201e40051615080526104506201e480610cb1565b6201e480516201e4605261cb4051806101a0528060051b806101c08261cb6060045afa50505061f36051806129c0528060051b806129e08261f38060045afa50505062011b80516151e05262011ba051806152005260608102806152208262011bc060045afa5050506201e4605161ca20526104ca610ae1565b6201e460516000556201e400516001557f09d40458cf931745f8d532ef13fa9c74bfb7fe0edcee88e0a677b0cbef88f0f96201e460516201e480526201e400516201e4a05260406201e480a1005b631029ae6f81186109cb576101e43610610d9b576004356004018035810180358101610140813511610d9b5780358061cb40526000816101408111610d9b57801561058557905b8060051b6020850101358060a01c610d9b578160051b61cb60015260010181811861055f575b5050505060208101358101610140813511610d9b5780358061f36052602082018160051b808261f3803750505050604081013562011b80525060208101358101610140813511610d9b5780358062011ba0526020820160608202808262011bc037505050505060243560040180358060a01c610d9b57620193c05260208101358101615000813511610d9b57803580620193e0526020820181816201940037505050506064354211156106a15760076201e400527f54696d656f7574000000000000000000000000000000000000000000000000006201e420526201e400506201e40051806201e42001601f826000031636823750506308c379a06201e3c05260206201e3e052601f19601f6201e4005101166044016201e3dcfd5b60026044356020526000526040600020541561072657600f6201e400527f55736564204d6573736167655f494400000000000000000000000000000000006201e420526201e400506201e40051806201e42001601f826000031636823750506308c379a06201e3c05260206201e3e052601f19601f6201e4005101166044016201e3dcfd5b6001600260443560205260005260406000205561cb4051806040528060051b8060608261cb6060045afa50505061f3605180612860528060051b806128808261f38060045afa50505062011b8051615080526107846201e400610cb1565b6201e4005160005418156108015760146201e420527f496e636f727265637420436865636b706f696e740000000000000000000000006201e440526201e420506201e42051806201e44001601f826000031636823750506308c379a06201e3e05260206201e40052601f19601f6201e4205101166044016201e3fcfd5b63980721b26201e4245260046080806201e44452806201e444016040620193c0518252806020830152808201620193e051808252602082018181836201940060045afa5050508051806020830101601f82600003163682375050601f19601f825160200101169050810190509050810190506044356201e464526020610dad6000396000516201e484526064356201e4a452016201e420526201e4208051602082012090506201e4005261cb4051806101a0528060051b806101c08261cb6060045afa50505061f36051806129c0528060051b806129e08261f38060045afa50505062011b80516151e05262011ba051806152005260608102806152208262011bc060045afa5050506201e4005161ca205261091b610ae1565b620193e0600060008251602084016000620193c0515af19050610943573d600060003e3d6000fd5b7f0d2bd340033bb64fd086788e6685b480a9bf10b98d63e9b8073eb5d0bd6c6ee96060620193c0516201e42052806201e44052806201e42001620193e051808252602082018181836201940060045afa5050508051806020830101601f82600003163682375050601f19601f825160200101169050810190506044356201e460526201e420a1005b63a9a4a98381186109ea5760043610610d9b5760005460405260206040f35b634da6ecc98118610a095760043610610d9b5760015460405260206040f35b6338d6172d8118610a365760243610610d9b57600260043560205260005260406000205460405260206040f35b505b60006000fd5b6000601c610100527f19457468657265756d205369676e6564204d6573736167653a0a333200000000610120526101008051602082018361016001815181525050808301925050506060518161016001526020810190508061014052610140905080516020820120905060e05260e051610100526080516101205260a0516101405260c0516101605260206000608061010060015afa5060005160405114815250565b60403661ca40376000615200516101408111610d9b578015610c3857905b6060810261522001805161ca8052602081015161caa052604081015161cac0525061ca805115610c185761ca40516101a051811015610d9b5760051b6101c0015160405261ca205160605261ca805160805261caa05160a05261cac05160c052610b6a61cae0610a3e565b61cae051610bd857601161cb00527f496e76616c6964205369676e617475726500000000000000000000000000000061cb205261cb005061cb00518061cb2001601f826000031636823750506308c379a061cac052602061cae052601f19601f61cb0051011660440161cadcfd5b61ca605161ca40516129c051811015610d9b5760051b6129e00151808201828110610d9b579050905061ca605263aaaaaaaa61ca605110610c1857610c38565b61ca405160018101818110610d9b57905061ca4052600101818118610aff575b505063aaaaaaaa61ca60511015610caf57601261ca80527f496e73756666696369656e7420506f776572000000000000000000000000000061caa05261ca805061ca80518061caa001601f826000031636823750506308c379a061ca4052602061ca6052601f19601f61ca8051011660440161ca5cfd5b565b63299018c26150a45260046080806150c452806150c40160006040518083528060051b6000826101408111610d9b578015610d0557905b8060051b606001518160051b602088010152600101818118610ce8575b50508201602001915050905081019050806150e452806150c4016000612860518083528060051b6000826101408111610d9b578015610d5e57905b8060051b61288001518160051b602088010152600101818118610d40575b5050820160200191505090508101905061508051615104526020610dad60003960005161512452016150a0526150a0805160208201209050815250565b600080fda165767970657283000307000b005b600080fd",
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

// TurnstoneId is a free data retrieval call binding the contract method 0xf0f40504.
//
// Solidity: function turnstone_id() pure returns(bytes32)
func (_Compass *CompassCaller) TurnstoneId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Compass.contract.Call(opts, &out, "turnstone_id")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// TurnstoneId is a free data retrieval call binding the contract method 0xf0f40504.
//
// Solidity: function turnstone_id() pure returns(bytes32)
func (_Compass *CompassSession) TurnstoneId() ([32]byte, error) {
	return _Compass.Contract.TurnstoneId(&_Compass.CallOpts)
}

// TurnstoneId is a free data retrieval call binding the contract method 0xf0f40504.
//
// Solidity: function turnstone_id() pure returns(bytes32)
func (_Compass *CompassCallerSession) TurnstoneId() ([32]byte, error) {
	return _Compass.Contract.TurnstoneId(&_Compass.CallOpts)
}

// SubmitLogicCall is a paid mutator transaction binding the contract method 0x1029ae6f.
//
// Solidity: function submit_logic_call(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, (address,bytes) args, uint256 message_id, uint256 deadline) returns()
func (_Compass *CompassTransactor) SubmitLogicCall(opts *bind.TransactOpts, consensus Struct2, args Struct3, message_id *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _Compass.contract.Transact(opts, "submit_logic_call", consensus, args, message_id, deadline)
}

// SubmitLogicCall is a paid mutator transaction binding the contract method 0x1029ae6f.
//
// Solidity: function submit_logic_call(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, (address,bytes) args, uint256 message_id, uint256 deadline) returns()
func (_Compass *CompassSession) SubmitLogicCall(consensus Struct2, args Struct3, message_id *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.SubmitLogicCall(&_Compass.TransactOpts, consensus, args, message_id, deadline)
}

// SubmitLogicCall is a paid mutator transaction binding the contract method 0x1029ae6f.
//
// Solidity: function submit_logic_call(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, (address,bytes) args, uint256 message_id, uint256 deadline) returns()
func (_Compass *CompassTransactorSession) SubmitLogicCall(consensus Struct2, args Struct3, message_id *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _Compass.Contract.SubmitLogicCall(&_Compass.TransactOpts, consensus, args, message_id, deadline)
}

// UpdateValset is a paid mutator transaction binding the contract method 0xeadf4af7.
//
// Solidity: function update_valset(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, (address[],uint256[],uint256) new_valset) returns()
func (_Compass *CompassTransactor) UpdateValset(opts *bind.TransactOpts, consensus Struct2, new_valset Struct0) (*types.Transaction, error) {
	return _Compass.contract.Transact(opts, "update_valset", consensus, new_valset)
}

// UpdateValset is a paid mutator transaction binding the contract method 0xeadf4af7.
//
// Solidity: function update_valset(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, (address[],uint256[],uint256) new_valset) returns()
func (_Compass *CompassSession) UpdateValset(consensus Struct2, new_valset Struct0) (*types.Transaction, error) {
	return _Compass.Contract.UpdateValset(&_Compass.TransactOpts, consensus, new_valset)
}

// UpdateValset is a paid mutator transaction binding the contract method 0xeadf4af7.
//
// Solidity: function update_valset(((address[],uint256[],uint256),(uint256,uint256,uint256)[]) consensus, (address[],uint256[],uint256) new_valset) returns()
func (_Compass *CompassTransactorSession) UpdateValset(consensus Struct2, new_valset Struct0) (*types.Transaction, error) {
	return _Compass.Contract.UpdateValset(&_Compass.TransactOpts, consensus, new_valset)
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
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterLogicCallEvent is a free log retrieval operation binding the contract event 0x0d2bd340033bb64fd086788e6685b480a9bf10b98d63e9b8073eb5d0bd6c6ee9.
//
// Solidity: event LogicCallEvent(address logic_contract_address, bytes payload, uint256 message_id)
func (_Compass *CompassFilterer) FilterLogicCallEvent(opts *bind.FilterOpts) (*CompassLogicCallEventIterator, error) {

	logs, sub, err := _Compass.contract.FilterLogs(opts, "LogicCallEvent")
	if err != nil {
		return nil, err
	}
	return &CompassLogicCallEventIterator{contract: _Compass.contract, event: "LogicCallEvent", logs: logs, sub: sub}, nil
}

// WatchLogicCallEvent is a free log subscription operation binding the contract event 0x0d2bd340033bb64fd086788e6685b480a9bf10b98d63e9b8073eb5d0bd6c6ee9.
//
// Solidity: event LogicCallEvent(address logic_contract_address, bytes payload, uint256 message_id)
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

// ParseLogicCallEvent is a log parse operation binding the contract event 0x0d2bd340033bb64fd086788e6685b480a9bf10b98d63e9b8073eb5d0bd6c6ee9.
//
// Solidity: event LogicCallEvent(address logic_contract_address, bytes payload, uint256 message_id)
func (_Compass *CompassFilterer) ParseLogicCallEvent(log types.Log) (*CompassLogicCallEvent, error) {
	event := new(CompassLogicCallEvent)
	if err := _Compass.contract.UnpackLog(event, "LogicCallEvent", log); err != nil {
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
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterValsetUpdated is a free log retrieval operation binding the contract event 0x09d40458cf931745f8d532ef13fa9c74bfb7fe0edcee88e0a677b0cbef88f0f9.
//
// Solidity: event ValsetUpdated(bytes32 checkpoint, uint256 valset_id)
func (_Compass *CompassFilterer) FilterValsetUpdated(opts *bind.FilterOpts) (*CompassValsetUpdatedIterator, error) {

	logs, sub, err := _Compass.contract.FilterLogs(opts, "ValsetUpdated")
	if err != nil {
		return nil, err
	}
	return &CompassValsetUpdatedIterator{contract: _Compass.contract, event: "ValsetUpdated", logs: logs, sub: sub}, nil
}

// WatchValsetUpdated is a free log subscription operation binding the contract event 0x09d40458cf931745f8d532ef13fa9c74bfb7fe0edcee88e0a677b0cbef88f0f9.
//
// Solidity: event ValsetUpdated(bytes32 checkpoint, uint256 valset_id)
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

// ParseValsetUpdated is a log parse operation binding the contract event 0x09d40458cf931745f8d532ef13fa9c74bfb7fe0edcee88e0a677b0cbef88f0f9.
//
// Solidity: event ValsetUpdated(bytes32 checkpoint, uint256 valset_id)
func (_Compass *CompassFilterer) ParseValsetUpdated(log types.Log) (*CompassValsetUpdated, error) {
	event := new(CompassValsetUpdated)
	if err := _Compass.contract.UnpackLog(event, "ValsetUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
