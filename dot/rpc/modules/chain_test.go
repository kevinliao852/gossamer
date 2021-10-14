// Copyright 2020 ChainSafe Systems (ON) Corp.
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

package modules

import (
	"github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"reflect"
	"testing"
)

func TestChainModule_GetBlock(t *testing.T) {
	mockBlockAPI := NewMockBlockAPI()

	mockBlockAPI.On("GetBlockByHash", mock.Anything).Return(&types.Block{}, nil)

	bhash, err := common.HexToHash("0xea374832a2c3997280d2772c10e6e5b0b493ccd3d09c0ab14050320e34076c2c")
	require.NoError(t, err)

	res := &ChainBlockResponse{}
	req := &ChainHashRequest{Bhash: &bhash}

	type fields struct {
		blockAPI BlockAPI
	}
	type args struct {
		r   *http.Request
		req *ChainHashRequest
		res *ChainBlockResponse
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "error path",
			fields:  fields{mockBlockAPI},
			args:    args{
				r: nil,
				req: req,
				res: res,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ChainModule{
				blockAPI: tt.fields.blockAPI,
			}
			if err := cm.GetBlock(tt.args.r, tt.args.req, tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("GetBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChainModule_GetBlockHash(t *testing.T) {
	type fields struct {
		blockAPI BlockAPI
	}
	type args struct {
		r   *http.Request
		req *ChainBlockNumberRequest
		res *ChainHashResponse
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ChainModule{
				blockAPI: tt.fields.blockAPI,
			}
			if err := cm.GetBlockHash(tt.args.r, tt.args.req, tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("GetBlockHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChainModule_GetFinalizedHead(t *testing.T) {
	type fields struct {
		blockAPI BlockAPI
	}
	type args struct {
		r   *http.Request
		req *EmptyRequest
		res *ChainHashResponse
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ChainModule{
				blockAPI: tt.fields.blockAPI,
			}
			if err := cm.GetFinalizedHead(tt.args.r, tt.args.req, tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("GetFinalizedHead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChainModule_GetFinalizedHeadByRound(t *testing.T) {
	type fields struct {
		blockAPI BlockAPI
	}
	type args struct {
		r   *http.Request
		req *ChainFinalizedHeadRequest
		res *ChainHashResponse
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ChainModule{
				blockAPI: tt.fields.blockAPI,
			}
			if err := cm.GetFinalizedHeadByRound(tt.args.r, tt.args.req, tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("GetFinalizedHeadByRound() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChainModule_GetHeader(t *testing.T) {
	type fields struct {
		blockAPI BlockAPI
	}
	type args struct {
		r   *http.Request
		req *ChainHashRequest
		res *ChainBlockHeaderResponse
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ChainModule{
				blockAPI: tt.fields.blockAPI,
			}
			if err := cm.GetHeader(tt.args.r, tt.args.req, tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("GetHeader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChainModule_SubscribeFinalizedHeads(t *testing.T) {
	type fields struct {
		blockAPI BlockAPI
	}
	type args struct {
		r   *http.Request
		req *EmptyRequest
		res *ChainBlockHeaderResponse
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ChainModule{
				blockAPI: tt.fields.blockAPI,
			}
			if err := cm.SubscribeFinalizedHeads(tt.args.r, tt.args.req, tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("SubscribeFinalizedHeads() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChainModule_SubscribeNewHead(t *testing.T) {
	type fields struct {
		blockAPI BlockAPI
	}
	type args struct {
		r   *http.Request
		req *EmptyRequest
		res *ChainBlockHeaderResponse
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ChainModule{
				blockAPI: tt.fields.blockAPI,
			}
			if err := cm.SubscribeNewHead(tt.args.r, tt.args.req, tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("SubscribeNewHead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChainModule_SubscribeNewHeads(t *testing.T) {
	type fields struct {
		blockAPI BlockAPI
	}
	type args struct {
		r   *http.Request
		req *EmptyRequest
		res *ChainBlockHeaderResponse
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ChainModule{
				blockAPI: tt.fields.blockAPI,
			}
			if err := cm.SubscribeNewHeads(tt.args.r, tt.args.req, tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("SubscribeNewHeads() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChainModule_hashLookup(t *testing.T) {
	type fields struct {
		blockAPI BlockAPI
	}
	type args struct {
		req *ChainHashRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   common.Hash
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ChainModule{
				blockAPI: tt.fields.blockAPI,
			}
			if got := cm.hashLookup(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hashLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainModule_lookupHashByInterface(t *testing.T) {
	type fields struct {
		blockAPI BlockAPI
	}
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ChainModule{
				blockAPI: tt.fields.blockAPI,
			}
			got, err := cm.lookupHashByInterface(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("lookupHashByInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("lookupHashByInterface() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainModule_unwindRequest(t *testing.T) {
	type fields struct {
		blockAPI BlockAPI
	}
	type args struct {
		req interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ChainModule{
				blockAPI: tt.fields.blockAPI,
			}
			got, err := cm.unwindRequest(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("unwindRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unwindRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeaderToJSON(t *testing.T) {
	type args struct {
		header types.Header
	}
	tests := []struct {
		name    string
		args    args
		want    ChainBlockHeaderResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HeaderToJSON(tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("HeaderToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HeaderToJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewChainModule(t *testing.T) {
	type args struct {
		api BlockAPI
	}
	tests := []struct {
		name string
		args args
		want *ChainModule
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChainModule(tt.args.api); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChainModule() = %v, want %v", got, tt.want)
			}
		})
	}
}