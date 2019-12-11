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
package org.apache.plc4x.java.amsads;

import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.amsads.connection.AdsConnectionFactory;
import org.apache.plc4x.java.amsads.readwrite.AmsNetId;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.spi.PlcDriver;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Stream;

/**
 * Implementation of the ADS protocol, based on:
 * - ADS Protocol
 * - TCP
 * - Serial
 */
public class AMSADSPlcDriver implements PlcDriver {

    public static final Pattern AMS_NET_ID_PATTERN =
        Pattern.compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}");

    public static final Pattern AMS_PORT_PATTERN = Pattern.compile("\\d+");

    public static final Pattern ADS_ADDRESS_PATTERN =
        Pattern.compile("(?<targetAmsNetId>" + AMS_NET_ID_PATTERN + "):(?<targetAmsPort>" + AMS_PORT_PATTERN + ")"
            + "(/"
            + "(?<sourceAmsNetId>" + AMS_NET_ID_PATTERN + "):(?<sourceAmsPort>" + AMS_PORT_PATTERN + ")"
            + ")?");
    public static final Pattern INET_ADDRESS_PATTERN = Pattern.compile("tcp://(?<host>[\\w.]+)(:(?<port>\\d*))?");
    public static final Pattern SERIAL_PATTERN = Pattern.compile("serial://(?<serialDefinition>((?!/\\d).)*)");
    public static final Pattern ADS_URI_PATTERN = Pattern.compile("^ads:(" + INET_ADDRESS_PATTERN + "|" + SERIAL_PATTERN + ")/" + ADS_ADDRESS_PATTERN + "(\\?.*)?");

    private AdsConnectionFactory adsConnectionFactory;

    public AMSADSPlcDriver() {
        this.adsConnectionFactory = new AdsConnectionFactory();
    }

    public AMSADSPlcDriver(AdsConnectionFactory adsConnectionFactory) {
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
        AmsNetId targetAmsNetId = AmsNetIdOf(matcher.group("targetAmsNetId"));
        int targetAmsPort = Integer.parseInt(matcher.group("targetAmsPort"));
        String sourceAmsNetIdString = matcher.group("sourceAmsNetId");
        AmsNetId sourceAmsNetId = StringUtils.isNotBlank(sourceAmsNetIdString) ? AmsNetIdOf(sourceAmsNetIdString) : null;
        String sourceAmsPortString = matcher.group("sourceAmsPort");
        int sourceAmsPort = StringUtils.isNotBlank(sourceAmsPortString) ? Integer.parseInt(sourceAmsPortString) : null;

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

    public static AmsNetId AmsNetIdOf(String address) {
        if (!AMS_NET_ID_PATTERN.matcher(address).matches()) {
            throw new IllegalArgumentException(address + " must match " + AMS_NET_ID_PATTERN);
        }
        String[] split = address.split("\\.");
        short[] shorts = ArrayUtils.toPrimitive(Stream.of(split).map(Integer::parseInt).map(Integer::shortValue).toArray(Short[]::new));
        return new AmsNetId(shorts[5], shorts[4], shorts[3], shorts[2], shorts[1], shorts[0]);
    }

}
