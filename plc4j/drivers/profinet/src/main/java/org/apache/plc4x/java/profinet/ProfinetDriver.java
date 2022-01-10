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
package org.apache.plc4x.java.profinet;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.metadata.PlcDriverMetadata;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.discovery.ProfinetPlcDiscoverer;
import org.apache.plc4x.java.profinet.field.ProfinetField;
import org.apache.plc4x.java.profinet.field.ProfinetFieldHandler;
import org.apache.plc4x.java.profinet.protocol.ProfinetProtocolLogic;
import org.apache.plc4x.java.profinet.readwrite.Ethernet_Frame;
import org.apache.plc4x.java.profinet.readwrite.io.Ethernet_FrameIO;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryRequest;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleFieldOptimizer;

import java.util.function.ToIntFunction;

public class ProfinetDriver extends GeneratedDriverBase<Ethernet_Frame> {

    public static final String DRIVER_CODE = "profinet";

    @Override
    public String getProtocolCode() {
        return DRIVER_CODE;
    }

    @Override
    public String getProtocolName() {
        return "Profinet";
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
        return new DefaultPlcDiscoveryRequest.Builder(new ProfinetPlcDiscoverer());
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return ProfinetConfiguration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "raw";
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
    protected ProfinetFieldHandler getFieldHandler() {
        return new ProfinetFieldHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new IEC61131ValueHandler();
    }

    @Override
    protected ProtocolStackConfigurer<Ethernet_Frame> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(Ethernet_Frame.class, Ethernet_FrameIO.class)
            .withProtocol(ProfinetProtocolLogic.class)
            .withDriverContext(ProfinetDriverContext.class)
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
    public ProfinetField prepareField(String query){
        return ProfinetField.of(query);
    }

}
