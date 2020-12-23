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
package org.apache.plc4x.java.opcua;

import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.opcua.connection.OpcuaConnectionFactory;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.opcua.protocol.*;
import org.apache.plc4x.java.opcua.config.*;
import org.apache.plc4x.java.opcua.readwrite.*;
import org.apache.plc4x.java.opcua.readwrite.io.*;
import org.apache.plc4x.java.opcua.readwrite.types.*;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleFieldOptimizer;
import io.netty.buffer.ByteBuf;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import java.util.function.ToIntFunction;

/**
 * Implementation of the OPC UA protocol, based on:
 * - Eclipse Milo (https://github.com/eclipse/milo)
 *
 * Created by Matthias Milan Strljic on 10.05.2019
 */
public class OpcuaPlcDriver extends GeneratedDriverBase<OpcuaAPU> {


    public static final Pattern INET_ADDRESS_PATTERN = Pattern.compile("(:(?<transport>tcp))?://(?<host>[\\w.-]+)(:(?<port>\\d*))?");
    public static final Pattern OPCUA_URI_PARAM_PATTERN = Pattern.compile("(?<param>[(\\?|\\&)([^=]+)\\=([^&]+)]+)?"); //later used for regex filtering of the params
    public static final Pattern OPCUA_URI_PATTERN = Pattern.compile("^opcua" + INET_ADDRESS_PATTERN + "(?<params>[\\w/=?&]+)?");
    private static final int requestTimeout = 10000;
    private OpcuaConnectionFactory opcuaConnectionFactory = new OpcuaConnectionFactory();


    @Override
    public String getProtocolCode() {
        return "opcua";
    }

    @Override
    public String getProtocolName() {
        return "Opcua";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return OpcuaConfiguration.class;
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
        return new SingleFieldOptimizer();
    }

    @Override
    protected OpcuaPlcFieldHandler getFieldHandler() {
        return new OpcuaPlcFieldHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new IEC61131ValueHandler();
    }

    @Override
    protected ProtocolStackConfigurer<OpcuaAPU> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(OpcuaAPU.class, OpcuaAPUIO.class)
            .withProtocol(OpcuaProtocolLogic.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            // Every incoming message is to be treated as a response.
            .withParserArgs(true)
            .littleEndian()
            .build();
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 6) {
                return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 256) + 6;
            }
            return -1;
        }
    }

    @Override
    public OpcuaField prepareField(String query){
        return OpcuaField.of(query);
    }

}
