/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, (byte) Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, (byte) WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, (byte) either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.opcua.readwrite.utils;

import java.nio.charset.StandardCharsets;
import org.apache.plc4x.java.opcua.readwrite.ExpandedNodeId;
import org.apache.plc4x.java.opcua.readwrite.NodeIdByteString;
import org.apache.plc4x.java.opcua.readwrite.NodeIdFourByte;
import org.apache.plc4x.java.opcua.readwrite.NodeIdGuid;
import org.apache.plc4x.java.opcua.readwrite.NodeIdNumeric;
import org.apache.plc4x.java.opcua.readwrite.NodeIdString;
import org.apache.plc4x.java.opcua.readwrite.NodeIdTwoByte;
import org.apache.plc4x.java.opcua.readwrite.NodeIdTypeDefinition;

public class StaticHelper {

    // Calculating length in UTF-8
    public static int utf8LengthToPascalLength(String stringValue) {
        if (stringValue == null) {
            return -1;
        }
        return stringValue.getBytes(StandardCharsets.UTF_8).length;
    }

    public static int pascalLengthToUtf8Length(int slength) {
        return Math.max(slength, 0);
    }

    public static int extensionId(ExpandedNodeId expandedNodeId) {
        try {
            return Integer.parseInt(expandedNodeId.getNodeId().getIdentifier());
        } catch (NumberFormatException e) {
            // Non-numeric encoding NodeId (string/guid/opaque). Well-known standard types always
            // use numeric ids, so this is a custom/user-defined type whose body is captured as raw
            // bytes (bodyKind = 2) and never dispatched through ExtensionObjectDefinition — the
            // extensionId is unused there. Return a sentinel that matches no known type instead of
            // failing the whole parse.
            return -1;
        }
    }

    /**
     * Whether the encoding NodeId of an ExtensionObject refers to a well-known OPC UA standard type,
     * i.e. it lives in namespace 0 and is not referenced by a namespace URI. Custom / user-defined
     * structure encodings live in namespace &ge; 1; their bodies cannot be parsed by the generated
     * {@code ExtensionObjectDefinition} dispatch and are instead captured as raw bytes (which the
     * driver decodes against the type's StructureDefinition).
     */
    public static boolean isStandardEncoding(ExpandedNodeId expandedNodeId) {
        if (expandedNodeId == null) {
            return true;
        }
        // A namespace referenced by URI is, by definition, not the standard OPC UA namespace (0).
        if (expandedNodeId.getNamespaceURISpecified()) {
            return false;
        }
        NodeIdTypeDefinition nodeId = expandedNodeId.getNodeId();
        if (nodeId instanceof NodeIdTwoByte) {
            // The two-byte form is always namespace 0.
            return true;
        } else if (nodeId instanceof NodeIdFourByte fourByte) {
            return fourByte.getNamespaceIndex() == 0;
        } else if (nodeId instanceof NodeIdNumeric numeric) {
            return numeric.getNamespaceIndex() == 0;
        } else if (nodeId instanceof NodeIdString string) {
            return string.getNamespaceIndex() == 0;
        } else if (nodeId instanceof NodeIdGuid guid) {
            return guid.getNamespaceIndex() == 0;
        } else if (nodeId instanceof NodeIdByteString byteString) {
            return byteString.getNamespaceIndex() == 0;
        }
        return false;
    }
}
