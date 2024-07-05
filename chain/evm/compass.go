package evm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/math"
	"github.com/VolumeFi/whoops"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	evmtypes "github.com/palomachain/paloma/x/evm/types"
	skywaytypes "github.com/palomachain/paloma/x/skyway/types"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/internal/ethfilter"
	"github.com/palomachain/pigeon/internal/liblog"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

const (
	SignedMessagePrefix             = "\x19Ethereum Signed Message:\n32"
	cEventQueryBlockHeightMinWindow = 10
)

var errValsetIDMismatch = errors.New("valset id mismatch")

//go:generate mockery --name=evmClienter --inpackage --testonly
type evmClienter interface {
	FilterLogs(ctx context.Context, fq ethereum.FilterQuery, currBlockHeight *big.Int, fn func(logs []ethtypes.Log) bool) (bool, error)
	ExecuteSmartContract(ctx context.Context, chainID *big.Int, contractAbi abi.ABI, addr common.Address, mevRelay bool, method string, arguments []any) (*ethtypes.Transaction, error)
	DeployContract(ctx context.Context, chainID *big.Int, rawABI string, bytecode, constructorInput []byte) (contractAddr common.Address, tx *ethtypes.Transaction, err error)
	TransactionByHash(ctx context.Context, txHash common.Hash) (*ethtypes.Transaction, bool, error)

	BalanceAt(ctx context.Context, address common.Address, blockHeight uint64) (*big.Int, error)
	FindBlockNearestToTime(ctx context.Context, startingHeight uint64, when time.Time) (uint64, error)
	FindCurrentBlockNumber(ctx context.Context) (*big.Int, error)
	LastValsetID(ctx context.Context, addr common.Address) (*big.Int, error)
	GetEthClient() ethClientConn
}

type observedHeights struct {
	batchSendEvent    int64
	sendToPalomaEvent int64
}

type compass struct {
	paloma                   PalomaClienter
	evm                      evmClienter
	compassAbi               *abi.ABI
	chainID                  *big.Int
	CompassID                string
	ChainReferenceID         string
	lastObservedBlockHeights observedHeights
	startingBlockHeight      int64
	smartContractAddr        common.Address
}

func newCompassClient(
	smartContractAddrStr,
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
		CompassID:         compassID,
		ChainReferenceID:  chainReferenceID,
		smartContractAddr: common.HexToAddress(smartContractAddrStr),
		chainID:           chainID,
		compassAbi:        compassAbi,
		paloma:            paloma,
		evm:               evm,
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
) (*ethtypes.Transaction, error) {
	currentValsetID, err := t.findLastValsetMessageID(ctx)
	if err != nil {
		return nil, err
	}
	logger := liblog.WithContext(ctx).WithFields(log.Fields{
		"chain-reference-id": t.ChainReferenceID,
		"current-valset-id":  currentValsetID,
	})
	logger.Debug("update_valset")

	currentValset, err := t.paloma.QueryGetEVMValsetByID(ctx, currentValsetID, t.ChainReferenceID)
	if err != nil {
		return nil, err
	}
	if currentValset == nil {
		logger.Error("current valset is empty")
		return nil, fmt.Errorf("current valset is empty")
	}

	consensusReached := isConsensusReached(ctx, currentValset, origMessage)
	if !consensusReached {
		return nil, ErrNoConsensus
	}

	tx, err := t.callCompass(ctx, false, "update_valset", []any{
		BuildCompassConsensus(currentValset, origMessage.Signatures),
		TransformValsetToCompassValset(newValset),
	})
	if err != nil {
		logger.WithError(err).Error("call_compass error")
		isSmartContractError, setErr := t.SetErrorData(ctx, queueTypeName, origMessage.ID, err)
		if setErr != nil {
			return nil, setErr
		}
		if isSmartContractError {
			logger.Debug("smart contract error. recovering...")
			return nil, nil
		}
		whoops.Assert(err)
	}

	return tx, nil
}

