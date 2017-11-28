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


import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.s7.connection.S7PlcConnection;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

public class S7PlcDriverTest {

    @Disabled("We first have to find/build some tool to help test these connections.")
    @Test
    @Tag("fast")
    void getConnectionTest() throws PlcException {
        S7PlcConnection s7Connection = (S7PlcConnection)
            new PlcDriverManager().getConnection("s7://localhost/1/2");
        Assertions.assertEquals(s7Connection.getHostName(), "localhost");
        Assertions.assertEquals(s7Connection.getRack(), 1);
        Assertions.assertEquals(s7Connection.getSlot(), 2);
    }

    /**
     * In this test case the 's7' driver should report an invalid url format.
     *
     * @throws PlcException something went wrong
     */
    @Test
    @Tag("fast")
    void getConnectionInvalidUrlTest() throws PlcException {
        Assertions.assertThrows(PlcConnectionException.class,
            () -> new PlcDriverManager().getConnection("s7://localhost/hurz/2"));
    }

    /**
     * In this test case the 's7' driver should report an error as this protocol
     * doesn't support authentication.
     *
     * @throws PlcException something went wrong
     */
    @Test
    @Tag("fast")
    void getConnectionWithAuthenticationTest() throws PlcException {
        Assertions.assertThrows(PlcConnectionException.class,
            () -> new PlcDriverManager().getConnection("s7://localhost/1/2",
                new PlcUsernamePasswordAuthentication("user", "pass")));
    }

}
