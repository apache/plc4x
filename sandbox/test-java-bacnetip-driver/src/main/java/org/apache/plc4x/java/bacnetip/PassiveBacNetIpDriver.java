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
package org.apache.plc4x.java.bacnetip;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.bacnetip.connection.PassiveBacNetIpPlcConnection;
import org.apache.plc4x.java.bacnetip.protocol.HelloWorldProtocol;
import org.apache.plc4x.java.spi.PlcDriver;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketIpAddress;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

import static org.apache.plc4x.java.utils.rawsockets.netty.RawSocketAddress.ALL_PROTOCOLS;

public class PassiveBacNetIpDriver implements PlcDriver {

    public static final int BACNET_IP_PORT = 47808;

    private static final Pattern PASSIVE_BACNET_IP_URI_PATTERN =
        Pattern.compile("^bachnet-ip-passive://(?<networkDevice>.*)(?<params>\\?.*)?");

    @Override
    public String getProtocolCode() {
        return "bacnet-ip";
    }

    @Override
    public String getProtocolName() {
        return "BACnet/IP";
    }

    @Override
    public PlcConnection connect(String url) throws PlcConnectionException {
        Matcher matcher = PASSIVE_BACNET_IP_URI_PATTERN.matcher(url);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection url doesn't match the format 'bacnet-ip://{host|ip}'");
        }
        String networkDevice = matcher.group("networkDevice");

        String params = matcher.group("params") != null ? matcher.group("params").substring(1) : null;

        try {
            RawSocketIpAddress rawSocketAddress = new RawSocketIpAddress(
                networkDevice, ALL_PROTOCOLS, null, BACNET_IP_PORT);
            return new PassiveBacNetIpPlcConnection(rawSocketAddress, params, new HelloWorldProtocol());
        } catch (Exception e) {
            throw new PlcConnectionException("Error connecting to host", e);
        }
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("BACnet/IP connections don't support authentication.");
    }

}
