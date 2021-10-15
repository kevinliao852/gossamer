// Copyright 2019 ChainSafe Systems (ON) Corp.
// This file is part of gossamer.
//
// The gossamer library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The gossamer library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the gossamer library. If not, see <http://www.gnu.org/licenses/>.

package sync

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/gossamer/dot/network"
	"github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
)

const (
	// maxResponseSize is maximum number of block data a BlockResponse message can contain
	maxResponseSize = 128
)

// CreateBlockResponse creates a block response message from a block request message
func (s *Service) CreateBlockResponse(req *network.BlockRequestMessage) (*network.BlockResponseMessage, error) {
	switch req.Direction {
	case network.Ascending:
		return s.handleAscendingRequest(req)
	case network.Descending:
		return s.handleDescendingRequest(req)
	default:
		return nil, errors.New("invalid request direction")
	}

	// var (
	// 	startHash, endHash     common.Hash
	// 	startHeader, endHeader *types.Header
	// 	err                    error
	// 	respSize               uint32
	// )

	// if blockRequest.Max != nil {
	// 	respSize = *blockRequest.Max
	// 	if respSize > maxResponseSize {
	// 		respSize = maxResponseSize
	// 	}
	// } else {
	// 	respSize = maxResponseSize
	// }

	// switch startBlock := blockRequest.StartingBlock.Value().(type) {
	// case uint64:
	// 	if startBlock == 0 {
	// 		startBlock = 1
	// 	}

	// 	block, err := s.blockState.GetBlockByNumber(big.NewInt(0).SetUint64(startBlock)) //nolint
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get start block %d for request: %w", startBlock, err)
	// 	}

	// 	startHeader = &block.Header
	// 	startHash = block.Header.Hash()
	// case common.Hash:
	// 	startHash = startBlock
	// 	startHeader, err = s.blockState.GetHeader(startHash)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get start block %s for request: %w", startHash, err)
	// 	}
	// default:
	// 	return nil, ErrInvalidBlockRequest
	// }

	// if blockRequest.EndBlockHash != nil {
	// 	endHash = *blockRequest.EndBlockHash
	// 	endHeader, err = s.blockState.GetHeader(endHash)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get end block %s for request: %w", endHash, err)
	// 	}
	// } else {
	// 	endNumber := big.NewInt(0).Add(startHeader.Number, big.NewInt(int64(respSize-1)))
	// 	bestBlockNumber, err := s.blockState.BestBlockNumber()
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get best block %d for request: %w", bestBlockNumber, err)
	// 	}

	// 	if endNumber.Cmp(bestBlockNumber) == 1 {
	// 		endNumber = bestBlockNumber
	// 	}

	// 	endBlock, err := s.blockState.GetBlockByNumber(endNumber)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get end block %d for request: %w", endNumber, err)
	// 	}
	// 	endHeader = &endBlock.Header
	// 	endHash = endHeader.Hash()
	// }

	// logger.Debug("handling BlockRequestMessage", "start", startHeader.Number, "end", endHeader.Number, "startHash", startHash, "endHash", endHash)

	// responseData := []*types.BlockData{}

	// switch blockRequest.Direction {
	// case network.Ascending:
	// 	for i := startHeader.Number.Int64(); i <= endHeader.Number.Int64(); i++ {
	// 		blockData, err := s.getBlockData(big.NewInt(i), blockRequest.RequestedData)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		responseData = append(responseData, blockData)
	// 	}
	// // case network.Descending:
	// // 	for i := startHeader.Number.Int64(); i >= endHeader.Number.Int64(); i-- {
	// // 		blockData, err := s.getBlockData(big.NewInt(i), blockRequest.RequestedData)
	// // 		if err != nil {
	// // 			return nil, err
	// // 		}
	// // 		responseData = append(responseData, blockData)
	// // 	}
	// default:
	// 	return nil, errors.New("invalid BlockRequest direction")
	// }

	// logger.Debug("sending BlockResponseMessage", "start", startHeader.Number, "end", endHeader.Number)
	// return &network.BlockResponseMessage{
	// 	BlockData: responseData,
	// }, nil
}

