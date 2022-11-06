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
package org.apache.plc4x.java.ads.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.model.DefaultArrayInfo;

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
public class DirectAdsField implements AdsField {

    private static final Pattern RESOURCE_ADDRESS_PATTERN = Pattern.compile("^((0[xX](?<indexGroupHex>[0-9a-fA-F]+))|(?<indexGroup>\\d+))/((0[xX](?<indexOffsetHex>[0-9a-fA-F]+))|(?<indexOffset>\\d+)):(?<adsDataType>\\w+)(\\[(?<numberOfElements>\\d+)])?");

    private final long indexGroup;

    private final long indexOffset;

    private final String adsDataTypeName;

    private final int numberOfElements;

    public DirectAdsField(long indexGroup, long indexOffset, String adsDataTypeName, Integer numberOfElements) {
        //ByteValue.checkUnsignedBounds(indexGroup, 4);
        this.indexGroup = indexGroup;
        //ByteValue.checkUnsignedBounds(indexOffset, 4);
        this.indexOffset = indexOffset;
        this.adsDataTypeName = Objects.requireNonNull(adsDataTypeName);
        this.numberOfElements = numberOfElements != null ? numberOfElements : 1;
        if (this.numberOfElements <= 0) {
            throw new IllegalArgumentException("numberOfElements must be greater then zero. Was " + this.numberOfElements);
        }
    }

    public static DirectAdsField of(long indexGroup, long indexOffset, String adsDataTypeName, Integer numberOfElements) {
        return new DirectAdsField(indexGroup, indexOffset, adsDataTypeName, numberOfElements);
    }

    public static DirectAdsField of(String address) {
        Matcher matcher = RESOURCE_ADDRESS_PATTERN.matcher(address);
        if (!matcher.matches()) {
            throw new PlcInvalidFieldException(address, RESOURCE_ADDRESS_PATTERN, "{indexGroup}/{indexOffset}:{adsDataType}([numberOfElements])?");
        }

        String indexGroupStringHex = matcher.group("indexGroupHex");
        String indexGroupString = matcher.group("indexGroup");

        String indexOffsetStringHex = matcher.group("indexOffsetHex");
        String indexOffsetString = matcher.group("indexOffset");

        long indexGroup;
        if (indexGroupStringHex != null) {
            indexGroup = Long.parseLong(indexGroupStringHex, 16);
        } else {
            indexGroup = Long.parseLong(indexGroupString);
        }

        long indexOffset;
        if (indexOffsetStringHex != null) {
            indexOffset = Long.parseLong(indexOffsetStringHex, 16);
        } else {
            indexOffset = Long.parseLong(indexOffsetString);
        }

        String adsDataTypeString = matcher.group("adsDataType");

        String numberOfElementsString = matcher.group("numberOfElements");
        Integer numberOfElements = numberOfElementsString != null ? Integer.valueOf(numberOfElementsString) : null;

        return new DirectAdsField(indexGroup, indexOffset, adsDataTypeString, numberOfElements);
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
        if (!(o instanceof DirectAdsField)) {
            return false;
        }
        DirectAdsField that = (DirectAdsField) o;
        return indexGroup == that.indexGroup &&
            indexOffset == that.indexOffset;
    }

    @Override
    public int hashCode() {
        return Objects.hash(indexGroup, indexOffset);
    }

    @Override
    public String toString() {
        return "DirectAdsField{" +
            "indexGroup=" + indexGroup +
            ", indexOffset=" + indexOffset +
            '}';
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());

        writeBuffer.writeUnsignedLong("indexGroup", 32, getIndexGroup());
        writeBuffer.writeUnsignedLong("indexOffset", 32, getIndexOffset());
        writeBuffer.writeUnsignedLong("numberOfElements", 32, getNumberOfElements());
        String plcDataType = getPlcDataType();
        writeBuffer.writeString("dataType", plcDataType.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), plcDataType);

        writeBuffer.popContext(getClass().getSimpleName());
    }

}
