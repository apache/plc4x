/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.modbus;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.modbus.config.ModbusConfiguration;
import org.apache.plc4x.java.modbus.field.ModbusField;
import org.apache.plc4x.java.modbus.field.ModbusFieldHandler;
import org.apache.plc4x.java.modbus.protocol.ModbusProtocolLogic;
import org.apache.plc4x.java.modbus.readwrite.ModbusTcpADU;
import org.apache.plc4x.java.modbus.readwrite.io.ModbusTcpADUIO;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleFieldOptimizer;

import java.util.function.ToIntFunction;

public class ModbusDriver extends GeneratedDriverBase<ModbusTcpADU> {

    @Override
    public String getProtocolCode() {
        return "modbus";
    }

    @Override
    public String getProtocolName() {
        return "Modbus";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return ModbusConfiguration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "tcp";
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
    protected PlcValueHandler getValueHandler() {
        return new IEC61131ValueHandler();
    }

    @Override
    protected ProtocolStackConfigurer<ModbusTcpADU> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(ModbusTcpADU.class, ModbusTcpADUIO.class)
            .withProtocol(ModbusProtocolLogic.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            // Every incoming message is to be treated as a response.
            .withParserArgs(true)
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
