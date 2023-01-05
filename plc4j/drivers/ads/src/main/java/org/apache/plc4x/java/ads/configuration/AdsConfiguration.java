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
package org.apache.plc4x.java.ads.configuration;

import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.ads.readwrite.AdsConstants;
import org.apache.plc4x.java.ads.readwrite.AmsNetId;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.ConfigurationParameterConverter;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ParameterConverter;
import org.apache.plc4x.java.spi.configuration.annotations.Required;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.transport.serial.SerialTransportConfiguration;
import org.apache.plc4x.java.transport.tcp.TcpTransportConfiguration;

import java.util.regex.Pattern;
import java.util.stream.Stream;

public class AdsConfiguration implements Configuration, TcpTransportConfiguration, SerialTransportConfiguration {

    public static final Pattern AMS_NET_ID_PATTERN =
        Pattern.compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}");

    @Required
    @ConfigurationParameter
    @ParameterConverter(AmsNetIdConverter.class)
    protected AmsNetId targetAmsNetId;

    @Required
    @ConfigurationParameter
    protected int targetAmsPort;

    @Required
    @ConfigurationParameter
    @ParameterConverter(AmsNetIdConverter.class)
    protected AmsNetId sourceAmsNetId;

    @Required
    @ConfigurationParameter
    protected int sourceAmsPort;

    @ConfigurationParameter("timeout-symbolic-address-resolution")
    @IntDefaultValue(2000)
    protected int timeoutSymbolicAddressResolution;

    @ConfigurationParameter("timeout-request")
    @IntDefaultValue(4000)
    protected int timeoutRequest;

    @ConfigurationParameter("load-symbol-and-data-type-tables")
    @BooleanDefaultValue(true)
    protected boolean loadSymbolAndDataTypeTables;

    public AmsNetId getTargetAmsNetId() {
        return targetAmsNetId;
    }

    public void setTargetAmsNetId(AmsNetId targetAmsNetId) {
        this.targetAmsNetId = targetAmsNetId;
    }

    public int getTargetAmsPort() {
        return targetAmsPort;
    }

    public void setTargetAmsPort(int targetAmsPort) {
        this.targetAmsPort = targetAmsPort;
    }

    public AmsNetId getSourceAmsNetId() {
        return sourceAmsNetId;
    }

    public void setSourceAmsNetId(AmsNetId sourceAmsNetId) {
        this.sourceAmsNetId = sourceAmsNetId;
    }

    public int getSourceAmsPort() {
        return sourceAmsPort;
    }

    public void setSourceAmsPort(int sourceAmsPort) {
        this.sourceAmsPort = sourceAmsPort;
    }

    public int getTimeoutSymbolicAddressResolution() {
        return timeoutSymbolicAddressResolution;
    }

    public void setTimeoutSymbolicAddressResolution(int timeoutSymbolicAddressResolution) {
        this.timeoutSymbolicAddressResolution = timeoutSymbolicAddressResolution;
    }

    public int getTimeoutRequest() {
        return timeoutRequest;
    }

    public void setTimeoutRequest(int timeoutRequest) {
        this.timeoutRequest = timeoutRequest;
    }

    public boolean isLoadSymbolAndDataTypeTables() {
        return loadSymbolAndDataTypeTables;
    }

    public void setLoadSymbolAndDataTypeTables(boolean loadSymbolAndDataTypeTables) {
        this.loadSymbolAndDataTypeTables = loadSymbolAndDataTypeTables;
    }

    @Override
    public int getDefaultPort() {
        return AdsConstants.ADSTCPDEFAULTPORT;
    }

    @Override
    public int getBaudRate() {
        return 57600;
    }

    public static class AmsNetIdConverter implements ConfigurationParameterConverter<AmsNetId> {

        @Override
        public Class<AmsNetId> getType() {
            return AmsNetId.class;
        }

        @Override
        public AmsNetId convert(String value) {
            return amsNetIdOf(value);
        }
    }

    public static AmsNetId amsNetIdOf(String address) {
        if (!AMS_NET_ID_PATTERN.matcher(address).matches()) {
            throw new IllegalArgumentException(address + " must match " + AMS_NET_ID_PATTERN);
        }
        String[] split = address.split("\\.");
        short[] shorts = ArrayUtils.toPrimitive(Stream.of(split).map(Integer::parseInt).map(Integer::shortValue).toArray(Short[]::new));
        return new AmsNetId(shorts[0], shorts[1], shorts[2], shorts[3], shorts[4], shorts[5]);
    }

}
