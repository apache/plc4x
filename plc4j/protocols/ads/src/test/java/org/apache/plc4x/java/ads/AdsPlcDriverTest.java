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


import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.ads.connection.AdsTcpPlcConnection;
import org.apache.plc4x.java.ads.util.TcpHexDumper;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.junit.Rule;
import org.junit.Test;

import java.util.regex.Pattern;

import static org.apache.plc4x.java.ads.AdsPlcDriver.*;
import static org.hamcrest.Matchers.is;
import static org.junit.Assert.*;

public class AdsPlcDriverTest {

    @Rule
    public TcpHexDumper tcpHexDumper = new TcpHexDumper(0, 2);

    @Test
    public void testAdsAddressPattern() throws Exception {
        assertMatching(ADS_ADDRESS_PATTERN, "0.0.0.0.0.0:13");
        assertMatching(ADS_ADDRESS_PATTERN, "0.0.0.0.0.0:13/0.0.0.0.0.0:13");

        assertMatching(INET_ADDRESS_PATTERN, "localhost");
        assertMatching(INET_ADDRESS_PATTERN, "localhost:3131");
        assertMatching(INET_ADDRESS_PATTERN, "www.google.de");
        assertMatching(INET_ADDRESS_PATTERN, "www.google.de:443");

        assertMatching(SERIAL_PATTERN, "serial:/dev/com1");
        assertMatching(SERIAL_PATTERN, "serial:COM1");
        assertMatching(SERIAL_PATTERN, "serial:/dev/ttyUSB0");

        assertMatching(ADS_URI_PATTERN, "ads://www.google.de/0.0.0.0.0.0:13");
        assertMatching(ADS_URI_PATTERN, "ads://www.google.de:443/0.0.0.0.0.0:13");
        assertMatching(ADS_URI_PATTERN, "ads://serial:/dev/com1/0.0.0.0.0.0:13");
        assertMatching(ADS_URI_PATTERN, "ads://serial:COM1/0.0.0.0.0.0:13");
        assertMatching(ADS_URI_PATTERN, "ads://serial:/dev/ttyUSB0/0.0.0.0.0.0:13");
    }

    private void assertMatching(Pattern pattern, String match) {
        if (!pattern.matcher(match).matches()) {
            fail(pattern + "doesn't match " + match);
        }
    }

    @Test
    public void getConnection() throws Exception {
        AdsTcpPlcConnection adsConnection = (AdsTcpPlcConnection)
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
        new PlcDriverManager().getConnection("ads://nowhere:8080/0.0.0.0.0.0:13");
    }

    @Test(expected = PlcConnectionException.class)
    public void getConnectionUnknownPort() throws Exception {
        new PlcDriverManager().getConnection("ads://nowhere:unknown/0.0.0.0.0.0:13");
    }

    /**
     * In this test case the 'ads' driver should report an invalid url format.
     *
     * @throws PlcException something went wrong
     */
    @Test(expected = PlcConnectionException.class)
    public void getConnectionInvalidUrl() throws PlcException {
        new PlcDriverManager().getConnection("ads://localhost/hurz/2");
    }

    @Test
    public void getProtocol() {
        AdsPlcDriver driver = new AdsPlcDriver();
        assertThat(driver.getProtocolCode(), is("ads"));
        assertThat(driver.getProtocolName(), is("Beckhoff Twincat ADS"));
    }


}
