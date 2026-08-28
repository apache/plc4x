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

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?<address>[\\w\\-. ]+)" + ArrayNotationParser.ARRAY_GROUP + "(:(?<datatype>[a-zA-Z_]+)){1}");
    private final String address;
    private final int quantity;

    /**
     * Whether the address wrote the selection as a range. A one-element range is still a range,
     * which no count can express.
     */
    private final boolean explicitRange;
    private final PlcValueType dataType;

    protected ProfinetTag(String address, Integer quantity, PlcValueType dataType) {
        this(address, quantity, dataType, (quantity != null) && (quantity > 1));
    }

    protected ProfinetTag(String address, Integer quantity, PlcValueType dataType, boolean explicitRange) {
        this.explicitRange = explicitRange;
        this.address = address;
        this.quantity = (quantity != null) ? quantity : 1;
        if (this.quantity <= 0) {
            throw new IllegalArgumentException("quantity must be greater than zero. Was " + this.quantity);
        }
        this.dataType = dataType;
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
                "{address}[selection]:{TYPE} - for example foo.bar[0..3]:DINT");
        }

        int quantity = elementsOf(matcher, addressString);
        PlcValueType plcValueType = PlcValueType.valueOf(matcher.group("datatype"));

        return new ProfinetTag(matcher.group("address"), quantity, plcValueType,
            rangeWritten(matcher, addressString));
    }

    @Override
    public String getAddressString() {
        // Unchanged from before the notation migration: the type is not part of what this
        // driver reports, so the result does not round-trip. Only the selection is added.
        return address + ArrayNotationParser.render(getArrayInfo());
    }

    @Override
    public PlcValueType getPlcValueType() {
        return dataType;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        // A range is an array even when it spans one element; the count cannot express that.
        if (explicitRange) {
            return Collections.singletonList(new DefaultArrayInfo(0, quantity - 1, 0, true));
        }
        return Collections.emptyList();
    }
}
