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
import org.apache.plc4x.java.modbus.readwrite.ModbusADU;
import org.apache.plc4x.java.modbus.rtu.context.ModbusRtuContext;
import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.spi.configuration.PlcTransportConfiguration;
import org.apache.plc4x.java.modbus.base.tag.ModbusTag;
import org.apache.plc4x.java.modbus.readwrite.DriverType;
import org.apache.plc4x.java.modbus.readwrite.ModbusRtuADU;
import org.apache.plc4x.java.modbus.rtu.config.ModbusRtuConfiguration;
import org.apache.plc4x.java.modbus.rtu.protocol.ModbusRtuProtocolLogic;
import org.apache.plc4x.java.modbus.tcp.config.ModbusTcpTransportConfiguration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleTagOptimizer;

import java.util.Arrays;
import java.util.List;
import java.util.Optional;
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
    protected Class<? extends PlcConnectionConfiguration> getConfigurationClass() {
        return ModbusRtuConfiguration.class;
    }

    @Override
    protected Optional<Class<? extends PlcTransportConfiguration>> getTransportConfigurationClass(String transportCode) {
        switch (transportCode) {
            case "tcp":
                return Optional.of(ModbusTcpTransportConfiguration.class);
        }
        return Optional.empty();
    }

    @Override
    protected Optional<String> getDefaultTransportCode() {
        return Optional.of("serial");
    }

    @Override
    protected List<String> getSupportedTransportCodes() {
        return Arrays.asList("tcp", "serial");
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
    protected boolean canPing() {
        return true;
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
    protected ProtocolStackConfigurer<ModbusRtuADU> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(ModbusRtuADU.class, io -> (ModbusRtuADU) ModbusRtuADU.staticParse(io, DriverType.MODBUS_RTU, true))
            .withProtocol(ModbusRtuProtocolLogic.class)
            .withDriverContext(ModbusRtuContext.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .build();
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            // A Modbus RTU packet has the absolute minimum size of 4 (if it has absolutely no payload)
            if (byteBuf.readableBytes() >= 4) {
                // Fetch what's currently in the buffer
                byte[] buf = new byte[byteBuf.readableBytes()];
                byteBuf.getBytes(byteBuf.readerIndex(), buf);
                ReadBufferByteBased reader = new ReadBufferByteBased(buf);

                // Try to parse the buffer content.
                try {
                    ModbusADU modbusADU = ModbusRtuADU.staticParse(reader, DriverType.MODBUS_RTU, true);

                    // Theoretically, the buffer could contain more than one message.
                    return modbusADU.getLengthInBytes();
                } catch (ParseException e) {
                    return -1;
                }
                // If we're getting this error, manually compact the buffer.
                // Hopefully now there will be enough space for another attempt.
                catch (ArrayIndexOutOfBoundsException e) {
                    byteBuf.discardReadBytes();
                    return -1;
                }
            }
            return -1;
        }
    }

    @Override
    public ModbusTag prepareTag(String tagAddress){
        return ModbusTag.of(tagAddress);
    }

}