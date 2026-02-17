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
import org.apache.plc4x.java.opcua.readwrite.ExtensionObjectDefinition;
import org.apache.plc4x.java.opcua.readwrite.UnknownExtensionObject;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

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
            // Non-numeric NodeIds (e.g., vendor-specific types like Siemens TE_DTL)
            // are not known to the parser. Return 0 to signal an unknown type.
            return 0;
        }
    }

    public static ExtensionObjectDefinition parseExtensionObjectBody(
            ReadBuffer readBuffer, int extensionId, int bodyLength) throws ParseException {
        // Maintain the same "body" context that a simple field would produce,
        // so the XML roundtrip tests stay consistent.
        readBuffer.pullContext("body");
        ExtensionObjectDefinition result;
        if (extensionId < 1 && bodyLength > 0) {
            // Unknown extension object type (e.g., vendor-specific like Siemens TE_DTL).
            // Read the raw body bytes so the buffer position stays correct.
            byte[] rawBytes = readBuffer.readByteArray("unknownBody", bodyLength);
            result = new UnknownExtensionObject(rawBytes);
        } else {
            result = ExtensionObjectDefinition.staticParse(readBuffer, extensionId);
        }
        readBuffer.closeContext("body");
        return result;
    }

    public static void serializeExtensionObjectBody(
            WriteBuffer writeBuffer, ExtensionObjectDefinition body) throws SerializationException {
        writeBuffer.pushContext("body");
        body.serialize(writeBuffer);
        writeBuffer.popContext("body");
    }
}
