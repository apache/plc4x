/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.eip.logix.configuration;

import org.apache.plc4x.java.eip.base.configuration.EIPConfiguration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Description;

public class LogixConfiguration extends EIPConfiguration {

    @ConfigurationParameter("communication-path")
    @Description("The communication path allows for connection routing across multiple backplanes. It uses a common format found in Logix controllers.\n" +
        "It consists of pairs of values, each pair begins with either 1 (Backplane) or 2 (Ethernet), followed by a slot in the case of a backplane address,\n" +
        "or if using Ethernet an ip address. e.g. [1,4,2,192.168.0.1,1,1] - Routes to the 4th slot in the first rack, which is a ethernet module, it then connects to the address 192.168.0.1, then finds the module in slot 1.")
    private String communicationPath;

    public String getCommunicationPath() { return this.communicationPath; }

    public void setCommunicationPath(String routingAddress) { this.communicationPath = communicationPath; }

}