func (t compass) submitLogicCall(
	ctx context.Context,
	queueTypeName string,
	msg *evmtypes.SubmitLogicCall,
	origMessage chain.MessageWithSignatures,
) (*ethtypes.Transaction, error) {
	logger := liblog.WithContext(ctx).WithFields(log.Fields{
		"chain-id": t.ChainReferenceID,
		"msg-id":   origMessage.ID,
	})
	logger.Info("submit logic call")

	executed, err := t.isArbitraryCallAlreadyExecuted(ctx, origMessage.ID)
	if err != nil {
		return nil, err
	}
	if executed {
		return nil, ErrCallAlreadyExecuted
	}

	valsetID, err := t.performValsetIDCrosscheck(ctx, t.ChainReferenceID)
	logger = logger.WithField("last-valset-id", valsetID)
	if err != nil {
		if errors.Is(err, errValsetIDMismatch) {
			logger.Warn("Valset ID mismatch. Swallowing error to retry message...")
			return nil, nil
		}

		return nil, err
	}

	valset, err := t.paloma.QueryGetEVMValsetByID(ctx, valsetID, t.ChainReferenceID)
	if err != nil {
		return nil, err
	}

	consensusReached := isConsensusReached(ctx, valset, origMessage)
	if !consensusReached {
		return nil, ErrNoConsensus
	}

	con := BuildCompassConsensus(valset, origMessage.Signatures)
	compassArgs := CompassLogicCallArgs{
		LogicContractAddress: common.HexToAddress(msg.GetHexContractAddress()),
		Payload:              msg.GetPayload(),
	}

	args := []any{
		con,
		compassArgs,
		new(big.Int).SetInt64(int64(origMessage.ID)),
		new(big.Int).SetInt64(msg.GetDeadline()),
	}

	logger.WithField("consensus", con).WithField("args", args).Debug("submitting logic call")
	tx, err := t.callCompass(ctx, msg.ExecutionRequirements.EnforceMEVRelay, "submit_logic_call", args)
	if err != nil {
		logger.WithError(err).Error("submitLogicCall: error calling DeployContract")
		isSmartContractError, setErr := t.SetErrorData(ctx, queueTypeName, origMessage.ID, err)
		if setErr != nil {
			return nil, setErr
		}
		if isSmartContractError {
			logger.Debug("smart contract error. recovering...")
			return nil, nil
		}

		return nil, err
	}

	return tx, nil
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
		logger.
			WithError(err).
			WithField("input", constructorInput).
			Error("uploadSmartContract: error calling DeployContract")

		isSmartContractError, setErr := t.SetErrorData(ctx, queueTypeName, origMessage.ID, err)
		if setErr != nil {
			return nil, setErr
		}
		if isSmartContractError {
			logger.Debug("smart contract error. recovering...")
			return nil, nil
		}

		return nil, err
	}

	return tx, nil
}

