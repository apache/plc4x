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

package main

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"time"
)

var plc4xBrowserConfigDir string
var configFile string
var config Config

type Config struct {
	History struct {
		Last10Hosts []string `yaml:"last_hosts"`
	}
	lastUpdated time.Time `yaml:"last_updated"`
}

func init() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	plc4xBrowserConfigDir = path.Join(userConfigDir, "plc4xbrowser")
	if _, err := os.Stat(plc4xBrowserConfigDir); os.IsNotExist(err) {
		err := os.Mkdir(plc4xBrowserConfigDir, os.ModeDir|os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	configFile = path.Join(plc4xBrowserConfigDir, "config.yml")
}

func loadConfig() {
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
	config.lastUpdated = time.Now()
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

func addHost(host string) {
	existingIndex := -1
	for i, lastHost := range config.History.Last10Hosts {
		if lastHost == host {
			existingIndex = i
			break
		}
	}
	if existingIndex >= 0 {
		config.History.Last10Hosts = append(config.History.Last10Hosts[:existingIndex], config.History.Last10Hosts[existingIndex+1:]...)
	}
	if len(config.History.Last10Hosts) >= 10 {
		config.History.Last10Hosts = config.History.Last10Hosts[1:]
	}
	config.History.Last10Hosts = append(config.History.Last10Hosts, host)
}
