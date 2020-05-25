/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.amsads.configuration;

import org.apache.plc4x.java.amsads.AMSADSPlcDriver;
import org.apache.plc4x.java.amsads.readwrite.AmsNetId;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.ConfigurationParameterConverter;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ParameterConverter;
import org.apache.plc4x.java.spi.configuration.annotations.Required;
import org.apache.plc4x.java.transport.serial.SerialTransportConfiguration;
import org.apache.plc4x.java.transport.tcp.TcpTransportConfiguration;

public class AdsConfiguration implements Configuration, TcpTransportConfiguration, SerialTransportConfiguration {

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

    @Override
    public int getDefaultPort() {
        return AMSADSPlcDriver.TCP_PORT;
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
            return AMSADSPlcDriver.AmsNetIdOf(value);
        }
    }
}
