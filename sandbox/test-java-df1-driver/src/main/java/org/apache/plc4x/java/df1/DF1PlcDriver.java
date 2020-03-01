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
package org.apache.plc4x.java.df1;

import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.df1.connection.SerialDf1Connection;
import org.apache.plc4x.java.spi.PlcDriver;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class DF1PlcDriver implements PlcDriver {

    // TODO there is only serial, I guess?
    public static final Pattern INET_ADDRESS_PATTERN = Pattern.compile("tcp://(?<host>[\\w.]+)(:(?<port>\\d*))?");
    public static final Pattern SERIAL_PATTERN = Pattern.compile("serial://(?<serialDefinition>/?[a-zA-Z0-9/]*)");
    public static final Pattern DF1_URI_PATTERN = Pattern.compile("^df1:(" + INET_ADDRESS_PATTERN + "|" + SERIAL_PATTERN + ")/?" + "(?<params>\\?.*)?");

    @Override
    public String getProtocolCode() {
        return "df1";
    }

    @Override
    public String getProtocolName() {
        return "Allen-Bradley DF1";
    }

    @Override
    public PlcConnection connect(String url) throws PlcConnectionException {
        Matcher matcher = DF1_URI_PATTERN.matcher(url);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection url doesn't match the format 'df1:{type}//{port|host}'");
        }

        String host = matcher.group("host");
        String serialDefinition = matcher.group("serialDefinition");
        String portString = matcher.group("port");
        Integer port = StringUtils.isNotBlank(portString) ? Integer.parseInt(portString) : null;
        String params = matcher.group("params") != null ? matcher.group("params").substring(1) : null;

        if (serialDefinition != null) {
            return new SerialDf1Connection(serialDefinition, params);
        } else {
            throw new PlcConnectionException("TCP DF1 connections not implemented yet.");
        }
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("DF1 connections doesn't support authentication.");
    }

}
