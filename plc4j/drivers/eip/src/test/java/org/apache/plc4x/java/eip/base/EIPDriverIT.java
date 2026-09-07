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
package org.apache.plc4x.java.eip.base;

import org.apache.plc4x.java.utils.testutils.driver.DriverTestsuiteRunner;

/**
 * Replays the recorded EIP handshake and read/write exchanges against the driver.
 *
 * <p>This was disabled for a long time; the cause was never the XML array-wrapper form the
 * old comment blamed, but {@code IncomingPlcMessageHandler} in test-utils not applying the
 * testsuite's declared byte order when serializing injected messages. Every injected packet
 * went out big-endian, so a little-endian driver such as this one read the encapsulation
 * length byte-swapped and waited forever for a packet that never arrived.
 */
public class EIPDriverIT extends DriverTestsuiteRunner {

    public EIPDriverIT() {
        super("/protocols/eip/DriverTestsuite.xml", "org.apache.plc4x.java.eip.readwrite");
    }

}
