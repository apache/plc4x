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
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.utils.testutils.driver.DriverTestsuiteRunner;

/**
 * Runs the Java specific S7 driver testsuite.
 * <p>
 * The shared testsuite in the protocols/s7 module is TPKT/COTP framed and only fits drivers
 * that speak COTP themselves (like plc4go's). This driver delegates that framing to the cotp
 * transport, which the testsuite runner replaces with the synthetic test transport - so its
 * testsuite describes the bare S7 message exchange instead.
 */
public class S7DriverIT extends DriverTestsuiteRunner {

    public S7DriverIT() {
        super("/testsuite/S7DriverTestsuite.xml", "org.apache.plc4x.java.s7.readwrite");
    }

}
