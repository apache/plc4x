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
package org.apache.plc4x.java.profinet.connection;

import org.apache.plc4x.java.connection.PlcConnection;

public class ProfinetPlcConnection implements PlcConnection {

    private final String hostName;
    private final int rack;
    private final int slot;

    public ProfinetPlcConnection(String hostName, int rack, int slot) {
        this.hostName = hostName;
        this.rack = rack;
        this.slot = slot;
    }

    public String getHostName() {
        return hostName;
    }

    public int getRack() {
        return rack;
    }

    public int getSlot() {
        return slot;
    }

}
