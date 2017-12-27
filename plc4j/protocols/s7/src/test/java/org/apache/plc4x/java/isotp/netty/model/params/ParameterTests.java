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

package org.apache.plc4x.java.isotp.netty.model.params;

import org.apache.plc4x.java.isotp.netty.model.types.DeviceGroup;
import org.apache.plc4x.java.isotp.netty.model.types.ParameterCode;
import org.apache.plc4x.java.isotp.netty.model.types.TpduSize;
import org.junit.jupiter.api.*;

import static org.junit.jupiter.api.Assertions.*;

class ParameterTest {
    private ParameterTest tsapParaameter;

    @BeforeEach
    void setUp() {
    }

    @AfterEach
    void tearDown() {
        tsapParaameter = null;
    }

    @Test
    @Tag("fast")
    void checksumPartameter() {
        ChecksumParameter checksumParameter = new ChecksumParameter((byte)1);

        assertTrue(checksumParameter.getChecksum() == (byte)1, "Checksum incorrect");
        assertTrue(checksumParameter.getType() == ParameterCode.CHECKSUM);
    }

    @Test
    @Tag("fast")
    void disconnectAdditionalInformationParameter() {
        byte[] data = {(byte)1, (byte)2};
        DisconnectAdditionalInformationParameter disconnectParameter = new DisconnectAdditionalInformationParameter(data);

        assertTrue(disconnectParameter.getData()[0] == (byte)1, "Return parameter incorrect");
        assertTrue(disconnectParameter.getData()[1] == (byte)2, "Return parameter incorrect");
        assertTrue(disconnectParameter.getType() == ParameterCode.DISCONNECT_ADDITIONAL_INFORMATION);
    }

    @Test
    @Tag("fast")
    void tpduSizeParameter() {
        TpduSizeParameter tpduSizeParameter = new TpduSizeParameter(TpduSize.SIZE_512);

        assertTrue(tpduSizeParameter.getTpduSize() == TpduSize.SIZE_512, "Tpdu size incorrect");
        assertTrue(tpduSizeParameter.getType() == ParameterCode.TPDU_SIZE);
    }
}