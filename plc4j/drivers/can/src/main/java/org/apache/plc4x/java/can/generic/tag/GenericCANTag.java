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
package org.apache.plc4x.java.can.generic.tag;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.genericcan.readwrite.GenericCANDataType;
import org.apache.plc4x.java.spi.model.DefaultArrayInfo;

public class GenericCANTag implements PlcTag {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?<nodeId>\\d+):(?<dataType>\\w+)(?:\\[(?<arraySize>\\d+)\\])?");
    private final int nodeId;
    private final GenericCANDataType dataType;
    private final int arraySize;

    public GenericCANTag(int nodeId, GenericCANDataType dataType, int arraySize) {
        this.nodeId = nodeId;
        this.dataType = dataType;
        this.arraySize = arraySize;
    }

    public int getNodeId() {
        return nodeId;
    }

    @Override
    public String getAddressString() {
        String address = nodeId + ":" + dataType.name();
        if(arraySize != 1) {
            address += "[" + arraySize + "]";
        }
        return address;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.valueOf(dataType.getPlcValueName());
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        if(arraySize != 1) {
            return Collections.singletonList(new DefaultArrayInfo(0, arraySize));
        }
        return Collections.emptyList();
    }

    public int getArraySize() {
        return arraySize;
    }

    public GenericCANDataType getDataType() {
        return dataType;
    }

    public static Optional<GenericCANTag> matches(String tagQuery) {
        Matcher matcher = ADDRESS_PATTERN.matcher(tagQuery);
        return matcher.matches() ? Optional.of(GenericCANTag.create(matcher)) : Optional.empty();
    }

    static GenericCANTag create(Matcher tagQuery) {
        int nodeId = Integer.parseInt(tagQuery.group("nodeId"));
        String type = tagQuery.group("dataType");

        GenericCANDataType dataType;
        try {
             dataType = GenericCANDataType.valueOf(type);
        } catch (IllegalArgumentException e) {
            throw new PlcRuntimeException("Could not create tag with data type " + type, e);
        }
        int arraySize = tagQuery.group("arraySize") != null ? Integer.parseInt(tagQuery.group("arraySize")) : 0;

        return new GenericCANTag(nodeId, dataType, arraySize);
    }

    public String toString() {
        return "GenericCANTag(" + nodeId + ":" + dataType.name() + (arraySize == 0 ? "" : "[" + arraySize + "]");
    }

}
