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

package org.apache.plc4x.java.isotp.netty.model.types;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import org.apache.plc4x.test.FastTests;
import org.junit.Test;
import org.junit.experimental.categories.Category;

public class IsotpTypeTests {

    @Test
    @Category(FastTests.class)
    public void deviceGroup() {
        DeviceGroup deviceGroup;

        deviceGroup = DeviceGroup.PG_OR_PC;
        assertThat(deviceGroup.getCode()).isEqualTo((byte)1);

        deviceGroup = DeviceGroup.OS;
        assertThat(deviceGroup.getCode()).isEqualTo((byte)2);

        deviceGroup = DeviceGroup.OTHERS;
        assertThat(deviceGroup.getCode()).isEqualTo((byte)3);
    }

    @Test
    @Category(FastTests.class)
    public void deviceGroupUnknown() {
        DeviceGroup deviceGroup = DeviceGroup.valueOf((byte)0x40);

        assertThat(deviceGroup).isNull();
    }


    @Test
    @Category(FastTests.class)
    public void disconnectReason() {
        DisconnectReason disconnectReason = DisconnectReason.ADDRESS_UNKNOWN;

        assertThat(DisconnectReason.valueOf((byte)3)).isEqualTo(DisconnectReason.ADDRESS_UNKNOWN).withFailMessage("3 incorrectly mapped");
        assertThat(disconnectReason.getCode()).isEqualTo((byte)3);
    }

    @Test
    @Category(FastTests.class)
    public void diosconectReasonUnknown() {
        DisconnectReason disconnectReason = DisconnectReason.valueOf((byte)4);

        assertThat(disconnectReason).isNull();
    }

    @Test
    @Category(FastTests.class)
    public void parameterCode() {
        ParameterCode parameterCode = ParameterCode.CALLING_TSAP;

        assertThat(ParameterCode.valueOf((byte)0xC1)).isEqualTo(ParameterCode.CALLING_TSAP);
        assertThat(parameterCode.getCode()).isEqualTo((byte)0xC1);
    }

    @Test
    @Category(FastTests.class)
    public void parameterCodeUnknown() {
        ParameterCode parameterCode = ParameterCode.valueOf((byte)0x90);

        assertThat(parameterCode).isNull();
    }

    @Test
    @Category(FastTests.class)
    public void protocolClass() {
        ProtocolClass protocolClass;

        protocolClass = ProtocolClass.CLASS_1;
        assertThat(protocolClass.getCode()).isEqualTo((byte)0x10);

        protocolClass = ProtocolClass.CLASS_2;
        assertThat(protocolClass.getCode()).isEqualTo((byte)0x20);

        protocolClass = ProtocolClass.CLASS_3;
        assertThat(protocolClass.getCode()).isEqualTo((byte)0x30);

        protocolClass = ProtocolClass.CLASS_4;
        assertThat(protocolClass.getCode()).isEqualTo((byte)0x40);
    }

    @Test
    @Category(FastTests.class)
    public void protocolClassUnknown() {
        ProtocolClass protocolClass = ProtocolClass.valueOf((byte)0x50);

        assertThat(protocolClass).isNull();
    }

    @Test
    @Category(FastTests.class)
    public void rejectCause() {
        RejectCause rejectCause = RejectCause.INVALID_PARAMETER_TYPE;

        assertThat(RejectCause.valueOf((byte)0x03)).isEqualTo(RejectCause.INVALID_PARAMETER_TYPE);
        assertThat(rejectCause.getCode()).isEqualTo((byte)0x03);
    }

    @Test
    @Category(FastTests.class)
    public void rejectClauseUnknown() {
        RejectCause rejectCause = RejectCause.valueOf((byte)0x90);

        assertThat(rejectCause).isNull();
    }

    @Test
    @Category(FastTests.class)
    public void tpduCode() {
        TpduCode tpduCode = TpduCode.DATA;

        assertThat(TpduCode.valueOf((byte)0xF0)).isEqualTo(TpduCode.DATA);
        assertThat(tpduCode.getCode()).isEqualTo((byte)0xF0);
    }

    @Test
    @Category(FastTests.class)
    public void tpduCodeUnknown() {
        TpduCode tpduCode = TpduCode.valueOf((byte)0x01);

        assertThat(TpduCode.valueOf((byte)0xFF)).isEqualTo(TpduCode.TPDU_UNKNOWN);
        assertThat(tpduCode.getCode()).isEqualTo((byte)0xFF);
    }
    
    @Test
    @Category(FastTests.class)
    public void typduSize() {
        TpduSize tpduSize = TpduSize.SIZE_128;

        assertThat(TpduSize.valueOf((byte)0x07)).isEqualTo(TpduSize.SIZE_128);
        assertThat(tpduSize.getCode()).isEqualTo((byte)0x07);
        assertThat(tpduSize.getValue()).isEqualTo(128);
    }

    @Test
    @Category(FastTests.class)
    public void tpduSizeUnknown() {
        TpduSize tpduSize = TpduSize.valueOf((byte)0x06);

        assertThat(tpduSize).isNull();
    }

    /**
     * If we are requesting exactly the size of one of the iso tp
     * pdu sizes, then exactly that box should be returned.
     */
    @Test
    @Category(FastTests.class)
    public void tpduValueForGivenExactFit() {
        TpduSize tpduSize = TpduSize.valueForGivenSize(256);

        assertThat(tpduSize).isEqualTo(TpduSize.SIZE_256);
    }

    /**
     * In this case we have a given value that is in-between the boundaries of
     * a pdu box, the method should return the next larger box.
     */
    @Test
    @Category(FastTests.class)
    public void tpduValueForGivenIntermediateSize() {
        TpduSize tpduSize = TpduSize.valueForGivenSize(222);

        assertThat(tpduSize).isEqualTo(TpduSize.SIZE_256);
        assertThat(tpduSize.getValue()).isNotEqualTo(222);
    }

    /**
     * This test should cause an exception as the tpdu size has to be greater
     * than 0 in any case.
     */
    @Test
    @Category(FastTests.class)
    public void tpduValueForGivenTooSmallSize() {
        assertThatThrownBy(() ->
            TpduSize.valueForGivenSize(-1))
            .isInstanceOf(IllegalArgumentException.class);

    }

    /**
     * In this test the tpdu size is greater than the maximum defined by the iso tp
     * protocol spec, so it is automatically downgraded to the maximum valid value.
     */
    @Test
    @Category(FastTests.class)
    public void tpduValueForGivenTooGreatSize() {
        TpduSize tpduSize = TpduSize.valueForGivenSize(10000);

        assertThat(tpduSize).isEqualTo(TpduSize.SIZE_8192);
    }

}