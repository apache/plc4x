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
package org.apache.plc4x.java.simulated.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.simulated.types.SimulatedTagType;
import org.apache.plc4x.java.spi.model.DefaultArrayInfo;

import java.util.Collections;
import java.util.List;
import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Test address for accessing values in virtual devices.
 */
public class SimulatedTag implements PlcTag {

    /**
     * Examples:
     * - {@code RANDOM/foo:INTEGER}
     * - {@code STDOUT/foo:STRING}
     */
    private static final Pattern ADDRESS_PATTERN = Pattern.compile("^(?<type>\\w+)/(?<name>[a-zA-Z0-9_\\\\.]+):(?<dataType>[a-zA-Z0-9]++)(\\[(?<numElements>\\d+)])?$");

    private final SimulatedTagType type;
    private final String name;
    private final PlcValueType dataType;
    private final int numElements;

    private SimulatedTag(SimulatedTagType type, String name, PlcValueType dataType, int numElements) {
        this.type = type;
        this.name = name;
        this.dataType = dataType;
        this.numElements = numElements;
    }

    public static SimulatedTag of(String tagString) throws PlcInvalidTagException {
        Matcher matcher = ADDRESS_PATTERN.matcher(tagString);
        if (matcher.matches()) {
            SimulatedTagType type = SimulatedTagType.valueOf(matcher.group("type"));
            String name = matcher.group("name");

            PlcValueType dataType;
            try {
                dataType = PlcValueType.valueOf(matcher.group("dataType").toUpperCase());
            } catch (Exception e) {
                throw new PlcInvalidTagException("Invalid data type: " + matcher.group("dataType"));
            }

            int numElements = 1;
            if (matcher.group("numElements") != null) {
                numElements = Integer.parseInt(matcher.group("numElements"));
            }
            return new SimulatedTag(type, name, dataType, numElements);
        }
        throw new PlcInvalidTagException("Unable to parse address: " + tagString);
    }

    static boolean matches(String tagString) {
        return ADDRESS_PATTERN.matcher(tagString).matches();
    }

    @Override
    public String getAddressString() {
        return String.format("%s/%s:%s[%d]", type.name(), name, dataType.name(), numElements);
    }

    @Override
    public PlcValueType getPlcValueType() {
        return dataType;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return Collections.singletonList(new DefaultArrayInfo(0, numElements));
    }

    public SimulatedTagType getType() {
        return type;
    }

    public String getName() {
        return name;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        SimulatedTag simulatedTag = (SimulatedTag) o;
        return numElements == simulatedTag.numElements &&
            type == simulatedTag.type &&
            Objects.equals(name, simulatedTag.name) &&
            Objects.equals(dataType, simulatedTag.dataType);
    }

    @Override
    public int hashCode() {
        return Objects.hash(type, name, dataType, numElements);
    }

    @Override
    public String toString() {
        return "SimulatedTag{" +
            "type=" + type +
            ", name='" + name + '\'' +
            ", dataType='" + dataType + '\'' +
            ", numElements=" + numElements +
            '}';
    }

}
