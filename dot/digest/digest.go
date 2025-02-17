// Copyright 2021 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package digest

import (
	"context"
	"errors"

	"github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/internal/log"
	"github.com/ChainSafe/gossamer/lib/services"
	"github.com/ChainSafe/gossamer/pkg/scale"
)

var (
	_ services.Service = &Handler{}
)

// Handler is used to handle consensus messages and relevant authority updates to BABE and GRANDPA
type Handler struct {
	ctx    context.Context
	cancel context.CancelFunc

	// interfaces
	blockState   BlockState
	epochState   EpochState
	grandpaState GrandpaState

	// block notification channels
	imported  chan *types.Block
	finalised chan *types.FinalisationInfo

	// GRANDPA changes
	grandpaScheduledChange *grandpaChange
	grandpaForcedChange    *grandpaChange
	grandpaPause           *pause
	grandpaResume          *resume

	logger log.LeveledLogger
}

type grandpaChange struct {
	auths   []types.Authority
	atBlock uint
}

type pause struct {
	atBlock uint
}

type resume struct {
	atBlock uint
}

// NewHandler returns a new Handler
func NewHandler(lvl log.Level, blockState BlockState, epochState EpochState,
	grandpaState GrandpaState) (*Handler, error) {
	imported := blockState.GetImportedBlockNotifierChannel()
	finalised := blockState.GetFinalisedNotifierChannel()

	logger := log.NewFromGlobal(log.AddContext("pkg", "digest"))
	logger.Patch(log.SetLevel(lvl))

	ctx, cancel := context.WithCancel(context.Background())
	return &Handler{
		ctx:          ctx,
		cancel:       cancel,
		blockState:   blockState,
		epochState:   epochState,
		grandpaState: grandpaState,
		imported:     imported,
		finalised:    finalised,
		logger:       logger,
	}, nil
}

// Start starts the Handler
func (h *Handler) Start() error {
	go h.handleBlockImport(h.ctx)
	go h.handleBlockFinalisation(h.ctx)
	return nil
}

// Stop stops the Handler
func (h *Handler) Stop() error {
	h.cancel()
	h.blockState.FreeImportedBlockNotifierChannel(h.imported)
	h.blockState.FreeFinalisedNotifierChannel(h.finalised)
	return nil
}

// NextGrandpaAuthorityChange returns the block number of the next upcoming grandpa authorities change.
// It returns 0 if no change is scheduled.
func (h *Handler) NextGrandpaAuthorityChange() (next uint) {
	next = ^uint(0)

	if h.grandpaScheduledChange != nil {
		next = h.grandpaScheduledChange.atBlock
	}

	if h.grandpaForcedChange != nil && h.grandpaForcedChange.atBlock < next {
		next = h.grandpaForcedChange.atBlock
	}

	if h.grandpaPause != nil && h.grandpaPause.atBlock < next {
		next = h.grandpaPause.atBlock
	}

	if h.grandpaResume != nil && h.grandpaResume.atBlock < next {
		next = h.grandpaResume.atBlock
	}

	return next
}

// HandleDigests handles consensus digests for an imported block
func (h *Handler) HandleDigests(header *types.Header) {
	for i, d := range header.Digest.Types {
		val, ok := d.Value().(types.ConsensusDigest)
		if ok {
			err := h.handleConsensusDigest(&val, header)
			if err != nil {
				h.logger.Errorf("cannot handle digests for block number %d, index %d, digest %s: %s",
					header.Number, i, d.Value(), err)
			}
		}
	}
}

