// Copyright 2021 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package main

import (
	"io"
	"testing"

	"github.com/ChainSafe/gossamer/chain/gssmr"
	"github.com/ChainSafe/gossamer/dot"

	ctoml "github.com/ChainSafe/gossamer/dot/config/toml"
	"github.com/ChainSafe/gossamer/internal/log"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

// TestExportCommand test "gossamer export --config"
func TestExportCommand(t *testing.T) {
	testCfg, testConfigFile := newTestConfigWithFile(t)
	testDir := testCfg.Global.BasePath
	genFile := dot.NewTestGenesisRawFile(t, testCfg)

	testApp := cli.NewApp()
	testApp.Writer = io.Discard

	testName := "testnode"
	testBootnode := "bootnode"
	testProtocol := "/protocol/test/0"
	testConfig := testConfigFile.Name()

	testcases := []struct {
		description string
		flags       []string
		values      []interface{}
		expected    *dot.Config
	}{
		{
			"Test gossamer export --config --genesis --basepath --name --log --force",
			[]string{"config", "genesis", "basepath", "name", "log", "force"},
			[]interface{}{testConfig, genFile, testDir, testName, log.Info.String(), "true"},
			&dot.Config{
				Global: dot.GlobalConfig{
					Name:           testName,
					ID:             testCfg.Global.ID,
					BasePath:       testCfg.Global.BasePath,
					LogLvl:         log.Info,
					PublishMetrics: testCfg.Global.PublishMetrics,
					MetricsPort:    testCfg.Global.MetricsPort,
				},
				Log: dot.LogConfig{
					CoreLvl:           log.Info,
					SyncLvl:           log.Info,
					NetworkLvl:        log.Info,
					RPCLvl:            log.Info,
					StateLvl:          log.Info,
					RuntimeLvl:        log.Info,
					BlockProducerLvl:  log.Info,
					FinalityGadgetLvl: log.Info,
				},
				Init: dot.InitConfig{
					Genesis: genFile,
				},
				Account: testCfg.Account,
				Core:    testCfg.Core,
				Network: dot.NetworkConfig{
					Port:              testCfg.Network.Port,
					Bootnodes:         testCfg.Network.Bootnodes,
					ProtocolID:        testCfg.Network.ProtocolID,
					NoBootstrap:       testCfg.Network.NoBootstrap,
					NoMDNS:            testCfg.Network.NoMDNS,
					DiscoveryInterval: testCfg.Network.DiscoveryInterval,
					MinPeers:          testCfg.Network.MinPeers,
					MaxPeers:          testCfg.Network.MaxPeers,
				},
				RPC:   testCfg.RPC,
				Pprof: testCfg.Pprof,
			},
		},
		{
			"Test gossamer export --config --genesis --bootnodes --log --force",
			[]string{"config", "genesis", "bootnodes", "name", "force", "pruning", "retain-blocks"},
			[]interface{}{
				testConfig, genFile, testBootnode,
				"Gossamer", "true", gssmr.DefaultPruningMode,
				gssmr.DefaultRetainBlocks},
			&dot.Config{
				Global: testCfg.Global,
				Init: dot.InitConfig{
					Genesis: genFile,
				},
				Log: dot.LogConfig{
					CoreLvl:           log.Info,
					SyncLvl:           log.Info,
					NetworkLvl:        log.Info,
					RPCLvl:            log.Info,
					StateLvl:          log.Info,
					RuntimeLvl:        log.Info,
					BlockProducerLvl:  log.Info,
					FinalityGadgetLvl: log.Info,
				},
				Account: testCfg.Account,
				Core:    testCfg.Core,
				Network: dot.NetworkConfig{
					Port:              testCfg.Network.Port,
					Bootnodes:         []string{testBootnode},
					ProtocolID:        testCfg.Network.ProtocolID,
					NoBootstrap:       testCfg.Network.NoBootstrap,
					NoMDNS:            testCfg.Network.NoMDNS,
					DiscoveryInterval: testCfg.Network.DiscoveryInterval,
					MinPeers:          testCfg.Network.MinPeers,
					MaxPeers:          testCfg.Network.MaxPeers,
				},
				RPC:   testCfg.RPC,
				Pprof: testCfg.Pprof,
			},
		},
		{
			"Test gossamer export --config --genesis --protocol --log --force",
			[]string{"config", "genesis", "protocol", "force", "name", "pruning", "retain-blocks"},
			[]interface{}{
				testConfig, genFile, testProtocol,
				"true", "Gossamer", gssmr.DefaultPruningMode,
				gssmr.DefaultRetainBlocks},
			&dot.Config{
				Global: testCfg.Global,
				Init: dot.InitConfig{
					Genesis: genFile,
				},
				Log: dot.LogConfig{
					CoreLvl:           log.Info,
					SyncLvl:           log.Info,
					NetworkLvl:        log.Info,
					RPCLvl:            log.Info,
					StateLvl:          log.Info,
					RuntimeLvl:        log.Info,
					BlockProducerLvl:  log.Info,
					FinalityGadgetLvl: log.Info,
				},
				Account: testCfg.Account,
				Core:    testCfg.Core,
				Network: dot.NetworkConfig{
					Port:              testCfg.Network.Port,
					Bootnodes:         []string{testBootnode},
					ProtocolID:        testProtocol,
					NoBootstrap:       testCfg.Network.NoBootstrap,
					NoMDNS:            testCfg.Network.NoMDNS,
					DiscoveryInterval: testCfg.Network.DiscoveryInterval,
					MinPeers:          testCfg.Network.MinPeers,
					MaxPeers:          testCfg.Network.MaxPeers,
				},
				RPC:   testCfg.RPC,
				Pprof: testCfg.Pprof,
			},
		},
	}

	for _, c := range testcases {
		c := c // bypass scopelint false positive
		t.Run(c.description, func(t *testing.T) {
			ctx, err := newTestContext(c.description, c.flags, c.values)
			require.NoError(t, err)

			err = exportAction(ctx)
			require.NoError(t, err)

			config := ctx.GlobalString(ConfigFlag.Name)

			cfg := new(ctoml.Config)
			err = loadConfig(cfg, config)
			require.NoError(t, err)

			require.Equal(t, dotConfigToToml(c.expected), cfg)
		})
	}
}
