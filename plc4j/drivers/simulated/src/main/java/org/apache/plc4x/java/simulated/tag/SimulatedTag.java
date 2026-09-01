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
import org.apache.plc4x.java.simulated.readwrite.SimulatedDataTypeSizes;
import org.apache.plc4x.java.simulated.types.SimulatedTagType;
import org.apache.plc4x.java.spi.drivers.model.AddressConstraints;
import org.apache.plc4x.java.spi.drivers.model.ArrayNotationParser;
import org.apache.plc4x.java.spi.drivers.model.DefaultArrayInfo;

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
    private static final Pattern ADDRESS_PATTERN = Pattern.compile(
        "^(?<type>\\w+)/(?<name>[a-zA-Z0-9_\\\\.]+)" + ArrayNotationParser.ARRAY_GROUP + ":(?<dataType>[a-zA-Z0-9]++)$");

    /**
     * Largest amount of made-up data a single tag may ask for.
     *
     * <p>There is no protocol here to take a limit from - the device invents whatever is asked of
     * it, and the count in the address decides how much. So this is simply a size past which a
     * test is not testing anything, set well above any real use of it.</p>
     */
    static final int MAX_TAG_BYTES = 16 * 1024 * 1024;

    private final SimulatedTagType type;
    private final String name;
    private final PlcValueType dataType;
    private final int numElements;

    /**
     * Whether the address wrote the selection as a range. {@code [0]} and {@code [0..0]} both make
     * one element, and only the range is an array, so the count cannot carry this.
     */
    private final boolean explicitRange;

    private SimulatedTag(SimulatedTagType type, String name, PlcValueType dataType, int numElements) {
        this(type, name, dataType, numElements, numElements > 1);
    }

    private SimulatedTag(SimulatedTagType type, String name, PlcValueType dataType, int numElements,
                         boolean explicitRange) {
        this.type = type;
        this.name = name;
        this.dataType = dataType;
        this.numElements = numElements;
        this.explicitRange = explicitRange;
    }


    /**
     * Resolves the address's array expression to the number of elements. This driver addresses a
     * named variable rather than a numeric offset, so a selection that does not start at the
     * first element has nothing to apply to and is reported rather than quietly ignored.
     */
    /** Whether the address wrote a range. Not derivable from the count - see {@link #explicitRange}. */
    private static boolean rangeWritten(Matcher matcher, String tagString) {
        String expression = matcher.group("array");
        if (expression == null) {
            return false;
        }
        return ArrayNotationParser.parse(expression, tagString, AddressConstraints.SINGLE_DIMENSION)
            .get(0).isRange();
    }

    private static int elementsOf(Matcher matcher, String tagString) {
        String expression = matcher.group("array");
        if (expression == null) {
            return 1;
        }
        ArrayInfo dimension = ArrayNotationParser
            .parse(expression, tagString, AddressConstraints.SINGLE_DIMENSION).get(0);
        if (dimension.getLowerBound() - dimension.getBase() != 0) {
            throw new PlcInvalidTagException("Array selection '" + expression + "' in tag '"
                + tagString + "' must start at the first element: this driver addresses a "
                + "named variable, so there is no offset to start from");
        }
        return dimension.getSize();
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

            int numElements = checkNumElements(dataType, elementsOf(matcher, tagString));
            return new SimulatedTag(type, name, dataType, numElements, rangeWritten(matcher, tagString));
        }
        throw ArrayNotationParser.invalidAddress(tagString,
            "{type}/{name}[selection]:{TYPE} - for example RANDOM/foo[0..3]:INT");
    }

    /**
     * Checks that the tag asks for an amount of data that could be made.
     *
     * <p>The count is multiplied by the size of one element to get the array the device fills, and
     * that multiplication is done in an int - so a large enough count does not produce a large
     * array, it produces a negative one. Measuring the product as a long first is what keeps that
     * from being the way anyone finds out.</p>
     */
    private static int checkNumElements(PlcValueType dataType, int numElements) throws PlcInvalidTagException {
        if (numElements < 1) {
            throw new PlcInvalidTagException("The number of elements must be greater than zero.");
        }
        long totalBytes = (long) numElements * elementSizeOf(dataType);
        if (totalBytes > MAX_TAG_BYTES) {
            throw new PlcInvalidTagException("A tag of " + numElements + " elements of type "
                + dataType.name() + " would take " + totalBytes + " bytes, more than the "
                + MAX_TAG_BYTES + " a simulated tag may take.");
        }
        return numElements;
    }

    /** What one element of this type occupies, or one byte for a type with no fixed size. */
    private static int elementSizeOf(PlcValueType dataType) {
        try {
            return Math.max(1, SimulatedDataTypeSizes.valueOf(dataType.name()).getDataTypeSize());
        } catch (IllegalArgumentException e) {
            return 1;
        }
    }

    static boolean matches(String tagString) {
        return ADDRESS_PATTERN.matcher(tagString).matches();
    }

    @Override
    public String getAddressString() {
        return String.format("%s/%s%s:%s", type.name(), name,
            ArrayNotationParser.render(getArrayInfo()), dataType.name());
    }

    @Override
    public PlcValueType getPlcValueType() {
        return dataType;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        // A range is an array even when it spans one element; the count cannot express that.
        if (explicitRange) {
            return Collections.singletonList(new DefaultArrayInfo(0, numElements - 1, 0, true));
        }
        return Collections.emptyList();
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
            explicitRange == simulatedTag.explicitRange &&
            type == simulatedTag.type &&
            Objects.equals(name, simulatedTag.name) &&
            Objects.equals(dataType, simulatedTag.dataType);
    }

    @Override
    public int hashCode() {
        return Objects.hash(type, name, dataType, numElements, explicitRange);
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
