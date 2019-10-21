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
package org.apache.plc4x.java.modbus;


import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.mock.connection.tcp.TcpHexDumper;
import org.apache.plc4x.java.modbus.connection.ModbusConnectionFactory;
import org.apache.plc4x.java.modbus.connection.ModbusTcpPlcConnection;
import org.junit.Rule;
import org.junit.Test;

import java.util.regex.Pattern;
import java.util.stream.Stream;

import static org.apache.plc4x.java.modbus.ModbusPlcDriver.INET_ADDRESS_PATTERN;
import static org.apache.plc4x.java.modbus.ModbusPlcDriver.SERIAL_PATTERN;
import static org.hamcrest.Matchers.is;
import static org.junit.Assert.assertThat;
import static org.junit.Assert.fail;
import static org.mockito.Mockito.mock;

public class ModbusPlcDriverTest {

    @Rule
    public TcpHexDumper tcpHexDumper = new TcpHexDumper(0, 2);

    @Test
    public void testModbusAddressPattern() {
        assertMatching(INET_ADDRESS_PATTERN, "tcp://localhost");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://localhost:3131");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://www.google.de");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://www.google.de:443");

        assertMatching(SERIAL_PATTERN, "serial:///dev/com1");
        assertMatching(SERIAL_PATTERN, "serial://COM1");
        assertMatching(SERIAL_PATTERN, "serial:///dev/ttyUSB0");
    }

    @Test
    public void testDriverWithCompleteUrls() {
        ModbusPlcDriver SUT = new ModbusPlcDriver(mock(ModbusConnectionFactory.class));
        Stream.of(
            "modbus:tcp://localhost",
            "modbus:tcp://localhost:443",
            "modbus:serial:///dev/com1",
            "modbus:serial://COM1",
            "modbus:serial:///dev/ttyUSB0"
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
        ModbusTcpPlcConnection modbusConnection = (ModbusTcpPlcConnection)
            new PlcDriverManager().getConnection("modbus:tcp://localhost:" + tcpHexDumper.getPort());
        modbusConnection.close();
    }

    @Test(expected = PlcConnectionException.class)
    public void getConnectionNoAuthSupported() throws Exception {
        new PlcDriverManager().getConnection("modbus:tcp://localhost:" + tcpHexDumper.getPort(),
            new PlcUsernamePasswordAuthentication("admin", "admin"));
    }

    @Test(expected = PlcConnectionException.class)
    public void getConnectionUnknownHost() throws Exception {
        new PlcDriverManager().getConnection("modbus:tcp://IHopeThisHostDoesntExistAAAAAAAAhhhhhhh:8080");
    }

    @Test(expected = PlcConnectionException.class)
    public void getConnectionUnknownPort() throws Exception {
        new PlcDriverManager().getConnection("modbus:tcp://localhost:unknown");
    }

    /**
     * In this test case the 'modbus' driver should report an invalid url format.
     *
     * @throws PlcException something went wrong
     */
    @Test(expected = PlcConnectionException.class)
    public void getConnectionInvalidUrl() throws PlcException {
        new PlcDriverManager().getConnection("modbus:tcp://localhost/hurz/2");
    }

    @Test
    public void getProtocol() {
        ModbusPlcDriver driver = new ModbusPlcDriver();
        assertThat(driver.getProtocolCode(), is("modbus"));
        assertThat(driver.getProtocolName(), is("Modbus (TCP / Serial)"));
    }


}
