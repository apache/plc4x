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

package org.apache.plc4x.java.opcua.protocol.chunk;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.util.Collections;
import org.apache.plc4x.java.opcua.readwrite.BinaryPayload;
import org.apache.plc4x.java.opcua.readwrite.ExpandedNodeId;
import org.apache.plc4x.java.opcua.readwrite.ExtensiblePayload;
import org.apache.plc4x.java.opcua.readwrite.ExtensionObject;
import org.apache.plc4x.java.opcua.readwrite.HistoryEvent;
import org.apache.plc4x.java.opcua.readwrite.NodeIdFourByte;
import org.apache.plc4x.java.opcua.readwrite.RootExtensionObject;
import org.apache.plc4x.java.opcua.readwrite.SequenceHeader;
import org.apache.plc4x.java.spi.utils.hex.Hex;
import org.junit.jupiter.api.Test;

class PayloadConverterTest {

    @Test
    void convert() throws Exception {
        ExpandedNodeId expandedNodeId = new ExpandedNodeId(
            false,           //Namespace Uri Specified
            false,            //Server Index Specified
            new NodeIdFourByte(
                (short) 0, 661
            ),
            null,
            null
        );

        RootExtensionObject extObject = new RootExtensionObject(
            expandedNodeId,
            new HistoryEvent(Collections.emptyList())
        );

        ExtensiblePayload payload = new ExtensiblePayload(
            new SequenceHeader(1, 2),
            extObject
        );

        BinaryPayload binary = PayloadConverter.toBinary(payload);
        ExtensiblePayload extensible = PayloadConverter.toExtensible(binary);

        String extensibleSrcHex = Hex.dump(PayloadConverter.toStream(payload));
        String binaryDstHex = Hex.dump(PayloadConverter.toStream(binary));
        String extensibleDstHex = Hex.dump(PayloadConverter.toStream(extensible));

        assertEquals(extensibleSrcHex, binaryDstHex);
        assertEquals(extensibleSrcHex, extensibleDstHex);
    }

}