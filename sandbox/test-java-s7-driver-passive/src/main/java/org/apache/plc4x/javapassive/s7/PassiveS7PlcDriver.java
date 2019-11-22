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
package org.apache.plc4x.javapassive.s7;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.spi.PlcDriver;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketIpAddress;
import org.apache.plc4x.javapassive.s7.connection.PassiveS7PlcConnection;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

import static org.apache.plc4x.java.utils.rawsockets.netty.RawSocketAddress.ALL_PROTOCOLS;

/**
 * Implementation of the S7 protocol, based on:
 * - S7 Protocol
 * - ISO Transport Protocol (Class 0) (https://tools.ietf.org/html/rfc905)
 * - ISO on TCP (https://tools.ietf.org/html/rfc1006)
 * - TCP
 */
public class PassiveS7PlcDriver implements PlcDriver {

    private static final int ISO_ON_TCP_PORT = 102;

    private static final Pattern S7_URI_PATTERN = Pattern.compile("^s7-passive://(?<networkDevice>.*)(?<params>\\?.*)?");

    @Override
    public String getProtocolCode() {
        return "s7-passive";
    }

    @Override
    public String getProtocolName() {
        return "Siemens S7 (Passive-Mode)";
    }

    @Override
    public PlcConnection connect(String url) throws PlcConnectionException {
        Matcher matcher = S7_URI_PATTERN.matcher(url);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection url doesn't match the format 's7-passive://{network-device}'");
        }
        String networkDevice = matcher.group("networkDevice");

        String params = matcher.group("params") != null ? matcher.group("params").substring(1) : null;

        try {
            RawSocketIpAddress rawSocketAddress = new RawSocketIpAddress(
                networkDevice, ALL_PROTOCOLS, null, ISO_ON_TCP_PORT);
            return new PassiveS7PlcConnection(rawSocketAddress, params);
        } catch (Exception e) {
            throw new PlcConnectionException("Error connecting to host", e);
        }
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("Basic S7 connections don't support authentication.");
    }

    public static void main(String[] args) throws Exception {
        try (PlcConnection connection = new PlcDriverManager().getConnection("s7-passive://en10")) {
            PlcReadRequest readRequest = connection.readRequestBuilder().addItem("hurz", "lalala").build();
            readRequest.execute().get();
        }
    }

}
