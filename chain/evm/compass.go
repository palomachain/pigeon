package evm

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"time"

	"cosmossdk.io/math"
	"github.com/VolumeFi/whoops"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	evmtypes "github.com/palomachain/paloma/v2/x/evm/types"
	skywaytypes "github.com/palomachain/paloma/v2/x/skyway/types"
	"github.com/palomachain/pigeon/chain"
	cabi "github.com/palomachain/pigeon/chain/evm/abi/compass"
	"github.com/palomachain/pigeon/internal/ethfilter"
	"github.com/palomachain/pigeon/internal/liblog"
	"github.com/palomachain/pigeon/util/slice"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const (
	SignedMessagePrefix                    = "\x19Ethereum Signed Message:\n32"
	cEventQueryBlockHeightMinWindow        = 10
	cConservativeDummyGasEstimate   uint64 = 300_000
)

var (
	// Node sale event fields:
	//   contract_addr, buyer_addr, paloma_addr, node_count, grain_amount,
	//   skyway_nonce, event_id
	nodeSaleEvent = crypto.Keccak256Hash([]byte(
		"NodeSaleEvent(address,address,bytes32,uint256,uint256,uint256,uint256)",
	))
	batchSendEvent = crypto.Keccak256Hash([]byte(
		"BatchSendEvent(address,uint256,uint256,uint256)",
	))
	sendToPalomaEvent = crypto.Keccak256Hash([]byte(
		"SendToPalomaEvent(address,address,bytes32,uint256,uint256,uint256)",
	))
)

var errValsetIDMismatch = errors.New("valset id mismatch")

