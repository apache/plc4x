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
package org.apache.plc4x.java.profinet;


import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.exception.PlcConnectionException;
import org.apache.plc4x.java.exception.PlcException;
import org.apache.plc4x.java.s7.connection.S7PlcConnection;
import org.testng.Assert;
import org.testng.annotations.Test;

public class S7PlcDriverTest {

    @Test(groups = { "fast" })
    public void getConnectionTest() throws PlcException{
        S7PlcConnection profinetConnection = (S7PlcConnection)
            new PlcDriverManager().getConnection("s7://localhost/1/2");
        Assert.assertEquals(profinetConnection.getHostName(), "localhost");
        Assert.assertEquals(profinetConnection.getRack(), 1);
        Assert.assertEquals(profinetConnection.getSlot(), 2);
    }

    /**
     * In this test case the 's7' driver should report an invalid url format.
     * @throws PlcException
     */
    @Test(groups = { "fast" }, expectedExceptions = {PlcConnectionException.class})
    public void getConnectionInvalidUrlTest() throws PlcException {
        new PlcDriverManager().getConnection("s7://localhost/hurz/2");
    }

    /**
     * In this test case the 's7' driver should report an error as this protocol
     * doesn't support authentication.
     * @throws PlcException
     */
    @Test(groups = { "fast" }, expectedExceptions = {PlcConnectionException.class})
    public void getConnectionWithAuthenticationTest() throws PlcException {
        new PlcDriverManager().getConnection("s7://localhost/1/2",
            new PlcUsernamePasswordAuthentication("user", "pass"));
    }

}
