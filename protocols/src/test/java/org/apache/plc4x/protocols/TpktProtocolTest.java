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

package org.apache.plc4x.protocols;

import org.apache.daffodil.tdml.DFDLTestSuite;
import org.apache.daffodil.tdml.Runner;
import org.apache.daffodil.util.Misc;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

public class TpktProtocolTest {

    private DFDLTestSuite testSuite;

    @BeforeEach
    public void setup() {
        System.out.println("Starting");
        String tpktProtocolTestCases = "org/apache/plc4x/protocols/tpkt-protocol.tdml";
        testSuite = new DFDLTestSuite(Misc.getRequiredResource(tpktProtocolTestCases), true, true, false,
            Runner.defaultRoundTripDefaultDefault(), Runner.defaultValidationDefaultDefault());
    }

    @Test
    public void tpktPacketContainingCotpConnectResponse() {
        System.out.println("Running");
        testSuite.runOneTest("tpktPacketContainingCotpConnectResponse", scala.Option.apply(null), false);
        System.out.println("Done");
    }

}
