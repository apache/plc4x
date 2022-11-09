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
package org.apache.plc4x.java.can.generic;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.can.adapter.CANDriverAdapter;
import org.apache.plc4x.java.can.generic.configuration.GenericCANConfiguration;
import org.apache.plc4x.java.can.generic.context.GenericCANDriverContext;
import org.apache.plc4x.java.can.generic.tag.GenericCANTagHandler;
import org.apache.plc4x.java.can.generic.protocol.GenericCANProtocolLogic;
import org.apache.plc4x.java.can.generic.transport.GenericCANFrameDataHandler;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.connection.CustomProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.transport.Transport;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.transport.can.CANTransport;

/**
 * A generic purpose CAN driver which is able to work with any compatible CAN transport.
 *
 * Main role of this driver is provisioning of quick and easy way to create user specific CAN bus applications.
 */
public class GenericCANDriver extends GeneratedDriverBase<Message> {

    @Override
    public String getProtocolCode() {
        return "genericcan";
    }

    @Override
    public String getProtocolName() {
        return "Generic CAN";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return GenericCANConfiguration.class;
    }

    @Override
    protected boolean canRead() {
        return false;
    }

    @Override
    protected boolean canSubscribe() {
        return true;
    }

    @Override
    protected boolean canWrite() {
        return true;
    }

    @Override
    protected String getDefaultTransport() {
        return "socketcan";
    }

    @Override
    protected GenericCANTagHandler getTagHandler() {
        return new GenericCANTagHandler();
    }

    @Override
    protected org.apache.plc4x.java.api.value.PlcValueHandler getValueHandler() {
        return new PlcValueHandler();
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
    protected BaseOptimizer getOptimizer() {
        return null;
    }

    @Override
    protected ProtocolStackConfigurer<Message> getStackConfigurer() {
        throw new PlcRuntimeException("Generic CAN driver requires access to transport layer.");
    }

    @Override
    protected ProtocolStackConfigurer<Message> getStackConfigurer(Transport transport) {
        if (!(transport instanceof CANTransport)) {
            throw new PlcRuntimeException("Generic CAN Driver requires CAN Transport instance");
        }

        CANTransport<Message> canTransport = (CANTransport<Message>) transport;
        return CustomProtocolStackConfigurer.builder(canTransport.getMessageType(), canTransport::getMessageInput)
            .withProtocol(cfg -> {
                GenericCANProtocolLogic protocolLogic = new GenericCANProtocolLogic();
                ConfigurationFactory.configure(cfg, protocolLogic);
                return new CANDriverAdapter<>(protocolLogic,
                    canTransport.getMessageType(), canTransport.adapter(),
                    new GenericCANFrameDataHandler(canTransport::getTransportFrameBuilder)
                );
            })
            .withDriverContext(cfg -> new GenericCANDriverContext())
            .withPacketSizeEstimator(cfg -> canTransport.getEstimator())
            .littleEndian()
            .build();
    }

}
