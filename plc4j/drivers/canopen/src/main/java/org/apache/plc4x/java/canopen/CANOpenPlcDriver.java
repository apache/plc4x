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
package org.apache.plc4x.java.canopen;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.can.adapter.CANDriverAdapter;
import org.apache.plc4x.java.canopen.configuration.CANOpenConfiguration;
import org.apache.plc4x.java.canopen.context.CANOpenDriverContext;
import org.apache.plc4x.java.canopen.field.CANOpenFieldHandler;
import org.apache.plc4x.java.canopen.protocol.CANOpenProtocolLogic;
import org.apache.plc4x.java.canopen.transport.CANOpenFrameDataHandler;
import org.apache.plc4x.java.spi.configuration.BaseConfiguration;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.connection.CustomProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleFieldOptimizer;
import org.apache.plc4x.java.spi.transport.Transport;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.transport.can.CANTransport;

/**
 */
public class CANOpenPlcDriver extends GeneratedDriverBase<Message> {

    @Override
    public String getProtocolCode() {
        return "canopen";
    }

    @Override
    public String getProtocolName() {
        return "CAN open";
    }

    @Override
    protected Class<? extends BaseConfiguration> getConfigurationType() {
        return CANOpenConfiguration.class;
    }

    @Override
    protected boolean canRead() {
        return true;
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
    protected CANOpenFieldHandler getFieldHandler() {
        return new CANOpenFieldHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new IEC61131ValueHandler() {
            @Override
            public PlcValue newPlcValue(PlcField field, Object[] values) {
                if (values[0] instanceof PlcList) {
                    return (PlcList) values[0];
                }
                return super.newPlcValue(field, values);
            }
        };
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
        return new SingleFieldOptimizer();
    }

    @Override
    protected ProtocolStackConfigurer<Message> getStackConfigurer() {
        throw new PlcRuntimeException("CANopen driver requires access to transport layer.");
    }

    @Override
    protected ProtocolStackConfigurer<Message> getStackConfigurer(Transport transport) {
        if (!(transport instanceof CANTransport)) {
            throw new PlcRuntimeException("CANopen driver requires a CAN transport instance");
        }

        final CANTransport<Message> canTransport = (CANTransport<Message>) transport;
        return CustomProtocolStackConfigurer.builder(canTransport.getMessageType(), canTransport::getMessageInput)
            .withProtocol(cfg -> {
                CANOpenProtocolLogic protocolLogic = new CANOpenProtocolLogic();
                ConfigurationFactory.configure(cfg, protocolLogic);
                return new CANDriverAdapter<>(protocolLogic,
                    canTransport.getMessageType(), canTransport.adapter(),
                    new CANOpenFrameDataHandler(canTransport::getTransportFrameBuilder)
                );
            })
            .withDriverContext(cfg -> new CANOpenDriverContext())
            .withPacketSizeEstimator(configuration1 -> canTransport.getEstimator())
            .littleEndian()
            .build();
    }

}
