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
package org.apache.plc4x.java.ads.tag;

import org.apache.plc4x.java.ads.readwrite.AdsDataType;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.drivers.model.AddressConstraints;
import org.apache.plc4x.java.spi.drivers.model.ArrayNotationParser;
import org.apache.plc4x.java.spi.drivers.model.DefaultArrayInfo;

import java.nio.charset.StandardCharsets;
import java.util.Collections;
import java.util.List;
import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * ADS address witch is defined by {@code indexGroup/indexOffset}. These values can be either supplied as int or hex
 * representation.
 */
public class DirectAdsTag implements AdsTag {

    private static final Pattern RESOURCE_ADDRESS_PATTERN = Pattern.compile(
        "^((0[xX](?<indexGroupHex>[0-9a-fA-F]{1,8}))|(?<indexGroup>\\d{1,10}))" +
            "/((0[xX](?<indexOffsetHex>[0-9a-fA-F]{1,8}))|(?<indexOffset>\\d{1,10}))" +
            ArrayNotationParser.ARRAY_GROUP +
            ":(?<adsDataType>\\w+)");

    /** An index group, an index offset and a length all travel ADS as four bytes. */
    private static final long MAX_UINT32 = 0xFFFFFFFFL;

    private final long indexGroup;

    private final long indexOffset;

    private final String adsDataTypeName;

    private final int numberOfElements;

    /**
     * Whether the address wrote the selection as a range. A one-element range is still a range,
     * and the count cannot say which was written.
     */
    private final boolean explicitRange;

    public DirectAdsTag(long indexGroup, long indexOffset, String adsDataTypeName, Integer numberOfElements) {
        this(indexGroup, indexOffset, adsDataTypeName, numberOfElements, (numberOfElements != null) && (numberOfElements > 1));
    }

    public DirectAdsTag(long indexGroup, long indexOffset, String adsDataTypeName, Integer numberOfElements, boolean explicitRange) {
        this.explicitRange = explicitRange;
        this.indexGroup = checkUint32("indexGroup", indexGroup);
        this.indexOffset = checkUint32("indexOffset", indexOffset);
        this.adsDataTypeName = Objects.requireNonNull(adsDataTypeName);
        this.numberOfElements = numberOfElements != null ? numberOfElements : 1;
        if (this.numberOfElements <= 0) {
            throw new IllegalArgumentException("numberOfElements must be greater then zero. Was " + this.numberOfElements);
        }
    }

    /**
     * An address that cannot be put on the wire is not an address. ADS carries each of these as
     * four bytes, so a value past that would be silently narrowed on serialisation and address
     * something other than what was asked for.
     */
    private static long checkUint32(String what, long value) {
        if (value < 0 || value > MAX_UINT32) {
            throw new PlcInvalidTagException("The " + what + " " + value +
                " does not fit the four bytes ADS carries it in.");
        }
        return value;
    }

    /**
     * The count is multiplied by the size of one element to give the request's length, so it has
     * to be a count that could be counted, rather than something narrowed into a negative and
     * then reported as a different complaint. The pattern caps how many digits reach here and
     * every value it admits fits a long, so the range is all that is left to check.
     */
    protected static int parseElementCount(String count) {
        long parsed = Long.parseLong(count);
        if (parsed < 1 || parsed > Integer.MAX_VALUE) {
            throw new PlcInvalidTagException("The number of elements " + parsed +
                " is not a count of elements that could be read.");
        }
        return (int) parsed;
    }

    /**
     * Reads one of the address' numbers. The pattern admits at most ten decimal or eight hex
     * digits, and all of those fit a long, so a number too wide to parse never reaches here - it
     * fails to match the address at all and is reported as the invalid address it is.
     */
    protected static long parseUint32(String what, String decimal, String hex) {
        return checkUint32(what, hex != null ? Long.parseLong(hex, 16) : Long.parseLong(decimal));
    }

    public static DirectAdsTag of(long indexGroup, long indexOffset, String adsDataTypeName, Integer numberOfElements) {
        return new DirectAdsTag(indexGroup, indexOffset, adsDataTypeName, numberOfElements);
    }

