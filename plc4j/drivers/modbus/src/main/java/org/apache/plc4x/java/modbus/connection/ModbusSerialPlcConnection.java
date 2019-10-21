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

import io.netty.channel.ChannelHandler;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.serial.connection.connection.SerialChannelFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;

public class ModbusSerialPlcConnection extends BaseModbusPlcConnection {

    private static final Logger logger = LoggerFactory.getLogger(ModbusSerialPlcConnection.class);

    private ModbusSerialPlcConnection(String port, String params) {
        super(new SerialChannelFactory(port), params);
        logger.info("Configured ModbusSerialPlcConnection with: serial-port {}", port);
    }

    public static ModbusSerialPlcConnection of(String serialPort, String params) {
        return new ModbusSerialPlcConnection(serialPort, params);
    }

    @Override
    protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
        throw new NotImplementedException("Not implemented yet");
    }

}
