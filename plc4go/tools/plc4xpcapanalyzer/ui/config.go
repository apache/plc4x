/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package ui

import (
	cliConfig "github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"time"
)

var plc4xpcapanalyzerConfigDir string
var configFile string
var config = Config{
	MaxConsoleLines: 500,
	MaxOutputLines:  500,
	CliConfigs:      allCliConfigsInstances,
}

type Config struct {
	HostIp  string `yaml:"host_ip"`
	History struct {
		Last10Files    []string `yaml:"last_hosts"`
		Last10Commands []string `yaml:"last_commands"`
	}
	AutoRegisterDrivers []string      `yaml:"auto_register_driver"`
	LastUpdated         time.Time     `yaml:"last_updated"`
	LogLevel            string        `yaml:"log_level"`
	MaxConsoleLines     int           `yaml:"max_console_lines"`
	MaxOutputLines      int           `yaml:"max_output_lines"`
	CliConfigs          allCliConfigs `yaml:"cli_configs"`
}

type allCliConfigs struct {
	RootConfig    *cliConfig.RootConfig
	AnalyzeConfig *cliConfig.AnalyzeConfig
	ExtractConfig *cliConfig.ExtractConfig
	BacnetConfig  *cliConfig.BacnetConfig
	CBusConfig    *cliConfig.CBusConfig
	PcapConfig    *cliConfig.PcapConfig
}

var allCliConfigsInstances = allCliConfigs{
	&cliConfig.RootConfigInstance,
	&cliConfig.AnalyzeConfigInstance,
	&cliConfig.ExtractConfigInstance,
	&cliConfig.BacnetConfigInstance,
	&cliConfig.CBusConfigInstance,
	&cliConfig.PcapConfigInstance,
}

func init() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	plc4xpcapanalyzerConfigDir = path.Join(userConfigDir, "plc4xpcapanalyzer")
	if _, err := os.Stat(plc4xpcapanalyzerConfigDir); os.IsNotExist(err) {
		err := os.Mkdir(plc4xpcapanalyzerConfigDir, os.ModeDir|os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	configFile = path.Join(plc4xpcapanalyzerConfigDir, "config.yml")
}

func LoadConfig() {
	f, err := os.Open(configFile)
	if err != nil {
		log.Info().Err(err).Msg("No config file found")
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error closing config file")
		}
	}(f)

	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&config); err != nil {
		log.Warn().Err(err).Msg("Can't decode config file")
		return
	}
}

func saveConfig() {
	config.LastUpdated = time.Now()
	f, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Warn().Err(err).Msg("Can't save config file")
		return
	}
	encoder := yaml.NewEncoder(f)
	defer func(encoder *yaml.Encoder) {
		err := encoder.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error closing config file")
		}
	}(encoder)
	if err := encoder.Encode(config); err != nil {
		log.Warn().Err(err).Msg("Can't encode config file")
		panic(err)
	}
}

func addRecentFilesEntry(pcapFile string) {
	existingIndex := -1
	for i, lastPcapFile := range config.History.Last10Files {
		if lastPcapFile == pcapFile {
			existingIndex = i
			break
		}
	}
	if existingIndex >= 0 {
		config.History.Last10Files = append(config.History.Last10Files[:existingIndex], config.History.Last10Files[existingIndex+1:]...)
	}
	if len(config.History.Last10Files) >= 10 {
		config.History.Last10Files = config.History.Last10Files[1:]
	}
	config.History.Last10Files = append(config.History.Last10Files, pcapFile)
}

func addCommandHistoryEntry(command string) {
	switch command {
	case "clear":
		return
	case "history":
		return
	}
	existingIndex := -1
	for i, lastCommand := range config.History.Last10Commands {
		if lastCommand == command {
			existingIndex = i
			break
		}
	}
	if existingIndex >= 0 {
		config.History.Last10Commands = append(config.History.Last10Commands[:existingIndex], config.History.Last10Commands[existingIndex+1:]...)
	}
	if len(config.History.Last10Commands) >= 10 {
		config.History.Last10Commands = config.History.Last10Commands[1:]
	}
	config.History.Last10Commands = append(config.History.Last10Commands, command)
}

func setLevel(level zerolog.Level) {
	config.LogLevel = level.String()
}

func enableAutoRegister(driver string) error {
	if err := validateDriverParam(driver); err != nil {
		return err
	}
	for _, autoRegisterDriver := range config.AutoRegisterDrivers {
		if autoRegisterDriver == driver {
			return errors.Errorf("%s already registered for auto register", driver)
		}
	}
	config.AutoRegisterDrivers = append(config.AutoRegisterDrivers, driver)
	log.Info().Str("driver", driver).Msg("Auto register enabled")
	return nil
}

func disableAutoRegister(driver string) error {
	if err := validateDriverParam(driver); err != nil {
		return err
	}
	index := -1
	for i, autoRegisterDriver := range config.AutoRegisterDrivers {
		if autoRegisterDriver == driver {
			index = i
			break
		}
	}
	if index < 0 {
		return errors.Errorf("%s not registered for auto register", driver)
	}
	config.AutoRegisterDrivers = append(config.AutoRegisterDrivers[:index], config.AutoRegisterDrivers[index+1:]...)
	log.Info().Str("driver", driver).Msg("Auto register disabled")
	return nil
}
