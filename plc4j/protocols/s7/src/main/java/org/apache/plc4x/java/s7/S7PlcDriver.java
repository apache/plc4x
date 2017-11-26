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
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.s7.connection.S7PlcConnection;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Implementation of the S7 protocol, based on:
 * - S7 Protocol
 * - ISO Transport Protocol (Class 0) (https://tools.ietf.org/html/rfc905)
 * - ISO on TCP (https://tools.ietf.org/html/rfc1006)
 * - TCP
 */
public class S7PlcDriver implements PlcDriver {

    private static final Pattern S7_URI_PATTERN = Pattern.compile("^s7://(?<host>.*)/(?<rack>\\d{1,4})/(?<slot>\\d{1,4})(\\?.*)?");

    @Override
    public String getProtocolCode() {
        return "s7";
    }

    @Override
    public String getProtocolName() {
        return "Siemens S7 (Basic)";
    }

    @Override
    public PlcConnection connect(String url) throws PlcConnectionException {
        Matcher matcher = S7_URI_PATTERN.matcher(url);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection url doesn't match the format 's7://{host|ip}/{rack}/{slot}'");
        }
        String host = matcher.group("host");
        int rack = Integer.valueOf(matcher.group("rack"));
        int slot = Integer.valueOf(matcher.group("slot"));
        return new S7PlcConnection(host, rack, slot);
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("Basic S7 connections don't support authentication.");
    }

}
