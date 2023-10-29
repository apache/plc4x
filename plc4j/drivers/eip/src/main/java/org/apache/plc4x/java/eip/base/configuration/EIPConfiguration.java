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
package org.apache.plc4x.java.eip.base.configuration;

import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ComplexConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.transport.TransportConfiguration;
import org.apache.plc4x.java.spi.transport.TransportConfigurationProvider;

public class EIPConfiguration implements Configuration, TransportConfigurationProvider {

    @ConfigurationParameter
    private int backplane = 1;

    @ConfigurationParameter
    private int slot = 0;

    @ConfigurationParameter
    private boolean bigEndian = true;

    @ComplexConfigurationParameter(prefix = "tcp", defaultOverrides = {}, requiredOverrides = {})
    private EipTcpTransportConfiguration tcpTransportConfiguration;

    public int getBackplane() {
        return backplane;
    }

    public void setBackplane(int backpane) {
        this.backplane = backpane;
    }

    public int getSlot() {
        return slot;
    }

    public void setSlot(int slot) {
        this.slot = slot;
    }

    public ByteOrder getByteOrder() {
        return this.bigEndian ? ByteOrder.BIG_ENDIAN : ByteOrder.LITTLE_ENDIAN;
    }

    public void setByteOrder(ByteOrder byteOrder) {
        this.bigEndian = byteOrder == ByteOrder.BIG_ENDIAN;
    }

    public EipTcpTransportConfiguration getTcpTransportConfiguration() {
        return tcpTransportConfiguration;
    }

    public void setTcpTransportConfiguration(EipTcpTransportConfiguration tcpTransportConfiguration) {
        this.tcpTransportConfiguration = tcpTransportConfiguration;
    }

    @Override
    public TransportConfiguration getTransportConfiguration(String transportCode) {
        switch (transportCode) {
            case "tcp":
                return tcpTransportConfiguration;
        }
        return null;
    }

}
