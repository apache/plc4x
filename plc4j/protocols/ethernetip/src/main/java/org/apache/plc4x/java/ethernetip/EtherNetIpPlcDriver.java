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
package org.apache.plc4x.java.ethernetip;

import org.apache.plc4x.java.spi.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.ethernetip.connection.EtherNetIpTcpPlcConnection;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Implementation of the Ethernet/IP protocol, based on the driver implementation available at:
 * https://github.com/digitalpetri/ethernet-ip/
 *
 * Spec:
 * http://read.pudn.com/downloads166/ebook/763212/EIP-CIP-V2-1.0.pdf
 */
public class EtherNetIpPlcDriver implements PlcDriver {

    private static final Pattern ETHERNETIP_URI_PATTERN = Pattern.compile("^eip://(?<host>[\\w.]+)(:(?<port>\\d*))?(?<params>\\?.*)?");

    @Override
    public String getProtocolCode() {
        return "eip";
    }

    @Override
    public String getProtocolName() {
        return "EtherNet/IP (TCP)";
    }

    @Override
    public PlcConnection connect(String url) throws PlcConnectionException {
        Matcher matcher = ETHERNETIP_URI_PATTERN.matcher(url);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection url doesn't match the format 'eip//{port|host}'");
        }

        String host = matcher.group("host");
        String port = matcher.group("port");
        String params = matcher.group("params") != null ? matcher.group("params").substring(1) : null;
        try {
            InetAddress inetAddress = InetAddress.getByName(host);
            if (port == null) {
                return new EtherNetIpTcpPlcConnection(inetAddress, params);
            } else {
                return new EtherNetIpTcpPlcConnection(inetAddress, Integer.valueOf(port), params);
            }
        } catch (UnknownHostException e) {
            throw new PlcConnectionException("Unknown host" + host, e);
        }
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("EtherNet/IP connections don't support authentication.");
    }

}
