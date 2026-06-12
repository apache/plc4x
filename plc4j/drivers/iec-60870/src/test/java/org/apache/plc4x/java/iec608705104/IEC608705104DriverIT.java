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
package org.apache.plc4x.java.iec608705104;

import org.apache.plc4x.java.utils.testutils.driver.DriverTestsuiteRunner;
import org.junit.jupiter.api.Disabled;

/**
 * Disabled while the test-utils {@code DriverTestsuiteRunner} is updated to
 * handle the array-wrapper XML form used by the driver testsuites (the same
 * issue that disables the matching EIP driver IT — see
 * {@code EIPDriverIT}). Parser-serializer coverage of the wire format is
 * provided by {@link IEC608705104ParserSerializerTest} in the meantime.
 */
@Disabled("Pending fix in test-utils DriverTestsuiteRunner — see class Javadoc")
public class IEC608705104DriverIT extends DriverTestsuiteRunner {

    public IEC608705104DriverIT() {
        super("/protocols/iec608705104/DriverTestsuite.xml", "org.apache.plc4x.java.iec608705104.readwrite");
    }

}