func (s *Service) handleAscendingRequest(req *network.BlockRequestMessage) (*network.BlockResponseMessage, error) {
	// var (
	// 	startHash, endHash     common.Hash
	// 	startHeader, endHeader *types.Header
	// 	err                    error
	// 	respSize               uint32
	// )

	// if blockRequest.Max != nil {
	// 	respSize = *blockRequest.Max
	// 	if respSize > maxResponseSize {
	// 		respSize = maxResponseSize
	// 	}
	// } else {
	// 	respSize = maxResponseSize
	// }

	// switch startBlock := blockRequest.StartingBlock.Value().(type) {
	// case uint64:
	// 	if startBlock == 0 {
	// 		startBlock = 1
	// 	}

	// 	block, err := s.blockState.GetBlockByNumber(big.NewInt(0).SetUint64(startBlock)) //nolint
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get start block %d for request: %w", startBlock, err)
	// 	}

	// 	startHeader = &block.Header
	// 	startHash = block.Header.Hash()
	// case common.Hash:
	// 	startHash = startBlock
	// 	startHeader, err = s.blockState.GetHeader(startHash)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get start block %s for request: %w", startHash, err)
	// 	}
	// default:
	// 	return nil, ErrInvalidBlockRequest
	// }

	// if blockRequest.EndBlockHash != nil {
	// 	endHash = *blockRequest.EndBlockHash
	// 	endHeader, err = s.blockState.GetHeader(endHash)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get end block %s for request: %w", endHash, err)
	// 	}
	// } else {
	// 	switch req.Direction {
	// 	case network.Ascending:
	// 		endNumber := big.NewInt(0).Add(startHeader.Number, big.NewInt(int64(respSize-1)))
	// 		bestBlockNumber, err := s.blockState.BestBlockNumber()
	// 		if err != nil {
	// 			return nil, fmt.Errorf("failed to get best block %d for request: %w", bestBlockNumber, err)
	// 		}

	// 		if endNumber.Cmp(bestBlockNumber) == 1 {
	// 			endNumber = bestBlockNumber
	// 		}

	// 		endBlock, err := s.blockState.GetBlockByNumber(endNumber)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("failed to get end block %d for request: %w", endNumber, err)
	// 		}
	// 		endHeader = &endBlock.Header
	// 		endHash = endHeader.Hash()
	// 	case network.Descending:

	// 	}

	// }

	// logger.Debug("handling BlockRequestMessage", "start", startHeader.Number, "end", endHeader.Number, "startHash", startHash, "endHash", endHash)

	// responseData := []*types.BlockData{}

	// switch blockRequest.Direction {
	// case network.Ascending:
	// 	for i := startHeader.Number.Int64(); i <= endHeader.Number.Int64(); i++ {
	// 		blockData, err := s.getBlockData(big.NewInt(i), blockRequest.RequestedData)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		responseData = append(responseData, blockData)
	// 	}
	// case network.Descending:
	// 	for i := startHeader.Number.Int64(); i >= endHeader.Number.Int64(); i-- {
	// 		blockData, err := s.getBlockData(big.NewInt(i), blockRequest.RequestedData)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		responseData = append(responseData, blockData)
	// 	}
	// default:
	// 	return nil, errors.New("invalid BlockRequest direction")
	// }

	// logger.Debug("sending BlockResponseMessage", "start", startHeader.Number, "end", endHeader.Number)
	// return &network.BlockResponseMessage{
	// 	BlockData: responseData,
	// }, nil

	return nil, nil
}

func (s *Service) handleDescendingRequest(req *network.BlockRequestMessage) (*network.BlockResponseMessage, error) {
	var (
		startHash              *common.Hash
		endHash                = req.EndBlockHash
		startNumber, endNumber uint64
		//err                    error
		max uint32 = maxResponseSize
	)

	if req.Max != nil && *req.Max < maxResponseSize {
		max = *req.Max
	}

	// if blockRequest.EndBlockHash != nil {
	// 	endHash = blockRequest.EndBlockHash
	// } else {
	// 	endNumber := big.NewInt(0).Sub(startHeader.Number, big.NewInt(int64(respSize+1)))
	// 	endBlock, err := s.blockState.GetBlockByNumber(endNumber)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get end block %d for request: %w", endNumber, err)
	// 	}
	// 	endHeader = &endBlock.Header
	// 	endHash = endHeader.Hash()
	// }

	switch startBlock := req.StartingBlock.Value().(type) {
	case uint64:
		if startBlock == 0 {
			startNumber = 1
		}

		bestBlockNumber, err := s.blockState.BestBlockNumber()
		if err != nil {
			return nil, fmt.Errorf("failed to get best block %d for request: %w", bestBlockNumber, err)
		}

		// if request start is higher than our best block, only return blocks from our best block and below
		if bestBlockNumber.Uint64() < startBlock {
			startNumber = bestBlockNumber.Uint64()
		}

		// // make sure we have the start block
		// block, err := s.blockState.GetBlockByNumber(big.NewInt(0).SetUint64(startBlock)) //nolint
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to get start block %d for request: %w", startBlock, err)
		// }

		// startHeader = &block.Header
		// startHash = block.Header.Hash()
		startNumber = startBlock
	case common.Hash:
		startHash = &startBlock

		// make sure we actually have the starting block
		header, err := s.blockState.GetHeader(*startHash)
		if err != nil {
			return nil, fmt.Errorf("failed to get start block %s for request: %w", startHash, err)
		}

		startNumber = header.Number.Uint64()
	default:
		return nil, ErrInvalidBlockRequest
	}

	// end hash provided, need to determine start hash that is descendant of end hash
	if endHash != nil {
		if startHash != nil {
			// check if provided start hash is descendant of provided end hash
			is, err := s.blockState.IsDescendantOf(*endHash, *startHash)
			if err != nil {
				return nil, err
			}

			if !is {
				return nil, errors.New("request start hash is not descendant of end hash")
			}
		}

		// get block on canonical chain by start number
		hash, err := s.blockState.GetHashByNumber(big.NewInt(int64(startNumber)))
		if err != nil {
			return nil, err
		}

		// check if it's a descendant of the provided end hash
		is, err := s.blockState.IsDescendantOf(*endHash, hash)
		if err != nil {
			return nil, err
		}

		if !is {
			// if it's not a descedant, search for a block that has number=startNumber that is
			// a descendant of the end block
			hashes, err := s.blockState.GetAllBlocksAtNumber(big.NewInt(int64(startNumber)))
			if err != nil {
				return nil, fmt.Errorf("failed to get blocks at number %d: %w", startNumber, err)
			}

			for _, hash := range hashes {
				is, err := s.blockState.IsDescendantOf(*endHash, hash)
				if err != nil {
					continue
				}

				if is {
					// this sets the startHash to whatever the first block we find with the startNumber
					// is, however there might be multiple blocks that fit this criteria
					startHash = &hash
					break
				}
			}
		} else {
			// if it is, set startHash to our block w/ startNumber
			startHash = &hash
		}
	}

	// end hash is not provided, calculate end by number
	if endHash == nil {
		endNumber = startNumber - uint64(max+1)
		endHeader, err := s.blockState.GetHeaderByNumber(big.NewInt(int64(endNumber)))
		if err != nil {
			return nil, fmt.Errorf("failed to get end block %d for request: %w", endNumber, err)
		}

		hash := endHeader.Hash()
		endHash = &hash
	}

	if startHash == nil && endHash == nil {
		return s.handleDescendingByNumber(startNumber, endNumber, req.RequestedData)
	}

	if startHash == nil {
		panic("startHash is nil!")
	}

	if endHash == nil {
		panic("endHash is nil!")
	}

	return s.handleDescendingByHash(*startHash, *endHash, max, req.RequestedData)
}

