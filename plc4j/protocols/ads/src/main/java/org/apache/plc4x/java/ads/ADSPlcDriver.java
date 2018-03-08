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
package org.apache.plc4x.java.ads;

import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.connection.ADSTcpPlcConnection;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Implementation of the ADS protocol, based on:
 * - ADS Protocol
 * - TCP
 */
public class ADSPlcDriver implements PlcDriver {

    private static final Pattern ADS_ADDRESS_PATTERN =
        Pattern.compile("(?<targetAmsNetId>" + AMSNetId.AMS_NET_ID_PATTERN + "):(?<targetAmsPort>" + AMSPort.AMS_PORT_PATTERN + ")"
            + "(/"
            + "(?<sourceAmsNetId>" + AMSNetId.AMS_NET_ID_PATTERN + "):(?<sourceAmsPort>" + AMSPort.AMS_PORT_PATTERN + ")"
            + ")?");
    private static final Pattern ADS_URI_PATTERN = Pattern.compile("^ads://(?<host>\\w+)(:(?<port>\\d*))?/" + ADS_ADDRESS_PATTERN);

    @Override
    public String getProtocolCode() {
        return "ads";
    }

    @Override
    public String getProtocolName() {
        return "Beckhoff Twincat ADS";
    }

    @Override
    public PlcConnection connect(String url) throws PlcConnectionException {
        Matcher matcher = ADS_URI_PATTERN.matcher(url);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection url " + url + " doesn't match 'ads://{host|ip}/{targetAmsNetId}:{targetAmsPort}/{sourceAmsNetId}:{sourceAmsPort}' RAW:" + ADS_URI_PATTERN);
        }
        String host = matcher.group("host");
        String portString = matcher.group("port");
        Integer port = StringUtils.isNotBlank(portString) ? Integer.parseInt(portString) : null;
        AMSNetId targetAmsNetId = AMSNetId.of(matcher.group("targetAmsNetId"));
        AMSPort targetAmsPort = AMSPort.of(matcher.group("targetAmsPort"));
        String sourceAmsNetIdString = matcher.group("sourceAmsNetId");
        AMSNetId sourceAmsNetId = StringUtils.isNotBlank(sourceAmsNetIdString) ? AMSNetId.of(sourceAmsNetIdString) : null;
        String sourceAmsPortString = matcher.group("sourceAmsPort");
        AMSPort sourceAmsPort = StringUtils.isNotBlank(sourceAmsPortString) ? AMSPort.of(sourceAmsPortString) : null;
        try {
            return new ADSTcpPlcConnection(InetAddress.getByName(host), port, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
        } catch (UnknownHostException e) {
            throw new PlcConnectionException(e);
        }
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("Basic ADS connections don't support authentication.");
    }

}
