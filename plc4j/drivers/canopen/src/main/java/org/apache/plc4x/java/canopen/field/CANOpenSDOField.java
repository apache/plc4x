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
package org.apache.plc4x.java.canopen.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.canopen.readwrite.types.CANOpenDataType;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.w3c.dom.Document;
import org.w3c.dom.Element;

import java.nio.charset.StandardCharsets;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class CANOpenSDOField extends CANOpenField {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("SDO:" + NODE_PATTERN + "(?:/(?<answerNodeId>\\d+))?:" + CANOpenField.ADDRESS_PATTERN);
    private final int answerNode;
    private final short index;
    private final short subIndex;
    private final CANOpenDataType canOpenDataType;

    public CANOpenSDOField(int node, short index, short subIndex, CANOpenDataType canOpenDataType) {
        this(node, node, index, subIndex, canOpenDataType);
    }

    public CANOpenSDOField(int node, int answerNode, short index, short subIndex, CANOpenDataType canOpenDataType) {
        super(node);
        this.answerNode = answerNode;
        this.index = index;
        this.subIndex = subIndex;
        this.canOpenDataType = canOpenDataType;
    }

    public int getAnswerNodeId() {
        return answerNode;
    }

    public short getIndex() {
        return index;
    }

    public short getSubIndex() {
        return subIndex;
    }

    public CANOpenDataType getCanOpenDataType() {
        return canOpenDataType;
    }

    public static boolean matches(String addressString) {
        return ADDRESS_PATTERN.matcher(addressString).matches();
    }

    public static Matcher getMatcher(String addressString) throws PlcInvalidFieldException {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (matcher.matches()) {
            return matcher;
        }

        throw new PlcInvalidFieldException(addressString, ADDRESS_PATTERN);
    }

    public static CANOpenSDOField of(String addressString) {
        Matcher matcher = getMatcher(addressString);
        int nodeId = Integer.parseInt(matcher.group("nodeId"));

        short index = parseHex(matcher.group("indexHex"), matcher.group("index"));
        short subIndex = parseHex(matcher.group("subIndexHex"), matcher.group("subIndex"));

        String canDataTypeString = matcher.group("canDataType");
        CANOpenDataType canOpenDataType = CANOpenDataType.valueOf(canDataTypeString);

        String answerNode = matcher.group("answerNodeId");
        if (answerNode != null) {
            int answerNodeId = Integer.parseInt(matcher.group("answerNodeId"));
            return new CANOpenSDOField(nodeId, answerNodeId, index, subIndex, canOpenDataType);
        }

        //String numberOfElementsString = matcher.group("numberOfElements");
        //Integer numberOfElements = numberOfElementsString != null ? Integer.valueOf(numberOfElementsString) : null;

        return new CANOpenSDOField(nodeId, index, subIndex, canOpenDataType);
    }

    private static Short parseHex(String hex, String dec) {
        if (hex != null) {
            return Short.parseShort(hex, 16);
        }
        return Short.parseShort(dec);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws ParseException {
        writeBuffer.pushContext(getClass().getSimpleName());

        writeBuffer.writeInt("node",64, getNodeId());
        writeBuffer.writeInt("index",64, getIndex());
        writeBuffer.writeInt("subIndex",64, getSubIndex());
        String dataTypeName = getCanOpenDataType().name();
        writeBuffer.writeString("dataType", dataTypeName.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), dataTypeName);

        writeBuffer.popContext(getClass().getSimpleName());
    }

}