func (s *Service) handleDescendingByNumber(start, end uint64, requestedData byte) (*network.BlockResponseMessage, error) {
	var err error
	data := make([]*types.BlockData, start-end)

	for i := start; i >= end; i-- {
		data[i], err = s.getBlockDataByNumber(big.NewInt(int64(i)), requestedData)
		if err != nil {
			return nil, err
		}
	}

	return &network.BlockResponseMessage{
		BlockData: data,
	}, nil
}

func (s *Service) handleDescendingByHash(start, end common.Hash, max uint32, requestedData byte) (*network.BlockResponseMessage, error) {
	subchain, err := s.blockState.SubChain(start, end)
	if uint32(len(subchain)) > max {
		subchain = subchain[:max]
	}

	data := make([]*types.BlockData, len(subchain))

	for i, hash := range subchain {
		data[i], err = s.getBlockData(hash, requestedData)
		if err != nil {
			return nil, err
		}
	}

	return &network.BlockResponseMessage{
		BlockData: data,
	}, nil
}

func (s *Service) getBlockDataByNumber(num *big.Int, requestedData byte) (*types.BlockData, error) {
	hash, err := s.blockState.GetHashByNumber(num)
	if err != nil {
		return nil, err
	}

	return s.getBlockData(hash, requestedData)
}

func (s *Service) getBlockData(hash common.Hash, requestedData byte) (*types.BlockData, error) {
	var err error
	blockData := &types.BlockData{
		Hash: hash,
	}

	if requestedData == 0 {
		return blockData, nil
	}

	if (requestedData & network.RequestedDataHeader) == 1 {
		blockData.Header, err = s.blockState.GetHeader(hash)
		if err != nil {
			logger.Debug("failed to get header for block", "hash", hash, "error", err)
		}
	}

	if (requestedData&network.RequestedDataBody)>>1 == 1 {
		blockData.Body, err = s.blockState.GetBlockBody(hash)
		if err != nil {
			logger.Debug("failed to get body for block", "hash", hash, "error", err)
		}
	}

	if (requestedData&network.RequestedDataReceipt)>>2 == 1 {
		retData, err := s.blockState.GetReceipt(hash)
		if err == nil && retData != nil {
			blockData.Receipt = &retData
		}
	}

	if (requestedData&network.RequestedDataMessageQueue)>>3 == 1 {
		retData, err := s.blockState.GetMessageQueue(hash)
		if err == nil && retData != nil {
			blockData.MessageQueue = &retData
		}
	}

	if (requestedData&network.RequestedDataJustification)>>4 == 1 {
		retData, err := s.blockState.GetJustification(hash)
		if err == nil && retData != nil {
			blockData.Justification = &retData
		}
	}

	return blockData, nil
}
