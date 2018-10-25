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
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.connection.AdsConnectionFactory;
import org.apache.plc4x.java.spi.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Implementation of the ADS protocol, based on:
 * - ADS Protocol
 * - TCP
 * - Serial
 */
public class AdsPlcDriver implements PlcDriver {

    public static final Pattern ADS_ADDRESS_PATTERN =
        Pattern.compile("(?<targetAmsNetId>" + AmsNetId.AMS_NET_ID_PATTERN + "):(?<targetAmsPort>" + AmsPort.AMS_PORT_PATTERN + ")"
            + "(/"
            + "(?<sourceAmsNetId>" + AmsNetId.AMS_NET_ID_PATTERN + "):(?<sourceAmsPort>" + AmsPort.AMS_PORT_PATTERN + ")"
            + ")?");
    public static final Pattern INET_ADDRESS_PATTERN = Pattern.compile("tcp://(?<host>[\\w.]+)(:(?<port>\\d*))?");
    public static final Pattern SERIAL_PATTERN = Pattern.compile("serial://(?<serialDefinition>((?!/\\d).)*)");
    public static final Pattern ADS_URI_PATTERN = Pattern.compile("^ads:(" + INET_ADDRESS_PATTERN + "|" + SERIAL_PATTERN + ")/" + ADS_ADDRESS_PATTERN + "(\\?.*)?");

    private AdsConnectionFactory adsConnectionFactory;

    public AdsPlcDriver() {
        this.adsConnectionFactory = new AdsConnectionFactory();
    }

    public AdsPlcDriver(AdsConnectionFactory adsConnectionFactory) {
        this.adsConnectionFactory = adsConnectionFactory;
    }

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
                "Connection url " + url + " doesn't match 'ads://{{host|ip}|serial:definition}/{targetAmsNetId}:{targetAmsPort}/{sourceAmsNetId}:{sourceAmsPort}' RAW:" + ADS_URI_PATTERN);
        }
        String host = matcher.group("host");
        String serialDefinition = matcher.group("serialDefinition");
        String portString = matcher.group("port");
        Integer port = StringUtils.isNotBlank(portString) ? Integer.parseInt(portString) : null;
        AmsNetId targetAmsNetId = AmsNetId.of(matcher.group("targetAmsNetId"));
        AmsPort targetAmsPort = AmsPort.of(matcher.group("targetAmsPort"));
        String sourceAmsNetIdString = matcher.group("sourceAmsNetId");
        AmsNetId sourceAmsNetId = StringUtils.isNotBlank(sourceAmsNetIdString) ? AmsNetId.of(sourceAmsNetIdString) : null;
        String sourceAmsPortString = matcher.group("sourceAmsPort");
        AmsPort sourceAmsPort = StringUtils.isNotBlank(sourceAmsPortString) ? AmsPort.of(sourceAmsPortString) : null;

        if (serialDefinition != null) {
            return adsConnectionFactory.adsSerialPlcConnectionOf(serialDefinition, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
        } else {
            try {
                return adsConnectionFactory.adsTcpPlcConnectionOf(InetAddress.getByName(host), port, targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort);
            } catch (UnknownHostException e) {
                throw new PlcConnectionException(e);
            }
        }
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("Basic ADS connections don't support authentication.");
    }

}
