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

package org.apache.plc4x.java.tools.ui.model;

import java.util.Map;

public class Driver {

    private final String code;
    private final String name;

    private final boolean supportsDiscovery;

    private final Map<String, ConfigurationOption> configurationOptions;
    private final Map<String, Transport> transports;

    public Driver(String code, String name, boolean supportsDiscovery, Map<String, ConfigurationOption> configurationOptions, Map<String, Transport> transports) {
        this.code = code;
        this.name = name;
        this.supportsDiscovery = supportsDiscovery;
        this.configurationOptions = configurationOptions;
        this.transports = transports;
    }

    public String getCode() {
        return code;
    }

    public String getName() {
        return name;
    }

    public boolean isSupportsDiscovery() {
        return supportsDiscovery;
    }

    public Map<String, ConfigurationOption> getConfigurationOptions() {
        return configurationOptions;
    }

    public Map<String, Transport> getTransports() {
        return transports;
    }

}
