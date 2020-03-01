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
package org.apache.plc4x.java.abeth;

import org.apache.plc4x.java.abeth.connection.AbEthPlcConnection;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.spi.PlcDriver;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class AbEthDriver implements PlcDriver {

    private static final Pattern ABETH_URI_PATTERN = Pattern.compile("^ab-eth://(?<host>.*)/(?<station>\\d{1,2})(?<params>\\?.*)?");

    @Override
    public String getProtocolCode() {
        return "ab-eth";
    }

    @Override
    public String getProtocolName() {
        return "Allen Bradley ETH";
    }

    @Override
    public PlcConnection connect(String url) throws PlcConnectionException {
        Matcher matcher = ABETH_URI_PATTERN.matcher(url);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection url doesn't match the format 'ab-eth://{host|ip}/{station}'");
        }
        int station = Integer.parseInt(matcher.group("station"));
        String host = matcher.group("host");
        String params = matcher.group("params") != null ? matcher.group("params").substring(1) : null;

        try {
            InetAddress serverInetAddress = InetAddress.getByName(host);
            return new AbEthPlcConnection(serverInetAddress, station, params);
        } catch (UnknownHostException e) {
            throw new PlcConnectionException("Error parsing address", e);
        } catch (Exception e) {
            throw new PlcConnectionException("Error connecting to host", e);
        }
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("AB-ETH connections don't support authentication.");
    }

}
