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

import static org.assertj.core.api.Assertions.assertThat;

import org.apache.plc4x.java.isotp.netty.model.types.ParameterCode;
import org.apache.plc4x.java.isotp.netty.model.types.TpduSize;
import org.apache.plc4x.test.FastTests;
import org.junit.Test;
import org.junit.experimental.categories.Category;

public class ParameterTests {

    @Test
    @Category(FastTests.class)
    public void checksumPartameter() {
        ChecksumParameter checksumParameter = new ChecksumParameter((byte)1);

        assertThat(checksumParameter.getChecksum()).isEqualTo((byte)1).withFailMessage("Checksum incorrect");
        assertThat(checksumParameter.getType()).isEqualTo(ParameterCode.CHECKSUM);
    }

    @Test
    @Category(FastTests.class)
    public void disconnectAdditionalInformationParameter() {
        byte[] data = {(byte)1, (byte)2};
        DisconnectAdditionalInformationParameter disconnectParameter = new DisconnectAdditionalInformationParameter(data);

        assertThat(disconnectParameter.getData()[0]).isEqualTo((byte)1).withFailMessage("Return parameter incorrect");
        assertThat(disconnectParameter.getData()[1]).isEqualTo((byte)2).withFailMessage("Return parameter incorrect");
        assertThat(disconnectParameter.getType()).isEqualTo(ParameterCode.DISCONNECT_ADDITIONAL_INFORMATION);
    }

    @Test
    @Category(FastTests.class)
    public void tpduSizeParameter() {
        TpduSizeParameter tpduSizeParameter = new TpduSizeParameter(TpduSize.SIZE_512);

        assertThat(tpduSizeParameter.getTpduSize()).isEqualTo(TpduSize.SIZE_512).withFailMessage("Tpdu size incorrect");
        assertThat(tpduSizeParameter.getType()).isEqualTo(ParameterCode.TPDU_SIZE);
    }

}