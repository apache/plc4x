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

import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.s7.connection.S7PlcConnection;
import org.apache.plc4x.test.FastTests;
import org.junit.Ignore;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;

public class S7PlcDriverTest {

    @Test
    @Category(FastTests.class)
    public void getProtocolCode() {
        assertThat(new S7PlcDriver().getProtocolCode(), equalTo("s7"));
    }

    @Test
    @Category(FastTests.class)
    public void getProtocolName() {
        assertThat(new S7PlcDriver().getProtocolName(), equalTo("Siemens S7 (Basic)"));
    }

    @Test
    @Category(FastTests.class)
    public void getConnection() throws PlcException {
        S7PlcConnection s7Connection = (S7PlcConnection) new S7PlcDriver().connect("s7://localhost/1/2");
        assertThat(s7Connection.getRack(), equalTo(1));
        assertThat(s7Connection.getSlot(), equalTo(2));
    }

    @Test(expected = PlcConnectionException.class)
    @Category(FastTests.class)
    @Ignore("This test tends to fail on systems with DNS providers that proivde default IPs")
    public void getConnectionToUnknownHost() throws PlcException {
        new S7PlcDriver().connect("s7://IHopeThisHostDoesntExistAAAAAAAAhhhhhhh/1/2");
    }

    /**
     * In this test case the 's7' driver should report an invalid url format.
     */
    @Test(expected = PlcConnectionException.class)
    @Category(FastTests.class)
    public void getConnectionInvalidUrl() throws PlcConnectionException {
        new S7PlcDriver().connect("s7://localhost/hurz/2");
    }

    /**
     * In this test case the 's7' driver should report an error as this protocol
     * doesn't support authentication.
     */
    @Test(expected = PlcConnectionException.class)
    @Category(FastTests.class)
    public void getConnectionWithAuthentication() throws PlcConnectionException {
        new S7PlcDriver().connect("s7://localhost/1/2", new PlcUsernamePasswordAuthentication("user", "pass"));
    }

}
