/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.ads.model;

import org.apache.plc4x.java.ads.api.util.ByteValue;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;

import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * ADS address witch is defined by {@code indexGroup/indexOffset}. These values can be either supplied as int or hex
 * representation.
 */
public class DirectAdsField implements AdsField {

    private static final Pattern RESOURCE_ADDRESS_PATTERN = Pattern.compile("^((0[xX](?<indexGroupHex>[0-9a-fA-F]+))|(?<indexGroup>\\d+))/((0[xX](?<indexOffsetHex>[0-9a-fA-F]+))|(?<indexOffset>\\d+)):(?<adsDataType>\\w+)(\\[(?<numberOfElements>\\d)])?");

    private final long indexGroup;

    private final long indexOffset;

    private final AdsDataType adsDataType;

    private final int numberOfElements;

    private DirectAdsField(long indexGroup, long indexOffset, AdsDataType adsDataType, Integer numberOfElements) {
        ByteValue.checkUnsignedBounds(indexGroup, 4);
        this.indexGroup = indexGroup;
        ByteValue.checkUnsignedBounds(indexOffset, 4);
        this.indexOffset = indexOffset;
        this.adsDataType = Objects.requireNonNull(adsDataType);
        this.numberOfElements = numberOfElements != null ? numberOfElements : 1;
        if (this.numberOfElements <= 0) {
            throw new IllegalArgumentException("numberOfElements must be greater then zero. Was " + this.numberOfElements);
        }
    }

    public static DirectAdsField of(long indexGroup, long indexOffset, AdsDataType adsDataType, Integer numberOfElements) {
        return new DirectAdsField(indexGroup, indexOffset, adsDataType, numberOfElements);
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
        AdsDataType adsDataType = AdsDataType.valueOf(adsDataTypeString);

        String numberOfElementsString = matcher.group("numberOfElements");
        Integer numberOfElements = numberOfElementsString != null ? Integer.valueOf(numberOfElementsString) : null;

        return new DirectAdsField(indexGroup, indexOffset, adsDataType, numberOfElements);
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

    @Override
    public AdsDataType getAdsDataType() {
        return adsDataType;
    }

    public int getNumberOfElements() {
        return numberOfElements;
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
}
