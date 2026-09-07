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
package org.apache.plc4x.java.profinet.device;

import org.apache.plc4x.java.profinet.readwrite.PnDcp_Pdu_AlarmLow;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.mockito.Mockito.mock;

/**
 * An alarm frame says how long its variable part is, so a frame claiming to be an alarm can carry
 * nothing at all. The reason code used to be read out of it without asking whether it was there.
 */
class ProfinetAlarmVarPartTest {

    private static PnDcp_Pdu_AlarmLow alarmWithVarPartOf(int length) {
        return new PnDcp_Pdu_AlarmLow(0x8001, 0, 0, (byte) 1, (byte) 0, (byte) 0, (byte) 0,
            0, 0, new byte[length]);
    }

    private static ProfinetDevice device() {
        return new ProfinetDevice(mock(MessageWrapper.class), "test-device", "", "", (a, b) -> null);
    }

    @Test
    void anAlarmCarryingNoReasonIsNotFatal() {
        ProfinetDevice device = device();
        // varPartLen 0, which is what the frame in the report declared, and every length short of
        // the byte the reason lives in.
        for (int length : new int[]{0, 1, 2, 3}) {
            assertDoesNotThrow(() -> device.handleAlarmResponse(alarmWithVarPartOf(length)),
                "an alarm declaring " + length + " bytes of variable part must not throw");
        }
    }

    @Test
    void anAlarmCarryingAReasonIsStillRead() {
        ProfinetDevice device = device();
        assertDoesNotThrow(() -> device.handleAlarmResponse(alarmWithVarPartOf(4)));
        assertDoesNotThrow(() -> device.handleAlarmResponse(alarmWithVarPartOf(16)));
    }
}
