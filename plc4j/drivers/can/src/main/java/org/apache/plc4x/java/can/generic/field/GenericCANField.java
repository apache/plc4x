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
package org.apache.plc4x.java.can.generic.field;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;

import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.apache.plc4x.java.genericcan.readwrite.GenericCANDataType;

public class GenericCANField implements PlcField {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?<nodeId>\\d+):(?<dataType>\\w+)(?:\\[(?<arraySize>\\d+)\\])?");
    private final int nodeId;
    private final GenericCANDataType dataType;
    private final int arraySize;

    public GenericCANField(int nodeId, GenericCANDataType dataType, int arraySize) {
        this.nodeId = nodeId;
        this.dataType = dataType;
        this.arraySize = arraySize;
    }

    public int getNodeId() {
        return nodeId;
    }

    @Override
    public String getPlcDataType() {
        return dataType.name();
    }

    public int getArraySize() {
        return arraySize;
    }

    public GenericCANDataType getDataType() {
        return dataType;
    }

    public static Optional<GenericCANField> matches(String fieldQuery) {
        Matcher matcher = ADDRESS_PATTERN.matcher(fieldQuery);
        return matcher.matches() ? Optional.ofNullable(GenericCANField.create(matcher)) : Optional.empty();
    }

    static GenericCANField create(Matcher fieldQuery) {
        int nodeId = Integer.parseInt(fieldQuery.group("nodeId"));
        String type = fieldQuery.group("dataType");

        GenericCANDataType dataType;
        try {
             dataType = GenericCANDataType.valueOf(type);
        } catch (IllegalArgumentException e) {
            throw new PlcRuntimeException("Could not create field with data type " + type, e);
        }
        int arraySize = fieldQuery.group("arraySize") != null ? Integer.parseInt(fieldQuery.group("arraySize")) : 0;

        return new GenericCANField(nodeId, dataType, arraySize);
    }

    public String toString() {
        return "GenericCANField(" + nodeId + ":" + dataType.name() + (arraySize == 0 ? "" : "[" + arraySize + "]");
    }

}
