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
package org.apache.plc4x.java.s7.readwrite;

import org.apache.plc4x.java.api.PlcConnection;

public class ManualS7DriverTest {

    public static void main(String[] args) throws Exception {
        final S7Driver driver = new S7Driver();
        //final PlcConnection connection = driver.getConnection("s7ng:pcap:///Users/christofer.dutz/Projects/Apache/PLC4X/plc4x/sandbox/test-java-s7-driver/src/test/resources/pcaps/s7-1200-reads-writes.pcapng");
        final PlcConnection connection = driver.getConnection("s7ng://10.10.64.20");
        connection.connect();
        System.out.println(connection.readRequestBuilder().addItem("inputs", "%I0:BYTE").build().execute().get().getObject("inputs"));
    }

}
