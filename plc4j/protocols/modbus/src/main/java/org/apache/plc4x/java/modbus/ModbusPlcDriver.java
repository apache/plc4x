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
package org.apache.plc4x.java.modbus;

import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.modbus.connection.ModbusSerialPlcConnection;
import org.apache.plc4x.java.modbus.connection.ModbusTcpPlcConnection;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Implementation of the Modbus protocol, based on:
 * - Modbus Protocol (http://www.modbus.org/docs/Modbus_Application_Protocol_V1_1b3.pdf)
 */
public class ModbusPlcDriver implements PlcDriver {

    public static final Pattern INET_ADDRESS_PATTERN = Pattern.compile("tcp://(?<host>[\\w.]+)(:(?<port>\\d*))?");
    public static final Pattern SERIAL_PATTERN = Pattern.compile("serial://(?<serialDefinition>((?!/\\d).)*)");
    public static final Pattern MODBUS_URI_PATTERN = Pattern.compile("^modbus:(" + INET_ADDRESS_PATTERN + "|" + SERIAL_PATTERN + ")/?" + "(?<params>\\?.*)?");

    @Override
    public String getProtocolCode() {
        return "modbus";
    }

    @Override
    public String getProtocolName() {
        return "Modbus (TCP / Serial)";
    }

    @Override
    public PlcConnection connect(String url) throws PlcConnectionException {
        Matcher matcher = MODBUS_URI_PATTERN.matcher(url);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection url doesn't match the format 'modbus:{type}//{port|host}'");
        }

        String host = matcher.group("host");
        String serialDefinition = matcher.group("serialDefinition");
        String port = matcher.group("port");
        String params = matcher.group("params") != null ? matcher.group("params").substring(1) : null;
        if (serialDefinition == null) {
            String hostName = matcher.group("host");
            try {
                InetAddress inetAddress = InetAddress.getByName(host);
                if (port == null) {
                    return new ModbusTcpPlcConnection(inetAddress, params);
                } else {
                    return new ModbusTcpPlcConnection(inetAddress, Integer.valueOf(port), params);
                }
            } catch (UnknownHostException e) {
                throw new PlcConnectionException("Unknown host" + hostName, e);
            }
        } else {
            return new ModbusSerialPlcConnection(serialDefinition, params);
        }
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("Modbus connections don't support authentication.");
    }

}
