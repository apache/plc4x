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
package org.apache.plc4x.java.modbus.tcp;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.metadata.PlcDriverMetadata;
import org.apache.plc4x.java.modbus.readwrite.DriverType;
import org.apache.plc4x.java.modbus.tcp.config.ModbusTcpConfiguration;
import org.apache.plc4x.java.modbus.tcp.discovery.ModbusPlcDiscoverer;
import org.apache.plc4x.java.modbus.base.field.ModbusField;
import org.apache.plc4x.java.modbus.base.field.ModbusFieldHandler;
import org.apache.plc4x.java.modbus.readwrite.ModbusTcpADU;
import org.apache.plc4x.java.modbus.tcp.protocol.ModbusTcpProtocolLogic;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryRequest;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleFieldOptimizer;

import java.util.function.ToIntFunction;

public class ModbusTcpDriver extends GeneratedDriverBase<ModbusTcpADU> {

    @Override
    public String getProtocolCode() {
        return "modbus-tcp";
    }

    @Override
    public String getProtocolName() {
        return "Modbus TCP";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return ModbusTcpConfiguration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "tcp";
    }

    @Override
    public PlcDriverMetadata getMetadata() {
        return new PlcDriverMetadata() {
            @Override
            public boolean canDiscover() {
                return true;
            }
        };
    }

    @Override
    public PlcDiscoveryRequest.Builder discoveryRequestBuilder() {
        return new DefaultPlcDiscoveryRequest.Builder(new ModbusPlcDiscoverer());
    }

    /**
     * Modbus doesn't have a login procedure, so there is no need to wait for a login to finish.
     * @return false
     */
    @Override
    protected boolean awaitSetupComplete() {
        return false;
    }

    /**
     * This protocol doesn't have a disconnect procedure, so there is no need to wait for a login to finish.
     * @return false
     */
    @Override
    protected boolean awaitDisconnectComplete() {
        return false;
    }

    @Override
    protected boolean canRead() {
        return true;
    }

    @Override
    protected boolean canWrite() {
        return true;
    }

    @Override
    protected BaseOptimizer getOptimizer() {
        return new SingleFieldOptimizer();
    }

    @Override
    protected ModbusFieldHandler getFieldHandler() {
        return new ModbusFieldHandler();
    }

    @Override
    protected org.apache.plc4x.java.api.value.PlcValueHandler getValueHandler() {
        return new PlcValueHandler();
    }

    @Override
    protected ProtocolStackConfigurer<ModbusTcpADU> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(ModbusTcpADU.class,
                (io, args) -> (ModbusTcpADU) ModbusTcpADU.staticParse(io, args))
            .withProtocol(ModbusTcpProtocolLogic.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            // Every incoming message is to be treated as a response.
            .withParserArgs(DriverType.MODBUS_TCP, true)
            .build();
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 6) {
                return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 4) + 6;
            }
            return -1;
        }
    }

    @Override
    public ModbusField prepareField(String query){
        return ModbusField.of(query);
    }

}
