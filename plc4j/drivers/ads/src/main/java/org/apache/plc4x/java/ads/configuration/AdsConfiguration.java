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
import org.apache.plc4x.java.ads.readwrite.AmsNetId;
import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.spi.configuration.ConfigurationParameterConverter;
import org.apache.plc4x.java.spi.configuration.annotations.*;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;

import java.util.regex.Pattern;
import java.util.stream.Stream;

public class AdsConfiguration implements PlcConnectionConfiguration {

    public static final Pattern AMS_NET_ID_PATTERN =
        Pattern.compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}");

    @Required
    @ConfigurationParameter("target-ams-net-id")
    @ParameterConverter(AmsNetIdConverter.class)
    @Description("AMS-Net-Id of the target. An AMS-Net-Id has the regular format of an IPv4 IP-Address, however with 6 segments instead of 4.")
    protected AmsNetId targetAmsNetId;

    @Required
    @ConfigurationParameter("target-ams-port")
    @Description("AMS port of the target.")
    protected int targetAmsPort;

    @Required
    @ConfigurationParameter("source-ams-net-id")
    @ParameterConverter(AmsNetIdConverter.class)
    @Description("AMS-Net-Id of the source. An AMS-Net-Id has the regular format of an IPv4 IP-Address, however with 6 segments instead of 4.")
    protected AmsNetId sourceAmsNetId;

    @Required
    @ConfigurationParameter("source-ams-port")
    @Description("AMS port of the source.")
    protected int sourceAmsPort;

    @ConfigurationParameter("timeout-request")
    @IntDefaultValue(4000)
    @Description("Default timeout for all types of requests.")
    protected int timeoutRequest;

    @ConfigurationParameter("load-symbol-and-data-type-tables")
    @BooleanDefaultValue(true)
    @Description("Configures, if when connecting the data-type- and symbol-table should be read. This is an optimization that can help in cases, where the PLC program is pretty large and downloading the full tables is causing problems. When disabled, symbolic addresses will manually be resolved as soon as an address is used.")
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
