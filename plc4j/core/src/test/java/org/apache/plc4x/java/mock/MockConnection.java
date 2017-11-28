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
package org.apache.plc4x.java.mock;

import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.connection.AbstractPlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.model.Address;

public class MockConnection extends AbstractPlcConnection {

    private final PlcAuthentication authentication;

    public MockConnection(PlcAuthentication authentication) {
        this.authentication = authentication;
    }

    public PlcAuthentication getAuthentication() {
        return authentication;
    }

    @Override
    public void connect() throws PlcConnectionException {

    }

    @Override
    public void close() throws Exception {

    }

    @Override
    public Address parseAddress(String addressString) throws PlcException {
        return null;
    }

}
