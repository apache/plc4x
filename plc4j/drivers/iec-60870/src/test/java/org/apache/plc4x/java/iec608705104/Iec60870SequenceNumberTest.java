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

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * An acknowledgement tells the station how far we have read, so it can release what it was holding
 * in case it had to send it again. Getting the number wrong is not a local matter: the station acts
 * on what we send it.
 */
class Iec60870SequenceNumberTest {

    @Test
    void whatWeExpectNextIsOnePastWhatArrived() {
        assertEquals(1, Iec60870Connection.nextExpectedAfter(0));
        assertEquals(8, Iec60870Connection.nextExpectedAfter(7));
        assertEquals(32767, Iec60870Connection.nextExpectedAfter(32766));
    }

    @Test
    void theSequenceWrapsRatherThanGrowing() {
        // Fifteen bits wide: past the top it starts again, it does not become a sixteen-bit number.
        assertEquals(0, Iec60870Connection.nextExpectedAfter(32767));
        assertEquals(1, Iec60870Connection.nextExpectedAfter(0 + 32768));
    }

    @Test
    void theNumberGoesOnTheWireShiftedUpByOne() {
        // The lowest bit of the control field is what marks the format, so it must stay clear.
        assertEquals(0b10, Iec60870Connection.encodeSequenceNumber(1));
        assertEquals(0b1110, Iec60870Connection.encodeSequenceNumber(7));
        assertEquals(0xFFFE, Iec60870Connection.encodeSequenceNumber(32767));
    }

    @Test
    void anEncodedSequenceNumberAlwaysFitsTheControlField() {
        for (int sequenceNumber = 0; sequenceNumber < 32768; sequenceNumber++) {
            int encoded = Iec60870Connection.encodeSequenceNumber(sequenceNumber);
            assertTrue(encoded >= 0 && encoded <= 0xFFFE,
                sequenceNumber + " encoded to " + encoded + ", which does not fit the field");
            assertEquals(0, encoded & 1,
                sequenceNumber + " encoded with its lowest bit set, which would change the format");
        }
    }

    @Test
    void theWholeCycleRoundTrips() {
        // Every number the station can send, acknowledged and decoded back the way the peer will.
        for (int sent = 0; sent < 32768; sent++) {
            int encoded = Iec60870Connection.encodeSequenceNumber(Iec60870Connection.nextExpectedAfter(sent));
            assertEquals((sent + 1) % 32768, encoded >> 1);
        }
    }
}
