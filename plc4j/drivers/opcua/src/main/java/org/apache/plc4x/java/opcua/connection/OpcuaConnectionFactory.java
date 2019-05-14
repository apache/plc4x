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
package org.apache.plc4x.java.opcua.connection;

import java.net.InetAddress;
import java.util.Objects;
/**
 * @author Matthias Milan Stlrljic
 * Created by Matthias Milan Stlrljic on 10.05.2019
 */
public class OpcuaConnectionFactory {

    public OpcuaTcpPlcConnection opcuaTcpPlcConnectionOf(InetAddress address, Integer port, String params, int requestTimeout) {
        Objects.requireNonNull(address);

        if (port == null) {
            return OpcuaTcpPlcConnection.of(address, params, requestTimeout);
        } else {
            return OpcuaTcpPlcConnection.of(address, port, params, requestTimeout);
        }
    }

}
