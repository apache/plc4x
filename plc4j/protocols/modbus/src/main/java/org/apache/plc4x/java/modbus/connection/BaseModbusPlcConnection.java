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
package org.apache.plc4x.java.modbus.connection;

import io.netty.channel.ChannelFuture;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.connection.AbstractPlcConnection;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.modbus.model.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;

public abstract class BaseModbusPlcConnection extends AbstractPlcConnection implements PlcReader, PlcWriter {

    private static final Logger logger = LoggerFactory.getLogger(BaseModbusPlcConnection.class);

    protected BaseModbusPlcConnection(ChannelFactory channelFactory, String params) {
        super(channelFactory);

        if (!StringUtils.isEmpty(params)) {
            for (String param : params.split("&")) {
                String[] paramElements = param.split("=");
                String paramName = paramElements[0];
                if (paramElements.length == 2) {
                    String paramValue = paramElements[1];
                    switch (paramName) {
                        default:
                            logger.debug("Unknown parameter {} with value {}", paramName, paramValue);
                    }
                } else {
                    logger.debug("Unknown no-value parameter {}", paramName);
                }
            }
        }
    }

    @Override
    public PlcField prepareField(String fieldString) throws PlcInvalidFieldException {
        if (MaskWriteRegisterModbusField.ADDRESS_PATTERN.matcher(fieldString).matches()) {
            return MaskWriteRegisterModbusField.of(fieldString);
        } else if (ReadDiscreteInputsModbusField.ADDRESS_PATTERN.matcher(fieldString).matches()) {
            return ReadDiscreteInputsModbusField.of(fieldString);
        } else if (ReadHoldingRegistersModbusField.ADDRESS_PATTERN.matcher(fieldString).matches()) {
            return ReadHoldingRegistersModbusField.of(fieldString);
        } else if (ReadInputRegistersModbusField.ADDRESS_PATTERN.matcher(fieldString).matches()) {
            return ReadInputRegistersModbusField.of(fieldString);
        } else if (CoilModbusField.ADDRESS_PATTERN.matcher(fieldString).matches()) {
            return CoilModbusField.of(fieldString);
        } else if (RegisterModbusField.ADDRESS_PATTERN.matcher(fieldString).matches()) {
            return RegisterModbusField.of(fieldString);
        }
        throw new PlcInvalidFieldException(fieldString);
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> readFuture = new CompletableFuture<>();
        ChannelFuture channelFuture = channel.writeAndFlush(new PlcRequestContainer<>(readRequest, readFuture));
        channelFuture.addListener(future -> {
            if (!future.isSuccess()) {
                readFuture.completeExceptionally(future.cause());
            }
        });
        return readFuture;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> writeFuture = new CompletableFuture<>();
        ChannelFuture channelFuture = channel.writeAndFlush(new PlcRequestContainer<>(writeRequest, writeFuture));
        channelFuture.addListener(future -> {
            if (!future.isSuccess()) {
                writeFuture.completeExceptionally(future.cause());
            }
        });
        return writeFuture;
    }
}
