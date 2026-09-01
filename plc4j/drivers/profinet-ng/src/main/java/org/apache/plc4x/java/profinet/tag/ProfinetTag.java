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
package org.apache.plc4x.java.profinet.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.spi.drivers.model.AddressConstraints;
import org.apache.plc4x.java.spi.drivers.model.DefaultArrayInfo;
import org.apache.plc4x.java.spi.drivers.model.ArrayNotationParser;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.Collections;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ProfinetTag implements PlcTag {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?<slot>\\d{1,5}).(?<subSlot>\\d{1,5}).(?<direction>INPUT|OUTPUT)(.(?<index>\\d{1,5}))?" + ArrayNotationParser.ARRAY_GROUP + ":(?<dataType>[a-zA-Z_]+)");
    private final int slot;
    private final int subSlot;
    private final Direction direction;
    private final int index;
    private final PlcValueType dataType;
    private final int numElements;

    /**
     * Whether the address wrote the selection as a range. A one-element range is still a range,
     * which no count can express.
     */
    private final boolean explicitRange;

    public ProfinetTag(int slot, int subSlot, Direction direction, int index, PlcValueType dataType, int numElements) {
        this(slot, subSlot, direction, index, dataType, numElements, numElements > 1);
    }

    public ProfinetTag(int slot, int subSlot, Direction direction, int index, PlcValueType dataType, int numElements, boolean explicitRange) {
        this.explicitRange = explicitRange;
        this.slot = slot;
        this.subSlot = subSlot;
        this.direction = direction;
        this.index = index;
        this.dataType = dataType;
        this.numElements = numElements;
        if (this.numElements <= 0) {
            throw new IllegalArgumentException("numElements must be greater than zero. Was " + this.numElements);
        }
    }


    /**
     * Resolves the address's array expression to the number of elements. This driver addresses a
     * named variable rather than a numeric offset, so a selection that does not start at the
     * first element has nothing to apply to and is reported rather than quietly ignored.
     */
    /**
     * Whether the address wrote a range. Not derivable from the element count: {@code [4]} and
     * {@code [4..4]} both select one element, and only the range is an array.
     */
    private static boolean rangeWritten(Matcher matcher, String addressString) {
        String expression = matcher.group("array");
        if (expression == null) {
            return false;
        }
        return ArrayNotationParser.parse(expression, addressString, AddressConstraints.SINGLE_DIMENSION)
            .get(0).isRange();
    }

    private static int elementsOf(Matcher matcher, String addressString) {
        String expression = matcher.group("array");
        if (expression == null) {
            return 1;
        }
        ArrayInfo dimension = ArrayNotationParser
            .parse(expression, addressString, AddressConstraints.SINGLE_DIMENSION).get(0);
        if (dimension.getLowerBound() - dimension.getBase() != 0) {
            throw new PlcInvalidTagException("Array selection '" + expression + "' in tag '"
                + addressString + "' must start at the first element: this driver addresses a "
                + "named variable, so there is no offset to start from");
        }
        return dimension.getSize();
    }

    public static ProfinetTag of(String addressString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (!matcher.matches()) {
            throw ArrayNotationParser.invalidAddress(addressString,
                "{slot}.{subSlot}.{INPUT|OUTPUT}.{index}[selection]:{TYPE}"
                    + " - for example 1.2.INPUT.0[0..3]:INT");
        }

        int slot = Integer.parseInt(matcher.group("slot"));
        int subSlot = Integer.parseInt(matcher.group("subSlot"));
        Direction direction = Direction.valueOf(matcher.group("direction"));
        // The index is optional in the pattern but has always been required in practice - it was
        // parsed unguarded, so an address without one failed with a NumberFormatException. Say so.
        String indexToken = matcher.group("index");
        if (indexToken == null) {
            throw new PlcInvalidTagException("Address '" + addressString + "' is missing the index:"
                + " expected {slot}.{subSlot}.{INPUT|OUTPUT}.{index}[selection]:{TYPE}");
        }
        int index = Integer.parseInt(indexToken);
        PlcValueType dataType = PlcValueType.valueOf(matcher.group("dataType"));
        int numElements = elementsOf(matcher, addressString);

        return new ProfinetTag(slot, subSlot, direction, index, dataType, numElements,
            rangeWritten(matcher, addressString));
    }

    public int getSlot() {
        return slot;
    }

    public int getSubSlot() {
        return subSlot;
    }

    public Direction getDirection() {
        return direction;
    }

    public int getIndex() {
        return index;
    }

    public int getNumElements() {
        return numElements;
    }

    @Override
    public String getAddressString() {
        return String.format("%d.%d.%s.%d:%s%s", slot, subSlot, direction, index, dataType, (numElements > 1) ? "[" + numElements + "]" : "");
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

    public static enum Direction {
        INPUT,
        OUTPUT
    }

}
