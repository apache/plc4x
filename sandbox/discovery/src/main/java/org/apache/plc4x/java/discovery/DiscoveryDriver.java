/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.discovery;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;

/**
 * Implementation of a virtual driver which wraps discovery operations in
 * PLC4X API calls.
 */
public class DiscoveryDriver implements PlcDriver {

    @Override
    public String getProtocolCode() {
        return "discovery";
    }

    @Override
    public String getProtocolName() {
        return "Apache PLC4X Discovery Driver";
    }

    @Override
    public PlcConnection getConnection(String url) {
        return new DiscoveryConnection();
    }

    @Override
    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("Authentication not supported.");
    }

}
