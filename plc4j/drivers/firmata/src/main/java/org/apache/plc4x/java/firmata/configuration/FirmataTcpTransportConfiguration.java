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

package org.apache.plc4x.java.firmata.configuration;

import org.apache.plc4x.java.transport.tcp.config.TcpTransportConfiguration;

/**
 * Firmata over TCP. Used by ESP32-style "Firmata-over-WiFi" sketches
 * (StandardFirmataWiFi, ConfigurableFirmata with EthernetClientStream),
 * and as a way to talk to socat-bridged virtual UART simulators where the
 * underlying serial transport's path validation can't open the bridge.
 *
 * <p>Default port is 3030 — the value used by StandardFirmataWiFi.ino /
 * StandardFirmataEthernet.ino. Users can override on the URL.</p>
 */
public class FirmataTcpTransportConfiguration extends TcpTransportConfiguration {

    @Override
    public int getDefaultPort() {
        return 3030;
    }

}
