/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.utils.connectionpool;

import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;

import java.net.URI;
import java.net.URISyntaxException;
import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class PoolKeyFactory {

    // TODO 11.12.18 jf: add a property to the drivers to keep this generic
    public PoolKey getPoolKey(String url, PlcAuthentication plcAuthentication) throws PlcConnectionException {
        Objects.requireNonNull(url);
        URI connectionUri;
        try {
            connectionUri = new URI(url);
        } catch (URISyntaxException e) {
            throw new PlcConnectionException("Invalid plc4j connection string '" + url + "'", e);
        }
        String protocol = connectionUri.getScheme().toLowerCase();
        switch (protocol) {
            // Currently this is disabled due to 2 reasons
            // First, see PLC4X-223 it needs to be migrated to new URI Syntax
            // Second, we have to decide which parameters uniquely identify a connection and which
            // not. See PLC4X-224
            /*
            case "s7":
                return new PoolKey(url, plcAuthentication) {
                    private final Pattern s7URIPattern = Pattern.compile("^(?<poolablePart>s7://((?<ip>[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})|(?<hostname>[a-zA-Z0-9\\.\\-]+))(:(?<port>[0-9]{1,5}))?)(?<params>\\?.*)?");

                    @Override
                    public String getPoolableKey() {
                        Matcher matcher = s7URIPattern.matcher(url);
                        if (!matcher.matches()) {
                            throw new IllegalArgumentException(url + " doesn't match " + s7URIPattern);
                        }
                        return Objects.requireNonNull(matcher.group("poolablePart"));
                    }
                };
            case "ads":
                return new PoolKey(url, plcAuthentication) {
                    private final Pattern amsPortPattern = Pattern.compile("\\d+");
                    private final Pattern amsNetIdPattern = Pattern.compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}");
                    private final Pattern adsAddressPattern =
                        Pattern.compile("(?<targetAmsNetId>" + amsNetIdPattern + "):(?<targetAmsPort>" + amsPortPattern + ")"
                            + "(/"
                            + "(?<sourceAmsNetId>" + amsNetIdPattern + "):(?<sourceAmsPort>" + amsPortPattern + ")"
                            + ")?");
                    private final Pattern inetAddressPattern = Pattern.compile("tcp://(?<host>[\\w.]+)(:(?<port>\\d*))?");
                    private final Pattern serialPattern = Pattern.compile("serial://(?<serialDefinition>((?!/\\d).)*)");
                    private final Pattern adsUriPattern = Pattern.compile("^(?<poolablePart>ads:(" + inetAddressPattern + "|" + serialPattern + "))/" + adsAddressPattern + "(\\?.*)?");

                    @Override
                    public String getPoolableKey() {
                        Matcher matcher =
                            adsUriPattern.matcher(url);
                        if (!matcher.matches()) {
                            throw new IllegalArgumentException(url + " doesn't match " + adsUriPattern);
                        }
                        return Objects.requireNonNull(matcher.group("poolablePart"));
                    }
                };
            case "modbus":
                return new PoolKey(url, plcAuthentication) {
                    private final Pattern inetAddressPattern = Pattern.compile("tcp://(?<host>[\\w.]+)(:(?<port>\\d*))?");
                    private final Pattern serialPattern = Pattern.compile("serial://(?<serialDefinition>((?!/\\d).)*)");
                    private final Pattern modbusUriPattern = Pattern.compile("^(?<poolablePart>modbus:(" + inetAddressPattern + "|" + serialPattern + "))/?" + "(?<params>\\?.*)?");

                    @Override
                    public String getPoolableKey() {
                        Matcher matcher = modbusUriPattern.matcher(url);
                        if (!matcher.matches()) {
                            throw new IllegalArgumentException(url + " doesn't match " + modbusUriPattern);
                        }
                        return Objects.requireNonNull(matcher.group("poolablePart"));
                    }
                };
             */
            default:
                return new PoolKey(url, plcAuthentication) {
                    @Override
                    public String getPoolableKey() {
                        return url;
                    }
                };
        }
    }
}
