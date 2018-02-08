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

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.s7.connection.S7PlcConnection;
import org.apache.plc4x.test.FastTests;
import org.junit.Ignore;
import org.junit.Test;
import org.junit.experimental.categories.Category;

public class S7PlcDriverTest {

    @Ignore("We first have to find/build some tool to help test these connections.")
    @Test
    @Category(FastTests.class)
    public void getConnection() throws PlcException {
        S7PlcConnection s7Connection = (S7PlcConnection)
            new PlcDriverManager().getConnection("s7://localhost/1/2");
        assertThat(s7Connection.getHostName()).isEqualTo("localhost");
        assertThat(s7Connection.getRack()).isEqualTo(1);
        assertThat(s7Connection.getSlot()).isEqualTo(2);
    }

    /**
     * In this test case the 's7' driver should report an invalid url format.
     *
     * @throws PlcException something went wrong
     */
    @Test
    @Category(FastTests.class)
    public void getConnectionInvalidUrl() throws PlcException {
        assertThatThrownBy(() ->
            new PlcDriverManager().getConnection("s7://localhost/hurz/2"))
            .isInstanceOf(PlcConnectionException.class);
    }

    /**
     * In this test case the 's7' driver should report an error as this protocol
     * doesn't support authentication.
     *
     * @throws PlcException something went wrong
     */
    @Test
    @Category(FastTests.class)
    public void getConnectionWithAuthentication() throws PlcException {
        assertThatThrownBy(() ->
            new PlcDriverManager().getConnection("s7://localhost/1/2",
                new PlcUsernamePasswordAuthentication("user", "pass")))
            .isInstanceOf(PlcConnectionException.class);
    }

}
