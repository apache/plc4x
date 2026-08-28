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
package org.apache.plc4x.java.firmata.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.spi.drivers.model.AddressConstraints;
import org.apache.plc4x.java.spi.drivers.model.ArrayNotationParser;
import org.apache.plc4x.java.api.model.PlcTag;

import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public abstract class FirmataTag implements PlcTag {

    public static final Pattern ADDRESS_PATTERN =
        Pattern.compile("(?<address>\\d{1,3})" + ArrayNotationParser.ARRAY_GROUP);

    /**
     * Number of pins the protocol can name at all: a pin travels the wire in eight bits.
     */
    public static final int PIN_COUNT = 256;

    private final int address;

    private final int quantity;

    /**
     * Whether the address wrote the selection as a range. {@code 3[4]} and {@code 3[4..4]} both
     * select one pin, and only the range is an array, so the count cannot carry this.
     */
    private final boolean explicitRange;

    public static FirmataTag of(String tagString) {
        Matcher matcher = FirmataTagAnalog.ADDRESS_PATTERN.matcher(tagString);
        if (matcher.matches()) {
            return FirmataTagAnalog.of(tagString);
        }
        matcher = FirmataTagDigital.ADDRESS_PATTERN.matcher(tagString);
        if (matcher.matches()) {
            return FirmataTagDigital.of(tagString);
        }
        throw new PlcInvalidTagException("Unable to parse address: " + tagString);
    }

    /**
     * Resolves the address's array expression to the offset of the first pin and the number of
     * pins. Firmata addresses carry no type, so the expression terminates the address.
     *
     * @return {@code {offset, quantity, rangeWritten}} - the last is 1 or 0, and is not
     *         derivable from the others: {@code [4]} and {@code [4..4]} both select one element
     */
    protected static int[] selectionOf(Matcher matcher, String address) {
        String expression = matcher.group("array");
        if (expression == null) {
            return new int[]{0, 1, 0};
        }
        ArrayInfo dimension = ArrayNotationParser
            .parse(expression, address, AddressConstraints.SINGLE_DIMENSION).getFirst();
        return new int[]{dimension.getLowerBound() - dimension.getBase(), dimension.getSize(),
            dimension.isRange() ? 1 : 0};
    }

    protected FirmataTag(int address, Integer quantity) {
        this(address, quantity, (quantity != null) && (quantity > 1));
    }

    protected FirmataTag(int address, Integer quantity, boolean explicitRange) {
        this.explicitRange = explicitRange;
        this.address = address;
        this.quantity = quantity != null ? quantity : 1;
        if (this.quantity <= 0) {
            throw new IllegalArgumentException("quantity must be greater then zero. Was " + this.quantity);
        }
        // A digital tag turns this span into one set bit per pin, so the count decides how long
        // that loop runs and how much the BitSet holds. Refuse a span that leaves the pin space
        // before anything is built from it.
        if (address < 0 || address >= PIN_COUNT || this.quantity > PIN_COUNT - address) {
            throw new PlcInvalidTagException("A tag of " + this.quantity + " elements at address " +
                address + " reaches past the " + PIN_COUNT + " pins the protocol can name.");
        }
    }

    public int getAddress() {
        return address;
    }

    /** Whether the address wrote a range, as opposed to selecting a single pin. */
    public boolean isExplicitRange() {
        return explicitRange;
    }

    public int getNumberOfElements() {
        return quantity;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o instanceof FirmataTag that) {
            return address == that.address;
        }
        return false;
    }

    @Override
    public int hashCode() {
        return Objects.hash(address);
    }

    @Override
    public String toString() {
        return "FirmataTag{" +
            "address=" + address +
            "quantity=" + quantity +
            '}';
    }


}
