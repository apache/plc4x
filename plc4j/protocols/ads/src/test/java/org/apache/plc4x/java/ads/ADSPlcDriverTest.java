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
import org.apache.plc4x.java.ads.connection.ADSPlcConnection;
import org.apache.plc4x.java.ads.util.ExtendWithTcpHexDumper;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@ExtendWithTcpHexDumper(value = ADSPlcConnection.TCP_PORT, shutdownTimeout = 3)
public class ADSPlcDriverTest {

    private static final Logger logger = LoggerFactory.getLogger(ADSPlcDriverTest.class);

    @Test
    @Tag("fast")
    void getConnection() throws Exception {
        ADSPlcConnection adsConnection = (ADSPlcConnection)
            new PlcDriverManager().getConnection("ads://localhost/0.0.0.0.0.0:13");
        Assertions.assertEquals(adsConnection.getHostName(), "localhost");
        Assertions.assertEquals(adsConnection.getTargetAmsNetId().toString(), "0.0.0.0.0.0");
        Assertions.assertEquals(adsConnection.getTargetAmsPort().toString(), "13");
        adsConnection.close();
    }

    /**
     * In this test case the 'ads' driver should report an invalid url format.
     *
     * @throws PlcException something went wrong
     */
    @Test
    @Tag("fast")
    void getConnectionInvalidUrl() throws PlcException {
        Assertions.assertThrows(PlcConnectionException.class,
            () -> new PlcDriverManager().getConnection("ads://localhost/hurz/2"));
    }

}
