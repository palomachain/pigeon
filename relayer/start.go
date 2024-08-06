package relayer

import (
	"context"
	"sync"
	"time"

	"github.com/palomachain/paloma/util/libvalid"
	"github.com/palomachain/pigeon/internal/liblog"
	"github.com/palomachain/pigeon/relayer/heartbeat"
	log "github.com/sirupsen/logrus"
)

const (
	updateExternalChainsLoopInterval = 1 * time.Minute
	signMessagesLoopInterval         = 500 * time.Millisecond
	relayMessagesLoopInterval        = 500 * time.Millisecond
	attestMessagesLoopInterval       = 500 * time.Millisecond
	checkStakingLoopInterval         = 5 * time.Second

	skywaySignBatchesLoopInterval  = 5 * time.Second
	skywayRelayBatchesLoopInterval = 5 * time.Second
	skywayEventWatcherLoopInterval = 1 * time.Minute
)

func (r *Relayer) checkStaking(ctx context.Context, locker sync.Locker) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	log.Info("checking if validator is staking")

	err := r.isStaking(ctx)
	if err == nil {
		log.Info("validator is staking")
		r.staking = true
	} else {
		log.Warn("validator is not staking... waiting")
		r.staking = false
	}
	return nil
}

func (r *Relayer) startProcess(ctx context.Context, name string, locker sync.Locker, tickerInterval time.Duration, requiresStaking bool, process func(context.Context, sync.Locker) error) {
	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	logger := liblog.WithContext(ctx).WithField("component", "procmon").WithField("process", name)
	for {
		select {
		case <-ctx.Done():
			logger.Warn("exiting due to context being done")
			return
		case <-ticker.C:
			if !requiresStaking || r.staking {
				jCtx := liblog.MustEnrichContext(ctx)
				err := process(jCtx, locker)
				if err != nil {
					liblog.WithContext(jCtx).WithField("component", "procmon").WithField("process", name).WithError(err).Errorf("Failed to execute process: %v", err)
				} else {
					liblog.WithContext(jCtx).WithField("component", "procmon").WithField("process", name).Debug("Process executed")
				}
			} else {
				logger.Debug("validor not staking, skipping process execution...")
			}
		}
	}
}

// Start starts the relayer. It's responsible for handling the communication
// with Paloma and other chains.
func (r *Relayer) Start(ctx context.Context) error {
	log.Info("starting pigeon")
	var locker sync.Mutex

	_ = r.checkStaking(ctx, &locker)

	// Start background goroutines to run separately from each other
	go r.startProcess(ctx, "Check staking", &locker, checkStakingLoopInterval, false, r.checkStaking)
	go r.startProcess(ctx, "Update external chain infos", &locker, updateExternalChainsLoopInterval, true, r.UpdateExternalChainInfos)
	go r.startProcess(ctx, "Sign messages", &locker, signMessagesLoopInterval, true, r.SignMessages)
	go r.startProcess(ctx, "Relay messages", &locker, relayMessagesLoopInterval, true, r.RelayMessages)
	go r.startProcess(ctx, "Attest messages", &locker, attestMessagesLoopInterval, true, r.AttestMessages)

	if !libvalid.IsNil(r.mevClient) {
		go r.startProcess(ctx, "[MEV] Client heartbeat", &locker, r.mevClient.GetHealthprobeInterval(), false, r.mevClient.KeepAlive)
	}

	// Start skyway background goroutines to run separately from each other
	go r.startProcess(ctx, "[Skyway] Sign batches", &locker, skywaySignBatchesLoopInterval, true, r.SkywaySignBatches)
	go r.startProcess(ctx, "[Skyway] Relay batches", &locker, skywayRelayBatchesLoopInterval, true, r.SkywayRelayBatches)
	go r.startProcess(ctx, "[Skyway] Handle Events", &locker, skywayEventWatcherLoopInterval, true, r.SkywayHandleEvents)

	// Setup heartbeat to Paloma
	heart := heartbeat.New(
		r.palomaClient.QueryGetValidatorAliveUntilBlockHeight,
		r.palomaClient.BlockHeight,
		r.palomaClient.KeepValidatorAlive,
		r.relayerConfig.KeepAliveBlockThreshold,
		r.appVersion,
		&locker)

	// Immediately send a keep alive to Paloma during startup
	_ = heart.Beat(liblog.MustEnrichContext(ctx), &locker)

	// Start the foreground process
	r.startProcess(ctx, "Keep alive", &locker, r.relayerConfig.KeepAliveLoopTimeout, false, heart.Beat)
	return nil
}