    /**
     * Resolves the address's array expression to the offset from the index offset and the number
     * of elements. An absent expression selects one element at the address itself.
     *
     * @return {@code {offset, numberOfElements}}
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

    public static DirectAdsTag of(String address) {
        Matcher matcher = RESOURCE_ADDRESS_PATTERN.matcher(address);
        if (!matcher.matches()) {
            throw ArrayNotationParser.invalidAddress(address,
                "{indexGroup}/{indexOffset}[selection]:{TYPE} - for example 0x4020/0[0..3]:DINT");
        }

        String indexGroupStringHex = matcher.group("indexGroupHex");
        String indexGroupString = matcher.group("indexGroup");

        String indexOffsetStringHex = matcher.group("indexOffsetHex");
        String indexOffsetString = matcher.group("indexOffset");

        long indexGroup = parseUint32("indexGroup", indexGroupString, indexGroupStringHex);
        long indexOffset = parseUint32("indexOffset", indexOffsetString, indexOffsetStringHex);

        String adsDataTypeString = matcher.group("adsDataType");

        int[] selection = selectionOf(matcher, address);
        // An index offset is a byte offset; the selection counts elements. They are the same
        // number only for a one-byte type, so 0x4020/0[3]:DINT would otherwise advance three
        // bytes and read from inside the first element.
        indexOffset += (long) selection[0] * bytesPerElement(adsDataTypeString, selection[0], address);
        Integer numberOfElements = selection[1];

        return new DirectAdsTag(indexGroup, indexOffset, adsDataTypeString, numberOfElements,
            selection[2] == 1);
    }

    /**
     * The storage size of one element of the named type.
     *
     * <p>The device's data-type table is not available while an address is being parsed, so only
     * the types ADS defines itself can be measured here. A selection on anything else cannot be
     * placed, and is refused rather than silently applied at the wrong offset - the offset is only
     * needed when something was selected, so an address without a selection is unaffected.</p>
     */
    private static int bytesPerElement(String typeName, int offset, String address) {
        try {
            return AdsDataType.valueOf(typeName).getNumBytes();
        } catch (IllegalArgumentException e) {
            if (offset == 0) {
                return 1;
            }
            throw new PlcInvalidTagException("Cannot place a selection in '" + address + "': the size"
                + " of type '" + typeName + "' is only known to the device, so the element's offset"
                + " cannot be computed here. Address the element directly instead.");
        }
    }

    public static boolean matches(String address) {
        return RESOURCE_ADDRESS_PATTERN.matcher(address).matches();
    }

    public long getIndexGroup() {
        return indexGroup;
    }

    public long getIndexOffset() {
        return indexOffset;
    }

    public String getPlcDataType() {
        return adsDataTypeName;
    }

    public int getNumberOfElements() {
        return numberOfElements;
    }

    @Override
    public String getAddressString() {
        // "0x%d" printed the group's decimal digits behind a hex prefix, so group 16416 came back
        // as 0x16416 - which re-parses as 91158, a different address entirely.
        return String.format("0x%X/%d%s:%s", getIndexGroup(), getIndexOffset(),
            ArrayNotationParser.render(getArrayInfo()), getPlcDataType());
    }

    @Override
    public PlcValueType getPlcValueType() {
        try {
            return PlcValueType.valueOf(adsDataTypeName);
        } catch (Exception e) {
            return PlcValueType.Struct;
        }
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        // A range is an array even when it spans one element; the count cannot express that.
        if (explicitRange) {
            return Collections.singletonList(new DefaultArrayInfo(0, getNumberOfElements() - 1, 0, true));
        }
        return Collections.emptyList();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o instanceof DirectAdsTag that) {
            return indexGroup == that.indexGroup &&
                indexOffset == that.indexOffset;
        }
        return false;
    }

    @Override
    public int hashCode() {
        return Objects.hash(indexGroup, indexOffset);
    }

    @Override
    public String toString() {
        return "DirectAdsTag{" +
            "indexGroup=" + indexGroup +
            ", indexOffset=" + indexOffset +
            '}';
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        writeBuffer.pushContext(WithOption.WithName(getClass().getSimpleName()));

        writeBuffer.writeUnsignedLong(32, getIndexGroup(), WithOption.WithName("indexGroup"));
        writeBuffer.writeUnsignedLong(32, getIndexOffset(), WithOption.WithName("indexOffset"));
        writeBuffer.writeUnsignedLong(32, getNumberOfElements(), WithOption.WithName("numberOfElements"));
        writeBuffer.writeString(
            getPlcDataType().getBytes(StandardCharsets.UTF_8).length * 8,
            getPlcDataType(),
            WithOption.WithName("dataType"),
            WithOption.WithEncoding("UTF8"));

        writeBuffer.popContext(WithOption.WithName(getClass().getSimpleName()));
    }

}
