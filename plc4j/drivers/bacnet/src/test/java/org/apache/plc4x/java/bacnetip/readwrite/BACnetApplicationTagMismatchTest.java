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

package org.apache.plc4x.java.bacnetip.readwrite;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * An application tag names its own type on the wire, so the tag a field gets is whichever one the
 * data selected - not necessarily the one the field was declared with. A field that expects an
 * unsigned integer and is handed a signed one has to fail as a parse error, because a
 * ClassCastException would travel straight past the receive path's error handling.
 */
class BACnetApplicationTagMismatchTest {

    private static final WithOption[] OPTS = {
        WithByteBasedOption.WithByteOrder("BIG_ENDIAN"),
        WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
        WithOption.WithSignedIntegerEncoding("twos-complement"),
        WithOption.WithFloatEncoding("IEEE754"),
        WithOption.WithStringEncoding("UTF8")
    };

    /** An I-Am request: device identifier, max APDU length, segmentation support, vendor id. */
    private static byte[] iAmRequest(int... maxApduLengthTag) {
        byte[] head = {
            0x00,                                     // service choice I-Am
            (byte) 0xC4, 0x02, 0x00, 0x00, 0x01,      // object identifier, application tag 12
        };
        byte[] tail = {
            (byte) 0x91, 0x00,                        // segmentation supported, enumerated tag 9
            0x21, 0x0F                                // vendor id, unsigned tag 2
        };
        byte[] data = new byte[head.length + maxApduLengthTag.length + tail.length];
        System.arraycopy(head, 0, data, 0, head.length);
        for (int i = 0; i < maxApduLengthTag.length; i++) {
            data[head.length + i] = (byte) maxApduLengthTag[i];
        }
        System.arraycopy(tail, 0, data, head.length + maxApduLengthTag.length, tail.length);
        return data;
    }

    private static BACnetUnconfirmedServiceRequest parse(byte[] data) throws Exception {
        return BACnetUnconfirmedServiceRequest.staticParse(
            new ReadBufferByteBased(data, OPTS), data.length - 1);
    }

    /** Application tag 2 is an unsigned integer, which is what the field is declared with. */
    @Test
    void unsignedIntegerInTheMaxApduLengthSlotParses() throws Exception {
        // 0x22 -> unsigned integer tag, two payload bytes.
        BACnetUnconfirmedServiceRequest request = parse(iAmRequest(0x22, 0x01, 0xE0));

        BACnetUnconfirmedServiceRequestIAm iAm = assertInstanceOf(BACnetUnconfirmedServiceRequestIAm.class, request);
        assertEquals(480, iAm.getMaximumApduLengthAcceptedLength().getActualValue().longValue());
    }

    /** Application tag 3 is a signed integer, so the declared unsigned integer is not what arrives. */
    @Test
    void signedIntegerInTheMaxApduLengthSlotIsAParseFailure() {
        // 0x31 -> signed integer tag, one payload byte.
        byte[] data = iAmRequest(0x31, 0x05);

        assertThrows(BufferException.class, () -> parse(data));
    }
}