//go:generate mockery --name=evmClienter --inpackage --testonly
type evmClienter interface {
	FilterLogs(ctx context.Context, fq ethereum.FilterQuery, currBlockHeight *big.Int, fn func(logs []ethtypes.Log) bool) (bool, error)
	ExecuteSmartContract(ctx context.Context, chainID *big.Int, contractAbi abi.ABI, addr common.Address, opts callOptions, method string, arguments []any) (*ethtypes.Transaction, error)
	DeployContract(ctx context.Context, chainID *big.Int, rawABI string, bytecode, constructorInput []byte) (contractAddr common.Address, tx *ethtypes.Transaction, err error)
	TransactionByHash(ctx context.Context, txHash common.Hash) (*ethtypes.Transaction, bool, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*ethtypes.Receipt, error)

	BalanceAt(ctx context.Context, address common.Address, blockHeight uint64) (*big.Int, error)
	FindBlockNearestToTime(ctx context.Context, startingHeight uint64, when time.Time) (uint64, error)
	FindCurrentBlockNumber(ctx context.Context) (*big.Int, error)
	LastValsetID(ctx context.Context, addr common.Address) (*big.Int, error)
	QueryUserFunds(ctx context.Context, feemgraddr common.Address, palomaAddress [32]byte) (*big.Int, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	GetEthClient() ethClientConn
}

type compass struct {
	paloma                  PalomaClienter
	evm                     evmClienter
	compassAbi              *abi.ABI
	chainID                 *big.Int
	CompassID               string
	ChainReferenceID        string
	lastObservedBlockHeight uint64
	startingBlockHeight     int64
	smartContractAddr       common.Address
	feeMgrContractAddr      common.Address
}

func newCompassClient(
	smartContractAddrStr,
	feeMgrContractAddrStr,
	compassID,
	chainReferenceID string,
	chainID *big.Int,
	compassAbi *abi.ABI,
	paloma PalomaClienter,
	evm evmClienter,
) compass {
	// if !ethcommon.IsHexAddress(smartContractAddrStr) {
	// 	whoops.Assert(errors.Unrecoverable(ErrInvalidAddress.Format(smartContractAddrStr)))
	// }
	return compass{
		CompassID:          compassID,
		ChainReferenceID:   chainReferenceID,
		smartContractAddr:  common.HexToAddress(smartContractAddrStr),
		feeMgrContractAddr: common.HexToAddress(feeMgrContractAddrStr),
		chainID:            chainID,
		compassAbi:         compassAbi,
		paloma:             paloma,
		evm:                evm,
	}
}

type Signature struct {
	V *big.Int
	R *big.Int
	S *big.Int
}
type CompassValset struct {
	ValsetId   *big.Int
	Validators []common.Address
	Powers     []*big.Int
}
type CompassConsensus struct {
	Valset     CompassValset
	Signatures []Signature

	originalSignatures [][]byte
}

type CompassLogicCallArgs struct {
	Payload              []byte
	LogicContractAddress common.Address
}

type FeeArgs cabi.Struct5

type CompassTokenSendArgs struct {
	Receiver []common.Address
	Amount   []*big.Int
}

func (c CompassConsensus) OriginalSignatures() [][]byte {
	return c.originalSignatures
}

func (t compass) updateValset(
	ctx context.Context,
	queueTypeName string,
	newValset *evmtypes.Valset,
	origMessage chain.MessageWithSignatures,
	ethSender common.Address,
	estimate *big.Int,
	opts callOptions,
) (*ethtypes.Transaction, uint64, error) {
	currentValsetID, err := t.findLastValsetMessageID(ctx)
	if err != nil {
		return nil, 0, err
	}
	logger := liblog.WithContext(ctx).WithFields(log.Fields{
		"chain-reference-id": t.ChainReferenceID,
		"current-valset-id":  currentValsetID,
	})
	logger.Debug("update_valset")

	currentValset, err := t.paloma.QueryGetEVMValsetByID(ctx, currentValsetID, t.ChainReferenceID)
	if err != nil {
		return nil, 0, err
	}
	if currentValset == nil {
		logger.Error("current valset is empty")
		return nil, 0, fmt.Errorf("current valset is empty")
	}

	consensusReached := isConsensusReached(ctx, currentValset, origMessage)
	if !consensusReached {
		return nil, 0, ErrNoConsensus
	}

	if opts.estimateOnly {
		// Simulate maximum gas estimate to ensure the transaction is not rejected
		estimate = big.NewInt(0).SetUint64(cConservativeDummyGasEstimate)
	}

	// TODO: Use generated contract code directly
	// compass 2.0.0
	// def update_valset(consensus: Consensus, new_valset: Valset, relayer: address, gas_estimate: uint256)
	tx, err := t.callCompass(ctx, opts, "update_valset", []any{
		BuildCompassConsensus(currentValset, origMessage.Signatures),
		TransformValsetToCompassValset(newValset),
		ethSender,
		estimate,
	})
	if err != nil {
		return nil, 0, err
	}

	return tx, currentValsetID, nil
}

func (t compass) submitLogicCall(
	ctx context.Context,
	queueTypeName string,
	msg *evmtypes.SubmitLogicCall,
	origMessage chain.MessageWithSignatures,
	ethSender common.Address,
	opts callOptions,
) (*ethtypes.Transaction, uint64, error) {
	logger := liblog.WithContext(ctx).WithFields(log.Fields{
		"chain-id": t.ChainReferenceID,
		"msg-id":   origMessage.ID,
	})
	logger.Info("submit logic call")

	// Skip already executed check in case of estimate only
	if !opts.estimateOnly {
		executed, err := t.isArbitraryCallAlreadyExecuted(ctx, origMessage.ID)
		if err != nil {
			return nil, 0, err
		}
		if executed {
			return nil, 0, ErrCallAlreadyExecuted
		}
	}

	valsetID, err := t.performValsetIDCrosscheck(ctx, t.ChainReferenceID)
	logger = logger.WithField("last-valset-id", valsetID)
	if err != nil {
		if errors.Is(err, errValsetIDMismatch) {
			logger.Warn("Valset ID mismatch. Swallowing error to retry message...")
			return nil, 0, nil
		}

		return nil, 0, err
	}

	valset, err := t.paloma.QueryGetEVMValsetByID(ctx, valsetID, t.ChainReferenceID)
	if err != nil {
		return nil, 0, err
	}

	consensusReached := isConsensusReached(ctx, valset, origMessage)
	if !consensusReached {
		return nil, 0, ErrNoConsensus
	}

	con := BuildCompassConsensus(valset, origMessage.Signatures)
	compassArgs := CompassLogicCallArgs{
		LogicContractAddress: common.HexToAddress(msg.GetHexContractAddress()),
		Payload:              msg.GetPayload(),
	}

	feeArgs, err := t.getFeeArgs(ctx, queueTypeName, msg.SenderAddress,
		msg.Fees, origMessage, opts)
	if err != nil {
		return nil, 0, err
	}

	// TODO: Use generated contract code directly
	// compass 2.0.0
	// def submit_logic_call(consensus: Consensus, args: LogicCallArgs, fee_args: FeeArgs, message_id: uint256, deadline: uint256, relayer: address)
	args := []any{
		con,
		compassArgs,
		feeArgs,
		new(big.Int).SetInt64(int64(origMessage.ID)),
		new(big.Int).SetInt64(msg.GetDeadline()),
		ethSender,
	}

	if msg.ExecutionRequirements.EnforceMEVRelay {
		opts.useMevRelay = true
	}
	logger.WithField("consensus", con).WithField("args", args).Debug("submitting logic call")
	tx, err := t.callCompass(ctx, opts, "submit_logic_call", args)
	if err != nil {
		return nil, 0, err
	}

	return tx, valsetID, nil
}

func (t compass) compass_handover(
	ctx context.Context,
	queueTypeName string,
	msg *evmtypes.CompassHandover,
	origMessage chain.MessageWithSignatures,
	ethSender common.Address,
	estimate *big.Int,
	opts callOptions,
) (*ethtypes.Transaction, uint64, error) {
	logger := liblog.WithContext(ctx).WithFields(log.Fields{
		"chain-id": t.ChainReferenceID,
		"msg-id":   origMessage.ID,
	})
	logger.Info("compass handover")

	valsetID, err := t.performValsetIDCrosscheck(ctx, t.ChainReferenceID)
	logger = logger.WithField("last-valset-id", valsetID)
	if err != nil {
		if errors.Is(err, errValsetIDMismatch) {
			logger.Warn("Valset ID mismatch. Swallowing error to retry message...")
			return nil, 0, nil
		}

		return nil, 0, err
	}

	valset, err := t.paloma.QueryGetEVMValsetByID(ctx, valsetID, t.ChainReferenceID)
	if err != nil {
		return nil, 0, err
	}

	consensusReached := isConsensusReached(ctx, valset, origMessage)
	if !consensusReached {
		return nil, 0, ErrNoConsensus
	}

	con := BuildCompassConsensus(valset, origMessage.Signatures)
	compassArgs := slice.Map(msg.GetForwardCallArgs(), func(arg evmtypes.CompassHandover_ForwardCallArgs) CompassLogicCallArgs {
		return CompassLogicCallArgs{arg.GetPayload(), common.HexToAddress(arg.GetHexContractAddress())}
	})

	if opts.estimateOnly {
		// Simulate maximum gas estimate to ensure the transaction is not rejected
		estimate = big.NewInt(0).SetUint64(cConservativeDummyGasEstimate)
	}

	// TODO: Use generated contract code directly
	// compass 2.0.0
	// def compass_update_batch(consensus: Consensus, update_compass_args: DynArray[LogicCallArgs, MAX_BATCH], deadline: uint256, gas_estimate: uint256, relayer: address):
	args := []any{
		con,
		compassArgs,
		new(big.Int).SetInt64(msg.GetDeadline()),
		estimate,
		ethSender,
	}

	logger.WithField("consensus", con).WithField("args", args).Debug("compass handover")
	tx, err := t.callCompass(ctx, opts, "compass_update_batch", args)
	if err != nil {
		return nil, 0, err
	}

	return tx, valsetID, nil
}

func (t compass) uploadSmartContract(
	ctx context.Context,
	queueTypeName string,
	msg *evmtypes.UploadSmartContract,
	origMessage chain.MessageWithSignatures,
) (*ethtypes.Transaction, error) {
	constructorInput := msg.GetConstructorInput()
	logger := liblog.WithContext(ctx).WithFields(log.Fields{
		"chain-id":          t.ChainReferenceID,
		"constructor-input": constructorInput,
	})
	logger.Info("upload smart contract")

	// 0 means to get the latest valset
	latestValset, err := t.paloma.QueryGetEVMValsetByID(ctx, 0, t.ChainReferenceID)
	if err != nil {
		logger.WithError(err).Error("uploadSmartContract: error querying valset from Paloma")
		return nil, err
	}

	consensusReached := isConsensusReached(ctx, latestValset, origMessage)
	if !consensusReached {
		return nil, ErrNoConsensus
	}

	_, tx, err := t.evm.DeployContract(
		ctx,
		t.chainID,
		msg.GetAbi(),
		msg.GetBytecode(),
		constructorInput,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (t compass) uploadUserSmartContract(
	ctx context.Context,
	queueTypeName string,
	msg *evmtypes.UploadUserSmartContract,
	origMessage chain.MessageWithSignatures,
	ethSender common.Address,
	opts callOptions,
) (*ethtypes.Transaction, uint64, error) {
	logger := liblog.WithContext(ctx).WithFields(log.Fields{
		"chain-id": t.ChainReferenceID,
		"msg-id":   origMessage.ID,
	})
	logger.Info("upload user smart contract")

	// Skip already executed check in case of estimate only
	if !opts.estimateOnly {
		executed, err := t.isUserSmartContractUploaded(ctx, origMessage.ID)
		if err != nil {
			return nil, 0, err
		}
		if executed {
			return nil, 0, ErrCallAlreadyExecuted
		}
	}

	valsetID, err := t.performValsetIDCrosscheck(ctx, t.ChainReferenceID)
	logger = logger.WithField("last-valset-id", valsetID)
	if err != nil {
		if errors.Is(err, errValsetIDMismatch) {
			logger.Warn("Valset ID mismatch. Swallowing error to retry message...")
			return nil, 0, nil
		}

		return nil, 0, err
	}

	valset, err := t.paloma.QueryGetEVMValsetByID(ctx, valsetID, t.ChainReferenceID)
	if err != nil {
		return nil, 0, err
	}

	consensusReached := isConsensusReached(ctx, valset, origMessage)
	if !consensusReached {
		return nil, 0, ErrNoConsensus
	}

	con := BuildCompassConsensus(valset, origMessage.Signatures)

	feeArgs, err := t.getFeeArgs(ctx, queueTypeName, msg.SenderAddress,
		msg.Fees, origMessage, opts)
	if err != nil {
		return nil, 0, err
	}

	// TODO: Use generated contract code directly
	// compass 2.0.0
	// def deploy_contract(consensus: Consensus, _deployer: address, _bytecode: Bytes[24576], fee_args: FeeArgs, message_id: uint256, deadline: uint256, relayer: address)
	args := []any{
		con,
		common.HexToAddress(msg.GetDeployerAddress()),
		msg.GetBytecode(),
		feeArgs,
		new(big.Int).SetInt64(int64(origMessage.ID)),
		new(big.Int).SetInt64(msg.GetDeadline()),
		ethSender,
	}

	logger.WithField("consensus", con).WithField("args", args).
		Debug("deploying user smart contract")
	tx, err := t.callCompass(ctx, opts, "deploy_contract", args)
	if err != nil {
		return nil, 0, err
	}

	return tx, valsetID, nil
}

func (t compass) SetErrorData(
	ctx context.Context,
	queueTypeName string,
	msgID uint64,
	errToProcess error,
) error {
	data := []byte(errToProcess.Error())

	var jsonRpcErr rpc.DataError
	if errors.As(errToProcess, &jsonRpcErr) {
		liblog.WithContext(ctx).WithFields(
			log.Fields{
				"queue-type-name": queueTypeName,
				"message-id":      msgID,
				"error-message":   jsonRpcErr.Error(),
			},
		).Warn("smart contract returned an error")

		data = []byte(jsonRpcErr.Error())
	}

	return t.paloma.SetErrorData(ctx, queueTypeName, msgID, data)
}

func (t compass) findLastValsetMessageID(ctx context.Context) (uint64, error) {
	logger := liblog.WithContext(ctx)
	logger.Debug("fetching last valset message id")
	id, err := t.evm.LastValsetID(ctx, t.smartContractAddr)
	if err != nil {
		logger.WithError(err).WithField("addr", t.smartContractAddr.String()).Error("error getting LastValsetID")
		return 0, fmt.Errorf("error getting LastValsetID: %w", err)
	}

	logger.WithField("last-valset-id", id).Debug("fetching last valset message id: done!")
	return id.Uint64(), nil
}

func (t compass) isArbitraryCallAlreadyExecuted(ctx context.Context, messageID uint64) (bool, error) {
	topics := [][]common.Hash{
		{
			crypto.Keccak256Hash([]byte("LogicCallEvent(address,bytes,uint256,uint256)")),
			common.Hash{},
			common.Hash{},
			crypto.Keccak256Hash(new(big.Int).SetInt64(int64(messageID)).Bytes()),
		},
	}
	filter, err := ethfilter.Factory().
		WithFromBlockNumberProvider(t.evm.FindCurrentBlockNumber).
		WithFromBlockNumberSafetyMargin(9999).
		WithTopics(topics...).
		WithAddresses(t.smartContractAddr).
		Filter(ctx)
	if err != nil {
		return false, err
	}

	var found bool
	_, err = t.evm.FilterLogs(ctx, filter, nil, func(logs []ethtypes.Log) bool {
		for _, ethLog := range logs {
			event, err := t.compassAbi.Unpack("LogicCallEvent", ethLog.Data)
			if err != nil {
				found = true
				return found
			}

			logMessageID, ok := event[2].(*big.Int)
			if !ok {
				found = true
			}
			found = messageID == logMessageID.Uint64()
			if found {
				return found
			}
		}

		return found
	})
	if err != nil {
		return false, err
	}

	return found, nil
}

func (t compass) isUserSmartContractUploaded(ctx context.Context, messageID uint64) (bool, error) {
	topics := [][]common.Hash{
		{
			crypto.Keccak256Hash([]byte("ContractDeployed(address,address,uint256)")),
			common.Hash{},
			common.Hash{},
			crypto.Keccak256Hash(new(big.Int).SetInt64(int64(messageID)).Bytes()),
		},
	}
	filter, err := ethfilter.Factory().
		WithFromBlockNumberProvider(t.evm.FindCurrentBlockNumber).
		WithFromBlockNumberSafetyMargin(9999).
		WithTopics(topics...).
		WithAddresses(t.smartContractAddr).
		Filter(ctx)
	if err != nil {
		return false, err
	}

	var found bool
	_, err = t.evm.FilterLogs(ctx, filter, nil, func(logs []ethtypes.Log) bool {
		for _, ethLog := range logs {
			event, err := t.compassAbi.Unpack("ContractDeployed", ethLog.Data)
			if err != nil {
				found = true
				return found
			}

			logMessageID, ok := event[2].(*big.Int)
			if !ok {
				found = true
			}
			found = messageID == logMessageID.Uint64()
			if found {
				return found
			}
		}

		return found
	})
	if err != nil {
		return false, err
	}

	return found, nil
}

func (t compass) skywayIsBatchAlreadyRelayed(ctx context.Context, batchNonce uint64) (bool, error) {
	filter, err := ethfilter.Factory().
		WithFromBlockNumberProvider(t.evm.FindCurrentBlockNumber).
		WithFromBlockNumberSafetyMargin(9999).
		WithTopics([]common.Hash{crypto.Keccak256Hash([]byte("BatchSendEvent(address,uint256,uint256)"))}).
		WithAddresses(t.smartContractAddr).
		Filter(ctx)
	if err != nil {
		return false, err
	}

	found, err := t.evm.FilterLogs(ctx, filter, nil, func(logs []ethtypes.Log) bool {
		for _, ethLog := range logs {
			event, err := t.compassAbi.Unpack("BatchSendEvent", ethLog.Data)
			if err != nil {
				// Failed to unpack. Let's assume the message has been sent already.
				liblog.WithContext(ctx).WithField("batch-nonce", batchNonce).WithField("event", event).WithError(err).Error("Failed to unpack log event.")
				return true
			}

			logBatchNonce, ok := event[1].(*big.Int)
			if !ok {
				// Failed to parse nonce. Let's assume the message has been sent already.
				liblog.WithContext(ctx).WithField("batch-nonce", batchNonce).WithField("log-batch-nonce", logBatchNonce).Error("Failed to parse nonce to *big.Int.")
				return true
			}

			if batchNonce == logBatchNonce.Uint64() {
				return true
			}
		}

		return false
	})
	if err != nil {
		return false, err
	}

	return found, nil
}

func BuildCompassConsensus(
	v *evmtypes.Valset,
	signatures []chain.ValidatorSignature,
) CompassConsensus {
	signatureMap := slice.MakeMapKeys(
		signatures,
		func(sig chain.ValidatorSignature) string {
			return sig.SignedByAddress
		},
	)
	con := CompassConsensus{
		Valset: TransformValsetToCompassValset(v),
	}

	for i := range v.GetValidators() {
		sig, ok := signatureMap[v.GetValidators()[i]]
		if !ok {
			con.Signatures = append(con.Signatures,
				Signature{
					V: big.NewInt(0),
					R: big.NewInt(0),
					S: big.NewInt(0),
				})
		} else {
			con.Signatures = append(con.Signatures,
				Signature{
					V: new(big.Int).SetInt64(int64(sig.Signature[64]) + 27),
					R: new(big.Int).SetBytes(sig.Signature[:32]),
					S: new(big.Int).SetBytes(sig.Signature[32:64]),
				},
			)
		}
		con.originalSignatures = append(con.originalSignatures, sig.Signature)
	}

	return con
}

func (t compass) processMessages(ctx context.Context, queueTypeName string, msgs []chain.MessageWithSignatures, opts callOptions) ([]*ethtypes.Transaction, error) {
	var gErr whoops.Group
	logger := liblog.WithContext(ctx).WithField("queue-type-name", queueTypeName)
	res := make([]*ethtypes.Transaction, 0, len(msgs))
	for i, rawMsg := range msgs {
		logger = logger.WithField("message-id", rawMsg.ID)

		if ctx.Err() != nil {
			logger.Debug("exiting processing message context")
			break
		}

		var processingErr error
		var tx *ethtypes.Transaction
		var valsetID uint64
		msg := rawMsg.Msg.(*evmtypes.Message)
		ethSender, err := func() (common.Address, error) {
			// Do not retrieve eth sender for UploadSmartContract messages
			switch msg.GetAction().(type) {
			case *evmtypes.Message_UploadSmartContract:
				return common.Address{}, nil
			}
			return t.findAssigneeEthAddress(ctx, msg.Assignee)
		}()
		if err != nil {
			return res, fmt.Errorf("failed to find assignee eth address: %w", err)
		}

		logger := logger.WithFields(log.Fields{
			"message-type":           fmt.Sprintf("%T", msg.GetAction()),
			"chain-reference-id":     t.ChainReferenceID,
			"queue-name":             queueTypeName,
			"msg-id":                 rawMsg.ID,
			"msg-bytes-to-sign":      rawMsg.BytesToSign,
			"msg-msg":                rawMsg.Msg,
			"msg-nonce":              rawMsg.Nonce,
			"msg-public-access-data": rawMsg.PublicAccessData,
			"msg-eth-sender":         ethSender.String(),
		})
		logger.Debug("processing")

		switch action := msg.GetAction().(type) {
		case *evmtypes.Message_SubmitLogicCall:
			logger := logger.WithFields(log.Fields{
				"message-type": "Message_SubmitLogicCall",
			})
			logger.Debug("switch-case-message-update-valset")

			tx, valsetID, processingErr = t.submitLogicCall(
				ctx,
				queueTypeName,
				action.SubmitLogicCall,
				rawMsg,
				ethSender,
				opts,
			)
		case *evmtypes.Message_UpdateValset:
			logger := logger.WithFields(log.Fields{
				"message-type": "Message_UpdateValset",
			})
			logger.Debug("switch-case-message-update-valset")
			tx, valsetID, processingErr = t.updateValset(
				ctx,
				queueTypeName,
				action.UpdateValset.Valset,
				rawMsg,
				ethSender,
				rawMsg.Estimate,
				opts,
			)
		case *evmtypes.Message_UploadSmartContract:
			logger := logger.WithFields(log.Fields{
				"message-type": "Message_UploadSmartContract",
			})
			logger.Debug("switch-case-message-upload-contract")
			tx, processingErr = t.uploadSmartContract(
				ctx,
				queueTypeName,
				action.UploadSmartContract,
				rawMsg,
			)
		case *evmtypes.Message_UploadUserSmartContract:
			logger := logger.WithFields(log.Fields{
				"message-type": "Message_UploadUserSmartContract",
			})
			logger.Debug("switch-case-message-upload-user-contract")

			tx, valsetID, processingErr = t.uploadUserSmartContract(
				ctx,
				queueTypeName,
				action.UploadUserSmartContract,
				rawMsg,
				ethSender,
				opts,
			)
		case *evmtypes.Message_CompassHandover:
			logger := logger.WithFields(log.Fields{
				"message-type": "Message_CompassHandover",
			})
			logger.Debug("switch-case-message-upload-user-contract")

			tx, valsetID, processingErr = t.compass_handover(
				ctx,
				queueTypeName,
				action.CompassHandover,
				rawMsg,
				ethSender,
				rawMsg.Estimate,
				opts,
			)
		default:
			return res, ErrUnsupportedMessageType.Format(action)
		}

		processingErr = whoops.Enrich(
			processingErr,
			FieldMessageID.Val(rawMsg.ID),
			FieldMessageType.Val(msg.GetAction()),
		)

		// Append all txs, even if they are nil
		// These values will have to be filtered out by the caller
		res = append(res, tx)

		switch {
		case processingErr == nil:
			if tx != nil && opts.estimateOnly == false {
				logger.WithFields(log.Fields{
					"msg-public-access-data": tx.Hash().Hex(),
				}).Debug("setting public access data")
				err := t.paloma.SetPublicAccessData(ctx, queueTypeName,
					rawMsg.ID, valsetID, tx.Hash().Bytes())
				if err != nil {
					gErr.Add(err)
					return res, gErr
				}
			}
		case errors.Is(processingErr, ErrNoConsensus):
			// Only log
			logger.Warn(ErrNoConsensus.Error())
			if err := t.paloma.NewStatus().
				WithChainReferenceID(t.ChainReferenceID).
				WithQueueType(queueTypeName).
				WithMsg(&msgs[i]).
				WithLog(ErrNoConsensus.Error()).
				Error(ctx); err != nil {
				logger.WithError(err).Error("failed to send paloma status update")
			}
		default:
			logger.WithError(processingErr).Error("processing error")

			if !opts.estimateOnly {
				// If we're not just estimating, we want to set the error data
				// on the message
				setErr := t.SetErrorData(ctx, queueTypeName, rawMsg.ID, processingErr)
				if setErr != nil {
					// If we got an error setting the error data, this is the error
					// we want to log
					processingErr = setErr
				}
			}

			if err := t.paloma.NewStatus().
				WithChainReferenceID(t.ChainReferenceID).
				WithQueueType(queueTypeName).
				WithMsg(&msgs[i]).
				WithLog(processingErr.Error()).
				Error(ctx); err != nil {
				logger.WithError(err).Error("failed to send paloma status update")
			}
			gErr.Add(processingErr)
		}
	}

	return res, gErr.Return()
}

func (t compass) provideEvidenceForValidatorBalance(ctx context.Context, queueTypeName string, msgs []chain.MessageWithSignatures) error {
	logger := liblog.WithContext(ctx).WithField("queue-type-name", queueTypeName)
	logger.Debug("start processing validator balance request")

	for _, msg := range msgs {
		vb := msg.Msg.(*evmtypes.ValidatorBalancesAttestation)
		height, err := t.evm.FindBlockNearestToTime(ctx, uint64(t.startingBlockHeight), vb.FromBlockTime)
		if err != nil {
			logger.WithError(err).Error("failed to find block nearest to time")
			return err
		}

		logger := logger.WithFields(
			log.Fields{
				"height":          height,
				"nearest-to-time": vb.FromBlockTime,
			},
		)
		logger.Debug("got height for time")

		res := &evmtypes.ValidatorBalancesAttestationRes{
			BlockHeight: height,
			Balances:    make([]string, len(vb.HexAddresses)),
		}

		success := false
		for i, addrHex := range vb.HexAddresses {
			addr := common.HexToAddress(addrHex)
			balance, err := t.evm.BalanceAt(ctx, addr, height)
			if err != nil {
				logger.WithError(err).WithFields(log.Fields{
					"addr": addr,
				}).Warn("failed to get balance")
				continue
			}

			success = true

			logger.WithFields(log.Fields{
				"evm-address": addr,
				"balance":     balance,
			}).Info("got balance")

			res.Balances[i] = balance.Text(10)
		}

		if !success {
			// If all requests fail, there may be something wrong with the RPC,
			// so we return an error to avoid all validators from being jailed
			err := errors.New("all balance requests failed")
			logger.WithError(err).Warn("failed to get balances")
			return err
		}

		err = t.paloma.AddMessageEvidence(ctx, queueTypeName, msg.ID, res)
		if err != nil {
			logger.WithError(err).Error("failed to add message evidence")
			return err
		}
	}

	return nil
}

func (t compass) provideEvidenceForReferenceBlock(ctx context.Context, queueTypeName string, msgs []chain.MessageWithSignatures) error {
	logger := liblog.WithContext(ctx).WithField("queue-type-name", queueTypeName)
	logger.Debug("start processing reference block request")

	for _, msg := range msgs {
		vb := msg.Msg.(*evmtypes.ReferenceBlockAttestation)
		height, err := t.evm.FindBlockNearestToTime(ctx, uint64(t.startingBlockHeight), vb.FromBlockTime)
		if err != nil {
			return err
		}

		logger := logger.WithFields(
			log.Fields{
				"height":          height,
				"nearest-to-time": vb.FromBlockTime,
			},
		)
		logger.Debug("got reference block height for time")

		h, err := t.evm.GetEthClient().HeaderByNumber(ctx, new(big.Int).SetUint64(height))
		if err != nil {
			return err
		}

		res := &evmtypes.ReferenceBlockAttestationRes{
			BlockHeight: height,
			BlockHash:   h.Hash().String(),
		}

		logger.WithFields(
			log.Fields{
				"hash": h.Hash().String(),
			},
		).Debug("got reference block hash")

		err = t.paloma.AddMessageEvidence(ctx, queueTypeName, msg.ID, res)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *compass) getLogs(
	ctx context.Context,
	logger *logrus.Entry,
	from, to uint64,
) ([]ethtypes.Log, uint64, error) {
	var filter ethereum.FilterQuery
	var err error

	if from == 0 && to == 0 {
		// If from and to are zero, we default to search the latest blocks
		filter, err = ethfilter.Factory().
			WithFromBlockNumberProvider(t.evm.FindCurrentBlockNumber).
			WithFromBlockNumberSafetyMargin(1).
			WithTopics([]common.Hash{
				batchSendEvent,
				sendToPalomaEvent,
				nodeSaleEvent,
			}).
			WithAddresses(t.smartContractAddr).
			Filter(ctx)
		if err != nil {
			return nil, 0, err
		}

		currentBlockNumber := filter.FromBlock.Uint64()
		if t.lastObservedBlockHeight == 0 {
			t.lastObservedBlockHeight = currentBlockNumber - 10_000
		}

		if currentBlockNumber < t.lastObservedBlockHeight {
			// Either we messed up tracking the current block, or the RPC is
			// having some issue. Either way, it's best to sync again.
			t.lastObservedBlockHeight = currentBlockNumber
		}

		filter.FromBlock = big.NewInt(0).SetUint64(t.lastObservedBlockHeight)
		filter.ToBlock = big.NewInt(0).SetUint64(min(t.lastObservedBlockHeight+10_000, currentBlockNumber))
	} else {
		filter = ethereum.FilterQuery{
			Addresses: []common.Address{t.smartContractAddr},
			Topics: [][]common.Hash{{
				batchSendEvent,
				sendToPalomaEvent,
				nodeSaleEvent,
			}},
			FromBlock: big.NewInt(0).SetUint64(from),
			ToBlock:   big.NewInt(0).SetUint64(to),
		}
	}

	logger.WithField("from", filter.FromBlock).
		WithField("to", filter.ToBlock).
		Debug("Filter is ready")

	logs, err := t.evm.GetEthClient().FilterLogs(ctx, filter)
	if err != nil {
		logger.WithField("from", filter.FromBlock).
			WithField("to", filter.ToBlock).
			WithField("req_from", from).
			WithField("req_to", to).
			WithField("chain_reference_id", t.ChainReferenceID).
			WithError(err).
			Warn("Failed to filter events")
	}

	return logs, filter.ToBlock.Uint64(), err
}

func (t *compass) GetSkywayEvents(
	ctx context.Context,
	orchestrator string,
) ([]chain.SkywayEventer, error) {
	logger := liblog.WithContext(ctx)

	logger.Debug("Querying compass events")

	blocks, err := t.paloma.QueryUnobservedBlocksByValidator(ctx, t.ChainReferenceID, orchestrator)
	if err != nil {
		logger.WithError(err).Warn("Failed to query unobserved blocks")
		return nil, err
	}

	var logs []ethtypes.Log
	var toBlock uint64

	if len(blocks) == 0 {
		logs, toBlock, err = t.getLogs(ctx, logger, 0, 0)
		if err != nil {
			return nil, err
		}
	} else {
		var moreLogs []ethtypes.Log
		for i := range blocks {
			logger.WithField("block", blocks[i]).
				Debug("Getting logs from block")

			moreLogs, toBlock, err = t.getLogs(ctx, logger, blocks[i], blocks[i])
			if err != nil {
				return nil, err
			}

			logs = slices.Concat(logs, moreLogs)
		}
	}

	lastSkywayNonce, err := t.paloma.QueryLastObservedSkywayNonceByAddr(ctx, t.ChainReferenceID, orchestrator)
	if err != nil {
		logger.WithError(err).Warn("Failed to query last observed nonce")
		return nil, err
	}

	logger.WithField("skwyway_nonce", lastSkywayNonce).
		WithField("logs", len(logs)).
		Debug("Ready to parse events")

	var events []chain.SkywayEventer
	var evt chain.SkywayEventer

	for _, ethLog := range logs {
		switch ethLog.Topics[0] {
		case nodeSaleEvent:
			logger.Info("Parsing light node sale event")
			evt, err = t.parseLightNodeSaleEvent(ethLog.Data, ethLog.BlockNumber)
		case batchSendEvent:
			logger.Info("Parsing batch send event")
			evt, err = t.parseBatchSendEvent(ethLog.Data, ethLog.BlockNumber)
		case sendToPalomaEvent:
			logger.Info("Parsing send to paloma event")
			evt, err = t.parseSendToPalomaEvent(ethLog.Data, ethLog.BlockNumber)
		default:
			logger.WithField("event", ethLog).Warn("Unknown event from compass")
			continue
		}

		if err != nil {
			logger.WithError(err).Warn("Error parsing event")
			return nil, err
		}

		if evt.GetSkywayNonce() <= lastSkywayNonce {
			logger.WithField("last-event-nonce", lastSkywayNonce).
				WithField("skyway-nonce", evt.GetSkywayNonce()).
				Info("Skipping already observed event...")
			continue
		}

		events = append(events, evt)
	}

	t.lastObservedBlockHeight = toBlock

	return events, err
}

func (t *compass) parseSendToPalomaEvent(
	data []byte,
	blockHeight uint64,
) (evt chain.SendToPalomaEvent, err error) {
	event, err := t.compassAbi.Unpack("SendToPalomaEvent", data)
	if err != nil {
		return evt, err
	}

	tokenContract, ok := event[0].(common.Address)
	if !ok {
		return evt, fmt.Errorf("invalid token contract")
	}

	ethSender, ok := event[1].(common.Address)
	if !ok {
		return evt, fmt.Errorf("invalid sender address")
	}

	palomaReceiver, err := compassBytesToPalomaAddress(event[2])
	if err != nil {
		return evt, fmt.Errorf("invalid paloma receiver: %w", err)
	}

	amount, ok := event[3].(*big.Int)
	if !ok {
		return evt, fmt.Errorf("invalid amount")
	}

	skywayNonce, ok := event[4].(*big.Int)
	if !ok {
		return evt, fmt.Errorf("invalid paloma nonce")
	}

	eventNonce, ok := event[5].(*big.Int)
	if !ok {
		return evt, fmt.Errorf("invalid event nonce")
	}

	return chain.SendToPalomaEvent{
		EthBlockHeight: blockHeight,
		EventNonce:     eventNonce.Uint64(),
		Amount:         amount.Uint64(),
		EthereumSender: ethSender.String(),
		PalomaReceiver: palomaReceiver.String(),
		TokenContract:  tokenContract.String(),
		SkywayNonce:    skywayNonce.Uint64(),
	}, nil
}

func (t *compass) parseBatchSendEvent(
	data []byte,
	blockHeight uint64,
) (evt chain.BatchSendEvent, err error) {
	event, err := t.compassAbi.Unpack("BatchSendEvent", data)
	if err != nil {
		return evt, err
	}

	tokenContract, ok := event[0].(common.Address)
	if !ok {
		return evt, fmt.Errorf("invalid token contract")
	}

	batchNonce, ok := event[1].(*big.Int)
	if !ok {
		return evt, fmt.Errorf("invalid batch nonce")
	}

	skywayNonce, ok := event[2].(*big.Int)
	if !ok {
		return evt, fmt.Errorf("invalid skyway nonce")
	}

	eventNonce, ok := event[3].(*big.Int)
	if !ok {
		return evt, fmt.Errorf("invalid event nonce")
	}

	return chain.BatchSendEvent{
		EthBlockHeight: blockHeight,
		EventNonce:     eventNonce.Uint64(),
		BatchNonce:     batchNonce.Uint64(),
		TokenContract:  tokenContract.String(),
		SkywayNonce:    skywayNonce.Uint64(),
	}, nil
}

func (t *compass) parseLightNodeSaleEvent(
	data []byte,
	blockHeight uint64,
) (evt chain.LightNodeSaleEvent, err error) {
	event, err := t.compassAbi.Unpack("NodeSaleEvent", data)
	if err != nil {
		return evt, err
	}

	contractAddress, ok := event[0].(common.Address)
	if !ok {
		return evt, fmt.Errorf("invalid smart contract address")
	}

	clientAddress, err := compassBytesToPalomaAddress(event[2])
	if err != nil {
		return evt, fmt.Errorf("invalid client address: %w", err)
	}

	amount, ok := event[4].(*big.Int)
	if !ok {
		return evt, fmt.Errorf("invalid amount")
	}

	skywayNonce, ok := event[5].(*big.Int)
	if !ok {
		return evt, fmt.Errorf("invalid paloma nonce")
	}

	eventNonce, ok := event[6].(*big.Int)
	if !ok {
		return evt, fmt.Errorf("invalid event nonce")
	}

	return chain.LightNodeSaleEvent{
		EthBlockHeight:       blockHeight,
		EventNonce:           eventNonce.Uint64(),
		Amount:               amount.Uint64(),
		ClientAddress:        clientAddress.String(),
		SkywayNonce:          skywayNonce.Uint64(),
		SmartContractAddress: contractAddress.String(),
	}, nil
}

// provideTxProof provides a very simple proof which is a transaction object
func (t compass) provideTxProof(ctx context.Context, queueTypeName string, rawMsg chain.MessageWithSignatures) error {
	liblog.WithContext(ctx).WithFields(log.Fields{
		"queue-type-name":    queueTypeName,
		"msg-id":             rawMsg.ID,
		"public-access-data": rawMsg.PublicAccessData,
	}).Debug("providing proof")
	txHash := common.BytesToHash(rawMsg.PublicAccessData)
	tx, _, err := t.evm.TransactionByHash(ctx, txHash)
	if err != nil {
		return err
	}

	txProof, err := tx.MarshalBinary()
	if err != nil {
		return err
	}

	var serializedReceipt []byte

	// If this is a turnstone message, we need to get the tx receipt
	if _, ok := rawMsg.Msg.(*evmtypes.Message); ok {
		receipt, err := t.evm.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			return err
		}

		serializedReceipt, err = receipt.MarshalBinary()
		if err != nil {
			return err
		}
	}

	return t.paloma.AddMessageEvidence(ctx, queueTypeName, rawMsg.ID, &evmtypes.TxExecutedProof{
		SerializedTX:      txProof,
		SerializedReceipt: serializedReceipt,
	})
}

func (t compass) submitBatchSendToEVMClaim(ctx context.Context, event chain.BatchSendEvent, orchestrator string) error {
	msg := skywaytypes.MsgBatchSendToRemoteClaim{
		EventNonce:       event.EventNonce,
		EthBlockHeight:   event.EthBlockHeight,
		BatchNonce:       event.BatchNonce,
		TokenContract:    event.TokenContract,
		ChainReferenceId: t.ChainReferenceID,
		Orchestrator:     orchestrator,
		SkywayNonce:      event.SkywayNonce,
		CompassId:        t.CompassID,
	}
	return t.paloma.SendBatchSendToEVMClaim(ctx, msg)
}

func (t compass) submitSendToPalomaClaim(ctx context.Context, event chain.SendToPalomaEvent, orchestrator string) error {
	msg := skywaytypes.MsgSendToPalomaClaim{
		EventNonce:       event.EventNonce,
		EthBlockHeight:   event.EthBlockHeight,
		TokenContract:    event.TokenContract,
		Amount:           math.NewInt(int64(event.Amount)),
		EthereumSender:   event.EthereumSender,
		PalomaReceiver:   event.PalomaReceiver,
		ChainReferenceId: t.ChainReferenceID,
		Orchestrator:     orchestrator,
		SkywayNonce:      event.SkywayNonce,
		CompassId:        t.CompassID,
	}
	return t.paloma.SendSendToPalomaClaim(ctx, msg)
}

func (t compass) submitLightNodeSaleClaim(ctx context.Context, event chain.LightNodeSaleEvent, orchestrator string) error {
	msg := skywaytypes.MsgLightNodeSaleClaim{
		EventNonce:           event.EventNonce,
		EthBlockHeight:       event.EthBlockHeight,
		ClientAddress:        event.ClientAddress,
		Amount:               math.NewInt(int64(event.Amount)),
		ChainReferenceId:     t.ChainReferenceID,
		Orchestrator:         orchestrator,
		SkywayNonce:          event.SkywayNonce,
		SmartContractAddress: event.SmartContractAddress,
		CompassId:            t.CompassID,
	}
	return t.paloma.SendLightNodeSaleClaim(ctx, msg)
}

// provideErrorProof provides a pass-through proof for an error during relaying
func (t compass) provideErrorProof(ctx context.Context, queueTypeName string, rawMsg chain.MessageWithSignatures) error {
	liblog.WithContext(ctx).WithFields(log.Fields{
		"queue-type-name":    queueTypeName,
		"msg-id":             rawMsg.ID,
		"public-access-data": rawMsg.PublicAccessData,
	}).Debug("providing error proof")

	return t.paloma.AddMessageEvidence(ctx, queueTypeName, rawMsg.ID, &evmtypes.SmartContractExecutionErrorProof{
		ErrorMessage: string(rawMsg.PublicAccessData),
	})
}

func TransformValsetToCompassValset(val *evmtypes.Valset) CompassValset {
	return CompassValset{
		Validators: slice.Map(val.GetValidators(), func(s string) common.Address {
			return common.HexToAddress(s)
		}),
		Powers: slice.Map(val.GetPowers(), func(p uint64) *big.Int {
			return big.NewInt(int64(p))
		}),
		ValsetId: big.NewInt(int64(val.GetValsetID())),
	}
}

func isConsensusReached(ctx context.Context, val *evmtypes.Valset, msg chain.SignedEntity) bool {
	signaturesMap := make(map[string]chain.ValidatorSignature)
	for _, sig := range msg.GetSignatures() {
		signaturesMap[sig.SignedByAddress] = sig
	}
	logger := liblog.WithContext(ctx).WithFields(
		log.Fields{
			"validators-size": len(val.Validators),
		})
	logger.Debug("confirming consensus reached")
	var totalPower uint64
	for _, pow := range val.Powers {
		totalPower += pow
	}
	powerThreshold := totalPower * 2 / 3

	var s uint64
	for i := range val.Validators {
		valHex, pow := val.Validators[i], val.Powers[i]
		sig, ok := signaturesMap[valHex]
		logger.WithFields(log.Fields{
			"i":         i,
			"validator": valHex,
			"power":     pow,
		}).Debug("checking consensus")
		if !ok {
			continue
		}
		bytesToVerify := crypto.Keccak256(append(
			[]byte(SignedMessagePrefix),
			msg.GetBytes()...,
		))
		recoveredPK, err := crypto.Ecrecover(bytesToVerify, sig.Signature)
		if err != nil {
			continue
		}
		logger.WithFields(log.Fields{
			"i": i,
		}).Debug("good ecrecover")
		pk, err := crypto.UnmarshalPubkey(recoveredPK)
		if err != nil {
			continue
		}
		logger.WithFields(log.Fields{
			"i": i,
		}).Debug("good unmarshal")
		recoveredAddr := crypto.PubkeyToAddress(*pk)
		recoveredAddrHex := recoveredAddr.Hex()
		if valHex != recoveredAddrHex {
			continue
		}
		s += pow
		logger.WithFields(log.Fields{
			"i": i,
		}).Debug("good consensus")
	}
	return s >= powerThreshold
}

type callOptions struct {
	useMevRelay  bool
	estimateOnly bool
}

func (c compass) callCompass(
	ctx context.Context,
	opts callOptions,
	method string,
	arguments []any,
) (*ethtypes.Transaction, error) {
	if c.compassAbi == nil {
		return nil, ErrABINotInitialized
	}
	return c.evm.ExecuteSmartContract(ctx, c.chainID, *c.compassAbi, c.smartContractAddr, opts, method, arguments)
}

func (t compass) skywayRelayBatches(ctx context.Context, batches []chain.SkywayBatchWithSignatures) error {
	var gErr whoops.Group
	logger := liblog.WithContext(ctx).WithField("chainReferenceID", t.ChainReferenceID)
	for _, batch := range batches {
		logger = logger.WithField("batch-nonce", batch.BatchNonce)

		if ctx.Err() != nil {
			logger.Debug("exiting processing batch context")
			break
		}

		var processingErr error
		logger := logger.WithFields(log.Fields{
			"batch-nonce": batch.BatchNonce,
		})
		logger.Debug("relaying")

		_, processingErr = t.skywayRelayBatch(ctx, batch, callOptions{})
		processingErr = whoops.Enrich(
			processingErr,
			FieldMessageID.Val(batch.BatchNonce),
		)

		switch {
		case processingErr == nil:
			// do nothing.  claiming happens in a different goroutine
		case errors.Is(processingErr, ErrNoConsensus):
			// does nothing
		default:
			logger.WithError(processingErr).Error("relay error")
			gErr.Add(processingErr)
		}
	}

	return gErr.Return()
}

func (t compass) skywayEstimateBatches(ctx context.Context, batches []chain.SkywayBatchWithSignatures) ([]uint64, error) {
	estimates := make([]uint64, 0, len(batches))
	logger := liblog.WithContext(ctx).WithField("chainReferenceID", t.ChainReferenceID)
	for _, batch := range batches {
		logger = logger.WithField("batch-nonce", batch.BatchNonce)

		if ctx.Err() != nil {
			logger.Debug("exiting processing batch context")
			break
		}

		logger := logger.WithFields(log.Fields{
			"batch-nonce": batch.BatchNonce,
		})
		logger.Debug("estimating")

		tx, err := t.skywayRelayBatch(ctx, batch, callOptions{estimateOnly: true})
		if err != nil {
			logger.WithError(err).Error("failed to estimate batch")
			return nil, fmt.Errorf("failed to estimate batch: %w", err)
		}

		logger.WithField("estimate", tx.Gas()).Debug("Estimated gas for batch")
		estimates = append(estimates, tx.Gas())
	}

	return estimates, nil
}

func (t compass) skywayRelayBatch(
	ctx context.Context,
	batch chain.SkywayBatchWithSignatures,
	opts callOptions,
) (*ethtypes.Transaction, error) {
	return whoops.TryVal(func() *ethtypes.Transaction {
		logger := liblog.WithContext(ctx).
			WithField("component", "skyway-relay-batch").
			WithField("skyway-batch-nonce", batch.BatchNonce).
			WithField("chain-reference-id", batch.ChainReferenceId)
		executed, err := t.skywayIsBatchAlreadyRelayed(ctx, batch.BatchNonce)
		whoops.Assert(err)
		if executed {
			logger.Warn("skyway batch already executed!")
			return nil
		}

		valsetID, err := t.performValsetIDCrosscheck(ctx, batch.ChainReferenceId)
		if err != nil {
			if errors.Is(err, errValsetIDMismatch) {
				logger.Warn("Valset ID mismatch. Swallowing error to retry message...")
				return nil
			}

			whoops.Assert(err)
		}

		valset, err := t.paloma.QueryGetEVMValsetByID(ctx, valsetID, t.ChainReferenceID)
		whoops.Assert(err)

		consensusReached := isConsensusReached(ctx, valset, batch)
		if !consensusReached {
			whoops.Assert(ErrNoConsensus)
		}

		con := BuildCompassConsensus(valset, batch.Signatures)

		receivers := make([]common.Address, len(batch.Transactions))
		amounts := make([]*big.Int, len(batch.Transactions))

		for i, transaction := range batch.Transactions {
			receivers[i] = common.HexToAddress(transaction.DestAddress)
			amounts[i] = transaction.Erc20Token.Amount.BigInt()
		}

		compassArgs := CompassTokenSendArgs{
			Receiver: receivers,
			Amount:   amounts,
		}

		// get relayer
		// get gas estimate
		ethSender, err := t.findAssigneeEthAddress(ctx, batch.Assignee)
		if err != nil {
			logger.WithError(err).Error("failed to retrieve assignee eth address")
			whoops.Assert(fmt.Errorf("failed to retrieve assignee eth address: %w", err))
		}

		var estimate *big.Int = big.NewInt(0).SetUint64(cConservativeDummyGasEstimate)
		if !opts.estimateOnly {
			if batch.GasEstimate < 1 {
				logger.WithField("gas-estimate", batch.GasEstimate).Error("invalid gas estimate")
				whoops.Assert(fmt.Errorf("invalid gas estimate: %d", batch.GasEstimate))
			}
			estimate.SetUint64(batch.GasEstimate)
		}

		// TODO: Use compiled contract instead
		// compass 2.0
		// def submit_batch(consensus: Consensus, token: address, args: TokenSendArgs, batch_id: uint256, deadline: uint256, relayer: address, gas_estimate: uint256)
		tx, err := t.callCompass(
			ctx,
			opts,
			"submit_batch",
			[]any{
				con,
				common.HexToAddress(batch.TokenContract),
				compassArgs,
				new(big.Int).SetInt64(int64(batch.BatchNonce)),
				new(big.Int).SetInt64(int64(batch.GetBatchTimeout())),
				ethSender,
				estimate,
			},
		)
		if err != nil {
			logger.WithError(err).Error("failed to relay batch")
			whoops.Assert(err)
		}

		return tx
	})
}

// performValsetIDCrosscheck fetches the latest valset ID from the target chain
// as well as the expected valset ID stored on Paloma for the given target chain.
// In case of a mismatch, an error is logged to Paloma and errValsetIDMismatch is returned.
func (t compass) performValsetIDCrosscheck(ctx context.Context, chainReferenceID string) (uint64, error) {
	logger := liblog.WithContext(ctx).WithField("chain-reference-id", chainReferenceID)
	valsetID, err := t.findLastValsetMessageID(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to get last valset ID from target chain.")
		return 0, err
	}

	expectedValset, err := t.paloma.QueryGetLatestPublishedSnapshot(ctx, chainReferenceID)
	if err != nil {
		logger.WithError(err).Error("Failed to get expected valset ID from paloma.")
		return 0, err
	}

	if valsetID != expectedValset.Id {
		err := fmt.Errorf("target chain valset mismatch, expected %d, got %v", expectedValset.Id, valsetID)
		logger.WithError(err).Error("Target chain valset mismatch!")
		if err := t.paloma.NewStatus().WithChainReferenceID(t.ChainReferenceID).WithLog(err.Error()).Error(ctx); err != nil {
			logger.WithError(err).Error("Failed to send log to Paloma.")
		}
		return 0, errValsetIDMismatch
	}

	return valsetID, nil
}

func (t compass) findAssigneeEthAddress(ctx context.Context,
	palomaAddress string,
) (common.Address, error) {
	snapshot, err := t.paloma.QueryGetSnapshotByID(ctx, 0)
	if err != nil {
		return common.Address{}, err
	}

	for _, v := range snapshot.Validators {
		if v.Address.String() == palomaAddress {
			for _, ci := range v.ExternalChainInfos {
				if ci.ChainReferenceID == t.ChainReferenceID {
					return common.HexToAddress(ci.Address), nil
				}
			}
			break
		}
	}

	return common.Address{}, errors.New("assignee's eth address not found")
}

func (t compass) getFeeArgs(
	ctx context.Context,
	queueTypeName string,
	senderAddress []byte,
	fees *evmtypes.Fees,
	origMessage chain.MessageWithSignatures,
	opts callOptions,
) (FeeArgs, error) {
	padding := bytes.Repeat([]byte{0}, 32-len(senderAddress))
	paddedSenderAddress := [32]byte(append(padding, senderAddress...))

	// We use dummy fee data during estimation. This is not enough
	// to process the transaction, but the user needs to have enough
	// funds to cover the fees, even in the estimation.
	feeArgs := FeeArgs{
		RelayerFee:            big.NewInt(100_000),
		CommunityFee:          big.NewInt(100_000),
		SecurityFee:           big.NewInt(100_000),
		FeePayerPalomaAddress: paddedSenderAddress,
	}

	if opts.estimateOnly {
		return feeArgs, nil
	}

	if fees == nil {
		return feeArgs, errors.New("fees not provided")
	}

	feeArgs.RelayerFee = big.NewInt(0).SetUint64(fees.RelayerFee)
	feeArgs.CommunityFee = big.NewInt(0).SetUint64(fees.CommunityFee)
	feeArgs.SecurityFee = big.NewInt(0).SetUint64(fees.SecurityFee)

	userFunds, err := t.evm.QueryUserFunds(ctx, t.feeMgrContractAddr, paddedSenderAddress)
	if err != nil {
		return feeArgs, fmt.Errorf("failed to query user funds: %w", err)
	}

	gasPrice, err := t.evm.SuggestGasPrice(ctx)
	if err != nil {
		return feeArgs, fmt.Errorf("failed to suggest gas price: %w", err)
	}

	// (relayerFee*gasPrice + communityFee*gasPrice + security*gasPrice)
	totalFundsNeeded := big.NewInt(0).Add(
		big.NewInt(0).Mul(feeArgs.RelayerFee, gasPrice),
		big.NewInt(0).Add(
			big.NewInt(0).Mul(feeArgs.CommunityFee, gasPrice),
			big.NewInt(0).Mul(feeArgs.SecurityFee, gasPrice)))

	if userFunds.Cmp(totalFundsNeeded) < 0 {
		err := fmt.Errorf("insufficient funds for fees: %s < %s", userFunds, totalFundsNeeded)
		if sendErr := t.SetErrorData(ctx, queueTypeName, origMessage.ID, err); sendErr != nil {
			err = fmt.Errorf("failed to set error data: %w", sendErr)
		}
		return feeArgs, err
	}

	return feeArgs, nil
}

func compassBytesToPalomaAddress(b any) (sdk.AccAddress, error) {
	var addr sdk.AccAddress
	rawBytes, ok := b.([32]byte)
	if !ok {
		return addr, fmt.Errorf("invalid paloma address bytes")
	}

	var addrBytes []byte
	for i := range rawBytes {
		if rawBytes[i] != 0 {
			// Allow addresses starting with 0 and align to either 20 or 32
			// bytes
			if i < 12 {
				addrBytes = rawBytes[:]
			} else {
				addrBytes = rawBytes[12:]
			}

			break
		}
	}

	// The Unmarshal function below does not check for errors, so we need to do
	// it beforehand
	if err := sdk.VerifyAddressFormat(addrBytes); err != nil {
		return addr, err
	}

	if err := addr.Unmarshal(addrBytes); err != nil {
		return addr, err
	}

	return addr, nil
}
