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


import io.netty.channel.ConnectTimeoutException;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.ads.connection.AdsConnectionFactory;
import org.apache.plc4x.java.ads.connection.AdsTcpPlcConnection;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.mock.connection.tcp.TcpHexDumper;
import org.apache.plc4x.test.RequireInternetConnection;
import org.junit.jupiter.api.*;

import java.net.ConnectException;
import java.util.regex.Pattern;
import java.util.stream.Stream;

import static org.apache.plc4x.java.ads.AdsPlcDriver.*;
import static org.hamcrest.Matchers.is;
import static org.junit.Assert.*;
import static org.mockito.Mockito.mock;

public class AdsPlcDriverTest {

    public TcpHexDumper tcpHexDumper = new TcpHexDumper(0, 2);

    @BeforeEach
    public void before() throws Throwable {
        tcpHexDumper.before();
    }

    @AfterEach
    public void after() {
        tcpHexDumper.after();
    }

    @Test
    public void testAdsAddressPattern() {
        assertMatching(ADS_ADDRESS_PATTERN, "0.0.0.0.0.0:13");
        assertMatching(ADS_ADDRESS_PATTERN, "0.0.0.0.0.0:13/0.0.0.0.0.0:13");

        assertMatching(INET_ADDRESS_PATTERN, "tcp://localhost");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://localhost:3131");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://www.google.de");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://www.google.de:443");

        assertMatching(SERIAL_PATTERN, "serial:///dev/com1");
        assertMatching(SERIAL_PATTERN, "serial://COM1");
        assertMatching(SERIAL_PATTERN, "serial:///dev/ttyUSB0");

        assertMatching(ADS_URI_PATTERN, "ads:tcp://www.google.de/0.0.0.0.0.0:13");
        assertMatching(ADS_URI_PATTERN, "ads:tcp://www.google.de:443/0.0.0.0.0.0:13");
        assertMatching(ADS_URI_PATTERN, "ads:serial:///dev/com1/0.0.0.0.0.0:13");
        assertMatching(ADS_URI_PATTERN, "ads:serial://COM1/0.0.0.0.0.0:13");
        assertMatching(ADS_URI_PATTERN, "ads:serial:///dev/ttyUSB0/0.0.0.0.0.0:13");

        assertMatching(ADS_URI_PATTERN, "ads:serial:///dev/ttyUSB0/0.0.0.0.0.0:13?some=random&parameters=true");
    }

    @Test
    @RequireInternetConnection
    public void testDriverWithCompleteUrls() {
        AdsPlcDriver SUT = new AdsPlcDriver(mock(AdsConnectionFactory.class));
        Stream.of(
            "ads:tcp://www.google.de/0.0.0.0.0.0:13",
            "ads:tcp://www.google.de:443/0.0.0.0.0.0:13",
            "ads:serial:///dev/com1/0.0.0.0.0.0:13",
            "ads:serial://COM1/0.0.0.0.0.0:13",
            "ads:serial:///dev/ttyUSB0/0.0.0.0.0.0:13"
        ).forEach(url -> {
            try {
                SUT.connect(url);
            } catch (PlcConnectionException e) {
                throw new PlcRuntimeException(e);
            }
        });
    }

    private void assertMatching(Pattern pattern, String match) {
        if (!pattern.matcher(match).matches()) {
            fail(pattern + "doesn't match " + match);
        }
    }

    @Test
    public void getConnection() throws Exception {
        AdsTcpPlcConnection adsConnection = (AdsTcpPlcConnection)
            new PlcDriverManager().getConnection("ads:tcp://localhost:" + tcpHexDumper.getPort() + "/0.0.0.0.0.0:13");
        assertEquals(adsConnection.getTargetAmsNetId().toString(), "0.0.0.0.0.0");
        assertEquals(adsConnection.getTargetAmsPort().toString(), "13");
        adsConnection.close();
    }

    @Test
    public void getConnectionNoAuthSupported() throws Exception {
        Assertions.assertThrows(PlcConnectionException.class, () ->
            new PlcDriverManager().getConnection("ads:tcp://localhost:" + tcpHexDumper.getPort() + "/0.0.0.0.0.0:13",
                new PlcUsernamePasswordAuthentication("admin", "admin"))
        );
    }

    @Test
    public void getConnectionUnknownHost() {
        try {
            new PlcDriverManager().getConnection("ads:tcp://nowhere:8080/0.0.0.0.0.0:13");
            fail();
        } catch (Exception e) {
            assert((e instanceof PlcConnectionException) ||
                // When running this test on systems with an internet provider that redirects
                // all unknown hosts to a default host (For showing adds), one of these will
                // be thrown instead.
                (e instanceof ConnectTimeoutException) ||
                (e.getCause() instanceof ConnectException));
        }
    }

    @Test
    public void getConnectionUnknownPort() throws Exception {
        Assertions.assertThrows(PlcConnectionException.class, () ->
            new PlcDriverManager().getConnection("ads:tcp://nowhere:unknown/0.0.0.0.0.0:13")
        );
    }

    /**
     * In this test case the 'ads' driver should report an invalid url format.
     *
     * @throws PlcException something went wrong
     */
    @Test
    public void getConnectionInvalidUrl() throws PlcException {
        Assertions.assertThrows(PlcConnectionException.class, () ->
            new PlcDriverManager().getConnection("ads:tcp://localhost/hurz/2")
        );
    }

    @Test
    public void getProtocol() {
        AdsPlcDriver driver = new AdsPlcDriver();
        assertThat(driver.getProtocolCode(), is("ads"));
        assertThat(driver.getProtocolName(), is("Beckhoff Twincat ADS"));
    }


}
