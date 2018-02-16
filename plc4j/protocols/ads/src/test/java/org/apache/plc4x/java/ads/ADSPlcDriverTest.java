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


import org.apache.commons.lang3.RandomStringUtils;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.ads.connection.ADSPlcConnection;
import org.apache.plc4x.java.ads.util.TcpHexDumper;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.junit.Rule;
import org.junit.Test;

import static org.apache.plc4x.java.ads.util.Junit5Backport.assertThrows;
import static org.junit.Assert.assertEquals;

public class ADSPlcDriverTest {

    @Rule
    public TcpHexDumper tcpHexDumper = new TcpHexDumper(0, 2);

    @Test
    public void getConnection() throws Exception {
        ADSPlcConnection adsConnection = (ADSPlcConnection)
            new PlcDriverManager().getConnection("ads://localhost:" + tcpHexDumper.getPort() + "/0.0.0.0.0.0:13");
        assertEquals(adsConnection.getTargetAmsNetId().toString(), "0.0.0.0.0.0");
        assertEquals(adsConnection.getTargetAmsPort().toString(), "13");
        adsConnection.close();
    }

    @Test(expected = PlcConnectionException.class)
    public void getConnectionNoAuthSupported() throws Exception {
        new PlcDriverManager().getConnection("ads://localhost:" + tcpHexDumper.getPort() + "/0.0.0.0.0.0:13",
            new PlcUsernamePasswordAuthentication("admin", "admin"));
    }

    @Test(expected = PlcConnectionException.class)
    public void getConnectionUnknownHost() throws Exception {
        new PlcDriverManager().getConnection("ads://:" + RandomStringUtils.randomAscii(12) + "/0.0.0.0.0.0:13",
            new PlcUsernamePasswordAuthentication("admin", "admin"));
    }

    /**
     * In this test case the 'ads' driver should report an invalid url format.
     *
     * @throws PlcException something went wrong
     */
    @Test
    public void getConnectionInvalidUrl() throws PlcException {
        assertThrows(PlcConnectionException.class,
            () -> new PlcDriverManager().getConnection("ads://localhost/hurz/2"));
    }

}
