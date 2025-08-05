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
package org.apache.plc4x.java.api.metadata;

import java.util.List;
import java.util.Optional;

/**
 * Information about driver capabilities.
 */
public interface PlcDriverMetadata {

    /**
     * @return The transport code of the default transport, if the driver supports this.
     */
    Optional<String> getDefaultTransportCode();

    /**
     * @return A list of all actively supported transports.
     */
    List<String> getSupportedTransportCodes();

    /**
     * @return Get the configuration options for the current driver.
     */
    Optional<OptionMetadata> getProtocolConfigurationOptionMetadata();

    /**
     * @param transportCode transport code for the transport we want to get the transport options for.
     * @return Get the transport options for the current driver and the given transport type.
     */
    Optional<OptionMetadata> getTransportConfigurationOptionMetadata(String transportCode);

    /**
     * @return true, if the current driver supports discovery.
     */
    boolean isDiscoverySupported();

}
