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

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
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
            ":(?<adsDataType>\\w+)(\\[(?<numberOfElements>\\d{1,10})])?");

    /** An index group, an index offset and a length all travel ADS as four bytes. */
    private static final long MAX_UINT32 = 0xFFFFFFFFL;

    private final long indexGroup;

    private final long indexOffset;

    private final String adsDataTypeName;

    private final int numberOfElements;

    public DirectAdsTag(long indexGroup, long indexOffset, String adsDataTypeName, Integer numberOfElements) {
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
     * to be a count that could be counted. Anything wider is refused as an invalid tag rather
     * than narrowed into a negative and reported as something else.
     */
    protected static int parseElementCount(String count) {
        try {
            long parsed = Long.parseLong(count);
            if (parsed < 1 || parsed > Integer.MAX_VALUE) {
                throw new PlcInvalidTagException("The number of elements " + parsed +
                    " is not a count of elements that could be read.");
            }
            return (int) parsed;
        } catch (NumberFormatException e) {
            throw new PlcInvalidTagException("The number of elements is not a number.", e);
        }
    }

    /**
     * Reads one of the address' numbers, refusing a value too wide for ADS rather than letting
     * it arrive as an unchecked NumberFormatException from somewhere inside the parse.
     */
    protected static long parseUint32(String what, String decimal, String hex) {
        try {
            return checkUint32(what, hex != null ? Long.parseLong(hex, 16) : Long.parseLong(decimal));
        } catch (NumberFormatException e) {
            throw new PlcInvalidTagException("The " + what + " is not a number ADS could carry.", e);
        }
    }

    public static DirectAdsTag of(long indexGroup, long indexOffset, String adsDataTypeName, Integer numberOfElements) {
        return new DirectAdsTag(indexGroup, indexOffset, adsDataTypeName, numberOfElements);
    }

    public static DirectAdsTag of(String address) {
        Matcher matcher = RESOURCE_ADDRESS_PATTERN.matcher(address);
        if (!matcher.matches()) {
            throw new PlcInvalidTagException(address, RESOURCE_ADDRESS_PATTERN, "{indexGroup}/{indexOffset}:{adsDataType}([numberOfElements])?");
        }

        String indexGroupStringHex = matcher.group("indexGroupHex");
        String indexGroupString = matcher.group("indexGroup");

        String indexOffsetStringHex = matcher.group("indexOffsetHex");
        String indexOffsetString = matcher.group("indexOffset");

        long indexGroup = parseUint32("indexGroup", indexGroupString, indexGroupStringHex);
        long indexOffset = parseUint32("indexOffset", indexOffsetString, indexOffsetStringHex);

        String adsDataTypeString = matcher.group("adsDataType");

        String numberOfElementsString = matcher.group("numberOfElements");
        Integer numberOfElements = numberOfElementsString != null
            ? parseElementCount(numberOfElementsString) : null;

        return new DirectAdsTag(indexGroup, indexOffset, adsDataTypeString, numberOfElements);
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
        String address = String.format("0x%d/%d:%s", getIndexGroup(), getIndexOffset(), getPlcDataType());
        if(getNumberOfElements() != 1) {
            address += "[" + getNumberOfElements() + "]";
        }
        return address;
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
        if(getNumberOfElements() != 1) {
            return Collections.singletonList(new DefaultArrayInfo(0, getNumberOfElements()));
        }
        return Collections.emptyList();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof DirectAdsTag)) {
            return false;
        }
        DirectAdsTag that = (DirectAdsTag) o;
        return indexGroup == that.indexGroup &&
            indexOffset == that.indexOffset;
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
