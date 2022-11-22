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
package org.apache.plc4x.java.utils.connectionpool;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.PlcDriver;

import static org.mockito.Mockito.RETURNS_DEEP_STUBS;
import static org.mockito.Mockito.mock;

public class PooledDummyDriver implements PlcDriver {

    private PlcDriver mockedPlcDriver = mock(PlcDriver.class, RETURNS_DEEP_STUBS);

    @Override
    public String getProtocolCode() {
        return PooledDummyDriver.class.getName();
    }

    @Override
    public String getProtocolName() {
        return mockedPlcDriver.getProtocolCode();
    }

    @Override
    public PlcConnection getConnection(String url) throws PlcConnectionException {
        return mockedPlcDriver.getConnection(url);
    }

    @Override
    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        return mockedPlcDriver.getConnection(url, authentication);
    }

}
