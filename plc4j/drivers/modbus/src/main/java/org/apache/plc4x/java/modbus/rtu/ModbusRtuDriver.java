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
package org.apache.plc4x.java.modbus.rtu;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.modbus.base.tag.ModbusTag;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagHandler;
import org.apache.plc4x.java.modbus.readwrite.DriverType;
import org.apache.plc4x.java.modbus.readwrite.ModbusRtuADU;
import org.apache.plc4x.java.modbus.rtu.config.ModbusRtuConfiguration;
import org.apache.plc4x.java.modbus.rtu.protocol.ModbusRtuProtocolLogic;
import org.apache.plc4x.java.spi.configuration.BaseConfiguration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleTagOptimizer;
import org.apache.plc4x.java.spi.values.PlcValueHandler;

import java.util.function.ToIntFunction;

public class ModbusRtuDriver extends GeneratedDriverBase<ModbusRtuADU> {

    @Override
    public String getProtocolCode() {
        return "modbus-rtu";
    }

    @Override
    public String getProtocolName() {
        return "Modbus RTU";
    }

    @Override
    protected Class<? extends BaseConfiguration> getConfigurationType() {
        return ModbusRtuConfiguration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "serial";
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
        return new SingleTagOptimizer();
    }

    @Override
    protected ModbusTagHandler getTagHandler() {
        return new ModbusTagHandler();
    }

    @Override
    protected org.apache.plc4x.java.api.value.PlcValueHandler getValueHandler() {
        return new PlcValueHandler();
    }

    @Override
    protected ProtocolStackConfigurer<ModbusRtuADU> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(ModbusRtuADU.class,
                (io, args) -> (ModbusRtuADU) ModbusRtuADU.staticParse(io, args))
            .withProtocol(ModbusRtuProtocolLogic.class)
            .withPacketSizeEstimator(ModbusRtuDriver.ByteLengthEstimator.class)
            // Every incoming message is to be treated as a response.
            .withParserArgs(DriverType.MODBUS_RTU, true)
            .build();
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 1) {
                return byteBuf.readableBytes();
            }
            return -1;
        }
    }

    @Override
    public ModbusTag prepareTag(String tagAddress){
        return ModbusTag.of(tagAddress);
    }

}