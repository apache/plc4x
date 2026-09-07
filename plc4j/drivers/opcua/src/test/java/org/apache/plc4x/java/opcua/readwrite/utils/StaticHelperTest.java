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
package org.apache.plc4x.java.opcua.readwrite.utils;

import org.apache.plc4x.java.opcua.readwrite.ExpandedNodeId;
import org.apache.plc4x.java.opcua.readwrite.NodeIdByteString;
import org.apache.plc4x.java.opcua.readwrite.NodeIdFourByte;
import org.apache.plc4x.java.opcua.readwrite.NodeIdGuid;
import org.apache.plc4x.java.opcua.readwrite.NodeIdNumeric;
import org.apache.plc4x.java.opcua.readwrite.NodeIdString;
import org.apache.plc4x.java.opcua.readwrite.NodeIdTwoByte;
import org.apache.plc4x.java.opcua.readwrite.NodeIdTypeDefinition;
import org.apache.plc4x.java.opcua.readwrite.PascalByteString;
import org.apache.plc4x.java.opcua.readwrite.PascalString;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class StaticHelperTest {

    private static ExpandedNodeId expanded(NodeIdTypeDefinition nodeId) {
        return new ExpandedNodeId(false, false, nodeId, null, null);
    }

    @Test
    void isStandardEncodingHandlesNull() {
        assertThat(StaticHelper.isStandardEncoding(null)).isTrue();
    }

    @Test
    void namespaceUriMakesItCustom() {
        // A node addressed by namespace URI is never the standard OPC UA namespace, regardless of id.
        ExpandedNodeId nodeId = new ExpandedNodeId(true, false,
            new NodeIdNumeric(0, 1L), new PascalString("urn:example:custom"), null);
        assertThat(StaticHelper.isStandardEncoding(nodeId)).isFalse();
    }

    @Test
    void twoByteIsAlwaysStandard() {
        assertThat(StaticHelper.isStandardEncoding(expanded(new NodeIdTwoByte((short) 5)))).isTrue();
    }

    @Test
    void fourByteStandardOnlyInNamespaceZero() {
        assertThat(StaticHelper.isStandardEncoding(expanded(new NodeIdFourByte((short) 0, 5)))).isTrue();
        assertThat(StaticHelper.isStandardEncoding(expanded(new NodeIdFourByte((short) 1, 5)))).isFalse();
    }

    @Test
    void numericStandardOnlyInNamespaceZero() {
        assertThat(StaticHelper.isStandardEncoding(expanded(new NodeIdNumeric(0, 101L)))).isTrue();
        assertThat(StaticHelper.isStandardEncoding(expanded(new NodeIdNumeric(2, 5001L)))).isFalse();
    }

    @Test
    void stringStandardOnlyInNamespaceZero() {
        assertThat(StaticHelper.isStandardEncoding(expanded(new NodeIdString(0, new PascalString("Foo"))))).isTrue();
        assertThat(StaticHelper.isStandardEncoding(expanded(new NodeIdString(3, new PascalString("Foo"))))).isFalse();
    }

    @Test
    void guidStandardOnlyInNamespaceZero() {
        assertThat(StaticHelper.isStandardEncoding(expanded(new NodeIdGuid(0, new byte[16])))).isTrue();
        assertThat(StaticHelper.isStandardEncoding(expanded(new NodeIdGuid(1, new byte[16])))).isFalse();
    }

    @Test
    void byteStringStandardOnlyInNamespaceZero() {
        assertThat(StaticHelper.isStandardEncoding(
            expanded(new NodeIdByteString(0, new PascalByteString(0, new byte[0]))))).isTrue();
        assertThat(StaticHelper.isStandardEncoding(
            expanded(new NodeIdByteString(1, new PascalByteString(0, new byte[0]))))).isFalse();
    }

    @Test
    void extensionIdParsesNumericIdAndFallsBackForNonNumeric() {
        assertThat(StaticHelper.extensionId(expanded(new NodeIdNumeric(0, 101L)))).isEqualTo(101);
        // Non-numeric (string) encoding ids — used by custom types — must yield the -1 sentinel
        // rather than throwing.
        assertThat(StaticHelper.extensionId(expanded(new NodeIdString(3, new PascalString("Foo"))))).isEqualTo(-1);
    }
}
