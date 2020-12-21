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
import org.apache.plc4x.java.opcua.protocol.OpcuaField;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;


/**
 * Implementation of the OPC UA protocol, based on:
 * - Eclipse Milo (https://github.com/eclipse/milo)
 *
 * Created by Matthias Milan Strljic on 10.05.2019
 */
public class OpcuaPlcDriver implements PlcDriver {


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
        return "OPC UA (TCP)";
    }

    @Override
    public PlcConnection getConnection(String url) throws PlcConnectionException {
        Matcher matcher = OPCUA_URI_PATTERN.matcher(url);

        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection url doesn't match the format 'opcua:{type}//{host|port}'");
        }

        String host = matcher.group("host");
        String portString = matcher.group("port");
        Integer port = StringUtils.isNotBlank(portString) ? Integer.parseInt(portString) : null;
        String params = matcher.group("params") != null ? matcher.group("params").substring(1) : "";

        try {
            return opcuaConnectionFactory.opcuaTcpPlcConnectionOf(InetAddress.getByName(host), port, params, requestTimeout);
        } catch (UnknownHostException e) {
            throw new PlcConnectionException(e);
        }
    }

    @Override
    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("opcua does not support Auth at this state");
    }

    @Override
    public OpcuaField prepareField(String query){
        return OpcuaField.of(query);
    }

}
