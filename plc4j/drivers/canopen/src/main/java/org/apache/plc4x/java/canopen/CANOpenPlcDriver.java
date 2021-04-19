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
package org.apache.plc4x.java.canopen;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.canopen.transport.CANOpenFrame;
import org.apache.plc4x.java.canopen.transport.CANTransport;
import org.apache.plc4x.java.canopen.transport.socketcan.io.CANOpenSocketCANFrameIO;
import org.apache.plc4x.java.canopen.configuration.CANOpenConfiguration;
import org.apache.plc4x.java.canopen.context.CANOpenDriverContext;
import org.apache.plc4x.java.canopen.field.CANOpenFieldHandler;
import org.apache.plc4x.java.canopen.protocol.CANOpenProtocolLogic;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.CustomProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleFieldOptimizer;
import org.apache.plc4x.java.spi.transport.Transport;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.apache.plc4x.java.spi.values.PlcList;

import java.util.function.ToIntFunction;

/**
 */
public class CANOpenPlcDriver extends GeneratedDriverBase<CANOpenFrame> {

    @Override
    public String getProtocolCode() {
        return "canopen";
    }

    @Override
    public String getProtocolName() {
        return "CAN open";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
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
    protected ProtocolStackConfigurer<CANOpenFrame> getStackConfigurer() {
        throw new PlcRuntimeException("CANopen driver requires access to transport layer.");
    }

    @Override
    protected ProtocolStackConfigurer<CANOpenFrame> getStackConfigurer(Transport transport) {
        if (!(transport instanceof CANTransport)) {
            throw new RuntimeException();
        }

        CANTransport transport1 = (CANTransport) transport;

        return CustomProtocolStackConfigurer.builder(transport1.getMessageType(), transport1::getMessageIO)
            //.withTransport(transport)
            .withProtocol((cfg)-> new CANOpenProtocolLogic())
            .withDriverContext((cfg) -> new CANOpenDriverContext())
            .withPacketSizeEstimator(configuration1 ->  (ByteBuf byteBuf) -> {
                if (byteBuf.readableBytes() >= 5) {
                    return 16; // socketcan transport always returns 16 bytes padded with zeros;
                }
                return -1; //discard
            })
            .littleEndian()
            .build();
    }

}