func (t compass) SetErrorData(ctx context.Context, queueTypeName string, msgID uint64, errToProcess error) (bool, error) {
	var jsonRpcErr rpc.DataError
	if !errors.As(errToProcess, &jsonRpcErr) {
		return false, t.paloma.SetErrorData(ctx, queueTypeName, msgID, []byte(errToProcess.Error()))
	} else {
		liblog.WithContext(ctx).WithFields(
			log.Fields{
				"queue-type-name": queueTypeName,
				"message-id":      msgID,
				"error-message":   jsonRpcErr.Error(),
			},
		).Warn("smart contract returned an error")

		err := t.paloma.SetErrorData(ctx, queueTypeName, msgID, []byte(jsonRpcErr.Error()))
		if err != nil {
			return false, err
		}

		return true, nil
	}
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

func (t compass) processMessages(ctx context.Context, queueTypeName string, msgs []chain.MessageWithSignatures) error {
	var gErr whoops.Group
	logger := liblog.WithContext(ctx).WithField("queue-type-name", queueTypeName)
	for i, rawMsg := range msgs {
		logger = logger.WithField("message-id", rawMsg.ID)

		if ctx.Err() != nil {
			logger.Debug("exiting processing message context")
			break
		}

		var processingErr error
		var tx *ethtypes.Transaction
		msg := rawMsg.Msg.(*evmtypes.Message)
		logger := logger.WithFields(log.Fields{
			"chain-reference-id": t.ChainReferenceID,
			"queue-name":         queueTypeName,
			"msg-id":             rawMsg.ID,
			"message-type":       fmt.Sprintf("%T", msg.GetAction()),
		})
		logger.Debug("processing")

		switch action := msg.GetAction().(type) {
		case *evmtypes.Message_SubmitLogicCall:
			tx, processingErr = t.submitLogicCall(
				ctx,
				queueTypeName,
				action.SubmitLogicCall,
				rawMsg,
			)
		case *evmtypes.Message_UpdateValset:
			logger := logger.WithFields(log.Fields{
				"chain-reference-id":     t.ChainReferenceID,
				"queue-name":             queueTypeName,
				"msg-id":                 rawMsg.ID,
				"msg-bytes-to-sign":      rawMsg.BytesToSign,
				"msg-msg":                rawMsg.Msg,
				"msg-nonce":              rawMsg.Nonce,
				"msg-public-access-data": rawMsg.PublicAccessData,
				"message-type":           "Message_UpdateValset",
			})
			logger.Debug("switch-case-message-update-valset")
			tx, processingErr = t.updateValset(
				ctx,
				queueTypeName,
				action.UpdateValset.Valset,
				rawMsg,
			)
		case *evmtypes.Message_UploadSmartContract:
			logger := logger.WithFields(log.Fields{
				"chain-reference-id":     t.ChainReferenceID,
				"queue-name":             queueTypeName,
				"msg-id":                 rawMsg.ID,
				"msg-bytes-to-sign":      rawMsg.BytesToSign,
				"msg-msg":                rawMsg.Msg,
				"msg-nonce":              rawMsg.Nonce,
				"msg-public-access-data": rawMsg.PublicAccessData,
				"message-type":           "Message_UploadSmartContract",
			})
			logger.Debug("switch-case-message-upload-contract")
			tx, processingErr = t.uploadSmartContract(
				ctx,
				queueTypeName,
				action.UploadSmartContract,
				rawMsg,
			)
		default:
			return ErrUnsupportedMessageType.Format(action)
		}

		processingErr = whoops.Enrich(
			processingErr,
			FieldMessageID.Val(rawMsg.ID),
			FieldMessageType.Val(msg.GetAction()),
		)

		switch {
		case processingErr == nil:
			if tx != nil {
				logger.Debug("setting public access data")
				if err := t.paloma.SetPublicAccessData(ctx, queueTypeName, rawMsg.ID, tx.Hash().Bytes()); err != nil {
					gErr.Add(err)
					return gErr
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

	return gErr.Return()
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

func (t *compass) GetBatchSendEvents(ctx context.Context, orchestrator string) ([]chain.BatchSendEvent, error) {
	filter, err := ethfilter.Factory().
		WithFromBlockNumberProvider(t.evm.FindCurrentBlockNumber).
		WithFromBlockNumberSafetyMargin(1).
		WithTopics([]common.Hash{crypto.Keccak256Hash([]byte("BatchSendEvent(address,uint256,uint256,uint256)"))}).
		WithAddresses(t.smartContractAddr).
		Filter(ctx)
	if err != nil {
		return nil, err
	}

	currentBlockNumber := filter.FromBlock.Int64()
	if t.lastObservedBlockHeights.batchSendEvent == 0 {
		t.lastObservedBlockHeights.batchSendEvent = currentBlockNumber - 10_000
	}

	filter.FromBlock = big.NewInt(t.lastObservedBlockHeights.batchSendEvent)
	filter.ToBlock = big.NewInt(min(t.lastObservedBlockHeights.batchSendEvent+10_000, currentBlockNumber))

	var events []chain.BatchSendEvent

	logs, err := t.evm.GetEthClient().FilterLogs(ctx, filter)
	if err != nil {
		return nil, err
	}

	lastSkywayNonce, err := t.paloma.QueryLastObservedSkywayNonceByAddr(ctx, t.ChainReferenceID, orchestrator)
	if err != nil {
		return nil, err
	}

	for _, ethLog := range logs {
		event, err := t.compassAbi.Unpack("BatchSendEvent", ethLog.Data)
		if err != nil {
			return nil, err
		}

		tokenContract, ok := event[0].(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid token contract")
		}

		batchNonce, ok := event[1].(*big.Int)
		if !ok {
			return nil, fmt.Errorf("invalid batch nonce")
		}

		skywayNonce, ok := event[2].(*big.Int)
		if !ok {
			return nil, fmt.Errorf("invalid skyway nonce")
		}

		eventNonce, ok := event[3].(*big.Int)
		if !ok {
			return nil, fmt.Errorf("invalid event nonce")
		}

		if skywayNonce.Uint64() <= lastSkywayNonce {
			liblog.WithContext(ctx).WithField("last-event-nonce", lastSkywayNonce).WithField("skyway-nonce", skywayNonce.Uint64()).Info("Skipping already observed event...")
			continue
		}

		events = append(events, chain.BatchSendEvent{
			EthBlockHeight: ethLog.BlockNumber,
			EventNonce:     eventNonce.Uint64(),
			BatchNonce:     batchNonce.Uint64(),
			TokenContract:  tokenContract.String(),
			SkywayNonce:    skywayNonce.Uint64(),
		})
	}

	t.lastObservedBlockHeights.batchSendEvent = filter.ToBlock.Int64()

	return events, err
}

func (t *compass) GetSendToPalomaEvents(ctx context.Context, orchestrator string) ([]chain.SendToPalomaEvent, error) {
	filter, err := ethfilter.Factory().
		WithFromBlockNumberProvider(t.evm.FindCurrentBlockNumber).
		WithFromBlockNumberSafetyMargin(1).
		WithTopics([]common.Hash{crypto.Keccak256Hash([]byte("SendToPalomaEvent(address,address,string,uint256,uint256,uint256)"))}).
		WithAddresses(t.smartContractAddr).
		Filter(ctx)
	if err != nil {
		return nil, err
	}

	currentBlockNumber := filter.FromBlock.Int64()
	if t.lastObservedBlockHeights.sendToPalomaEvent == 0 {
		t.lastObservedBlockHeights.sendToPalomaEvent = currentBlockNumber - 10_000
	}

	filter.FromBlock = big.NewInt(t.lastObservedBlockHeights.sendToPalomaEvent)
	filter.ToBlock = big.NewInt(min(t.lastObservedBlockHeights.sendToPalomaEvent+10_000, currentBlockNumber))

	var events []chain.SendToPalomaEvent

	logs, err := t.evm.GetEthClient().FilterLogs(ctx, filter)
	if err != nil {
		return nil, err
	}

	lastSkywayNonce, err := t.paloma.QueryLastObservedSkywayNonceByAddr(ctx, t.ChainReferenceID, orchestrator)
	if err != nil {
		return nil, err
	}

	for _, ethLog := range logs {
		event, err := t.compassAbi.Unpack("SendToPalomaEvent", ethLog.Data)
		if err != nil {
			return nil, err
		}

		tokenContract, ok := event[0].(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid token contract")
		}

		ethSender, ok := event[1].(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid sender address")
		}

		palomaReceiver := event[2].(string)
		if !ok {
			return nil, fmt.Errorf("invalid paloma receiver")
		}

		amount, ok := event[3].(*big.Int)
		if !ok {
			return nil, fmt.Errorf("invalid amount")
		}

		skywayNonce, ok := event[4].(*big.Int)
		if !ok {
			return nil, fmt.Errorf("invalid paloma nonce")
		}

		eventNonce, ok := event[5].(*big.Int)
		if !ok {
			return nil, fmt.Errorf("invalid event nonce")
		}

		if skywayNonce.Uint64() <= lastSkywayNonce {
			liblog.WithContext(ctx).WithField("last-event-nonce", lastSkywayNonce).WithField("skyway-nonce", skywayNonce.Uint64()).Info("Skipping already observed event...")
			continue
		}

		events = append(events, chain.SendToPalomaEvent{
			EthBlockHeight: ethLog.BlockNumber,
			EventNonce:     eventNonce.Uint64(),
			Amount:         amount.Uint64(),
			EthereumSender: ethSender.String(),
			PalomaReceiver: palomaReceiver,
			TokenContract:  tokenContract.String(),
			SkywayNonce:    skywayNonce.Uint64(),
		})
	}

	t.lastObservedBlockHeights.sendToPalomaEvent = filter.ToBlock.Int64()

	return events, err
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

	return t.paloma.AddMessageEvidence(ctx, queueTypeName, rawMsg.ID, &evmtypes.TxExecutedProof{
		SerializedTX: txProof,
	})
}

func (t compass) submitBatchSendToEVMClaim(ctx context.Context, event chain.BatchSendEvent, orchestrator string) error {
	msg := skywaytypes.MsgBatchSendToEthClaim{
		EventNonce:       event.EventNonce,
		EthBlockHeight:   event.EthBlockHeight,
		BatchNonce:       event.BatchNonce,
		TokenContract:    event.TokenContract,
		ChainReferenceId: t.ChainReferenceID,
		Orchestrator:     orchestrator,
		SkywayNonce:      event.SkywayNonce,
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
	}
	return t.paloma.SendSendToPalomaClaim(ctx, msg)
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

func (c compass) callCompass(
	ctx context.Context,
	useMevRelay bool,
	method string,
	arguments []any,
) (*ethtypes.Transaction, error) {
	if c.compassAbi == nil {
		return nil, ErrABINotInitialized
	}
	return c.evm.ExecuteSmartContract(ctx, c.chainID, *c.compassAbi, c.smartContractAddr, useMevRelay, method, arguments)
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

		_, processingErr = t.skywayRelayBatch(ctx, batch)

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

func (t compass) skywayRelayBatch(
	ctx context.Context,
	batch chain.SkywayBatchWithSignatures,
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

		tx, err := t.callCompass(
			ctx,
			false,
			"submit_batch",
			[]any{
				con,
				common.HexToAddress(batch.TokenContract),
				compassArgs,
				new(big.Int).SetInt64(int64(batch.BatchNonce)),
				new(big.Int).SetInt64(int64(batch.GetBatchTimeout())), // TODO : Deadline
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