func (h *Handler) handleConsensusDigest(d *types.ConsensusDigest, header *types.Header) error {
	switch d.ConsensusEngineID {
	case types.GrandpaEngineID:
		data := types.NewGrandpaConsensusDigest()
		err := scale.Unmarshal(d.Data, &data)
		if err != nil {
			return err
		}
		err = h.handleGrandpaConsensusDigest(data, header)
		if err != nil {
			return err
		}
		return nil
	case types.BabeEngineID:
		data := types.NewBabeConsensusDigest()
		err := scale.Unmarshal(d.Data, &data)
		if err != nil {
			return err
		}
		err = h.handleBabeConsensusDigest(data, header)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("unknown consensus engine ID")
}

func (h *Handler) handleGrandpaConsensusDigest(digest scale.VaryingDataType, header *types.Header) error {
	switch val := digest.Value().(type) {
	case types.GrandpaScheduledChange:
		return h.handleScheduledChange(val, header)
	case types.GrandpaForcedChange:
		return h.handleForcedChange(val, header)
	case types.GrandpaOnDisabled:
		return nil // do nothing, as this is not implemented in substrate
	case types.GrandpaPause:
		return h.handlePause(val)
	case types.GrandpaResume:
		return h.handleResume(val)
	}

	return errors.New("invalid consensus digest data")
}

func (h *Handler) handleBabeConsensusDigest(digest scale.VaryingDataType, header *types.Header) error {
	switch val := digest.Value().(type) {
	case types.NextEpochData:
		h.logger.Debugf("handling BABENextEpochData data: %v", digest)
		return h.handleNextEpochData(val, header)
	case types.BABEOnDisabled:
		return h.handleBABEOnDisabled(val, header)
	case types.NextConfigData:
		h.logger.Debugf("handling BABENextConfigData data: %v", digest)
		return h.handleNextConfigData(val, header)
	}

	return errors.New("invalid consensus digest data")
}

func (h *Handler) handleBlockImport(ctx context.Context) {
	for {
		select {
		case block := <-h.imported:
			if block == nil {
				continue
			}

			err := h.handleGrandpaChangesOnImport(block.Header.Number)
			if err != nil {
				h.logger.Errorf("failed to handle grandpa changes on block import: %s", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (h *Handler) handleBlockFinalisation(ctx context.Context) {
	for {
		select {
		case info := <-h.finalised:
			if info == nil {
				continue
			}

			err := h.handleGrandpaChangesOnFinalization(info.Header.Number)
			if err != nil {
				h.logger.Errorf("failed to handle grandpa changes on block finalisation: %s", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (h *Handler) handleGrandpaChangesOnImport(num uint) error {
	resume := h.grandpaResume
	if resume != nil && num >= resume.atBlock {
		h.grandpaResume = nil
	}

	fc := h.grandpaForcedChange
	if fc != nil && num >= fc.atBlock {
		err := h.grandpaState.IncrementSetID()
		if err != nil {
			return err
		}

		h.grandpaForcedChange = nil
		curr, err := h.grandpaState.GetCurrentSetID()
		if err != nil {
			return err
		}

		h.logger.Debugf("incremented grandpa set id %d", curr)
	}

	return nil
}

func (h *Handler) handleGrandpaChangesOnFinalization(num uint) error {
	pause := h.grandpaPause
	if pause != nil && num >= pause.atBlock {
		h.grandpaPause = nil
	}

	sc := h.grandpaScheduledChange
	if sc != nil && num >= sc.atBlock {
		err := h.grandpaState.IncrementSetID()
		if err != nil {
			return err
		}

		h.grandpaScheduledChange = nil
		curr, err := h.grandpaState.GetCurrentSetID()
		if err != nil {
			return err
		}

		h.logger.Debugf("incremented grandpa set id %d", curr)
	}

	// if blocks get finalised before forced change takes place, disregard it
	h.grandpaForcedChange = nil
	return nil
}

func (h *Handler) handleScheduledChange(sc types.GrandpaScheduledChange, header *types.Header) error {
	curr, err := h.blockState.BestBlockHeader()
	if err != nil {
		return err
	}

	if h.grandpaScheduledChange != nil {
		return nil
	}

	h.logger.Debugf("handling GrandpaScheduledChange data: %v", sc)

	c, err := newGrandpaChange(sc.Auths, sc.Delay, curr.Number)
	if err != nil {
		return err
	}

	h.grandpaScheduledChange = c

	auths, err := types.GrandpaAuthoritiesRawToAuthorities(sc.Auths)
	if err != nil {
		return err
	}
	h.logger.Debugf("setting GrandpaScheduledChange at block %d",
		header.Number+uint(sc.Delay))
	return h.grandpaState.SetNextChange(
		types.NewGrandpaVotersFromAuthorities(auths),
		header.Number+uint(sc.Delay),
	)
}

func (h *Handler) handleForcedChange(fc types.GrandpaForcedChange, header *types.Header) error {
	if header == nil {
		return errors.New("header is nil")
	}

	if h.grandpaForcedChange != nil {
		return errors.New("already have forced change scheduled")
	}

	h.logger.Debugf("handling GrandpaForcedChange with data %v", fc)

	c, err := newGrandpaChange(fc.Auths, fc.Delay, header.Number)
	if err != nil {
		return err
	}

	h.grandpaForcedChange = c

	auths, err := types.GrandpaAuthoritiesRawToAuthorities(fc.Auths)
	if err != nil {
		return err
	}

	h.logger.Debugf("setting GrandpaForcedChange at block %d",
		header.Number+uint(fc.Delay))
	return h.grandpaState.SetNextChange(
		types.NewGrandpaVotersFromAuthorities(auths),
		header.Number+uint(fc.Delay),
	)
}

func (h *Handler) handlePause(p types.GrandpaPause) error {
	curr, err := h.blockState.BestBlockHeader()
	if err != nil {
		return err
	}

	h.grandpaPause = &pause{
		atBlock: curr.Number + uint(p.Delay),
	}

	return h.grandpaState.SetNextPause(h.grandpaPause.atBlock)
}

func (h *Handler) handleResume(r types.GrandpaResume) error {
	curr, err := h.blockState.BestBlockHeader()
	if err != nil {
		return err
	}

	h.grandpaResume = &resume{
		atBlock: curr.Number + uint(r.Delay),
	}

	return h.grandpaState.SetNextResume(h.grandpaResume.atBlock)
}

func newGrandpaChange(raw []types.GrandpaAuthoritiesRaw, delay uint32, currBlock uint) (*grandpaChange, error) {
	auths, err := types.GrandpaAuthoritiesRawToAuthorities(raw)
	if err != nil {
		return nil, err
	}

	return &grandpaChange{
		auths:   auths,
		atBlock: currBlock + uint(delay),
	}, nil
}

func (h *Handler) handleBABEOnDisabled(_ types.BABEOnDisabled, _ *types.Header) error {
	h.logger.Debug("handling BABEOnDisabled")
	return nil
}

func (h *Handler) handleNextEpochData(act types.NextEpochData, header *types.Header) error {
	currEpoch, err := h.epochState.GetEpochForBlock(header)
	if err != nil {
		return err
	}

	// set EpochState epoch data for upcoming epoch
	data, err := act.ToEpochData()
	if err != nil {
		return err
	}

	h.logger.Debugf("setting data for block number %d and epoch %d with data: %v",
		header.Number, currEpoch+1, data)
	return h.epochState.SetEpochData(currEpoch+1, data)
}

func (h *Handler) handleNextConfigData(config types.NextConfigData, header *types.Header) error {
	currEpoch, err := h.epochState.GetEpochForBlock(header)
	if err != nil {
		return err
	}

	h.logger.Debugf("setting BABE config data for block number %d and epoch %d with data: %v",
		header.Number, currEpoch+1, config.ToConfigData())
	// set EpochState config data for upcoming epoch
	return h.epochState.SetConfigData(currEpoch+1, config.ToConfigData())
}
