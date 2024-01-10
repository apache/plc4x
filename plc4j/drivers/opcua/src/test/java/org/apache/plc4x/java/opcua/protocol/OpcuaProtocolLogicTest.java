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

package org.apache.plc4x.java.opcua.protocol;

import static org.junit.jupiter.api.Assertions.assertArrayEquals;
import static org.junit.jupiter.api.Assertions.assertEquals;

import org.apache.plc4x.java.opcua.readwrite.NodeId;
import org.apache.plc4x.java.opcua.readwrite.NodeIdGuid;
import org.apache.plc4x.java.opcua.readwrite.NodeIdType;
import org.apache.plc4x.java.opcua.tag.OpcuaTag;
import org.junit.jupiter.api.Test;

public class OpcuaProtocolLogicTest {

    @Test
    public void testGenerateNodeId() {
        OpcuaTag tag = OpcuaTag.of("ns=2;g=00112233-4455-6677-8899-aabbccddeeff");
        NodeId nodeId = OpcuaProtocolLogic.generateNodeId(tag);
        assertEquals(NodeIdType.nodeIdTypeGuid, nodeId.getNodeId().getNodeType());

        NodeIdGuid nodeIdGuid = (NodeIdGuid) nodeId.getNodeId();
        assertEquals(2, nodeIdGuid.getNamespaceIndex());
        assertArrayEquals(
                new byte[] {
                        (byte) 0x33,
                        (byte) 0x22,
                        (byte) 0x11,
                        (byte) 0x00,
                        (byte) 0x55,
                        (byte) 0x44,
                        (byte) 0x77,
                        (byte) 0x66,
                        (byte) 0x88,
                        (byte) 0x99,
                        (byte) 0xaa,
                        (byte) 0xbb,
                        (byte) 0xcc,
                        (byte) 0xdd,
                        (byte) 0xee,
                        (byte) 0xff
                },
                nodeIdGuid.getId()
        );
    }
}
