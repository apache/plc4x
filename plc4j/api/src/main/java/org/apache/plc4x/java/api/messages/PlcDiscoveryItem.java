/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.api.messages;

import java.net.URL;
import java.util.Map;

public interface PlcDiscoveryItem {

    /**
     * @return returns the protocol-code part of the url (s7, modbus, ads, ...)
     */
    String getProtocolCode();

    /**
     * @return returns the transport part of the url (tcp, udp, serial, raw, ...)
     */
    String getTransportCode();

    /**
     * @return returns the part of the url, the given transport needs in order to connect (plc.mycompany.de, 192.168.42.23, /dev/serial, COM1)
     */
    URL getTransportUrl();

    /**
     * @return returns a map of all configuration options (usually encoded after the transport url's "?" character (rack=1&slot=1, little-endian=true, ...)
     */
    Map<String, String> getOptions();

    /**
     * @return returns something I bet made sense some time, but I have forgotten why I added it to plc4go ;-)
     */
    String getName();

    /**
     * @return returns a plc4x connection string that can be used in any PLC4X driver to connect to the given device (Generally just a concatenation of the other parts of this object)
     */
    String getConnectionUrl();

}
