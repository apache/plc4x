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
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;

import java.nio.charset.StandardCharsets;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * ADS address witch is defined by {@code indexGroup/indexOffset}. These values can be either supplied as int or hex
 * representation.
 */
public class DirectAdsStringTag extends DirectAdsTag implements AdsStringTag {

    private static final Pattern RESOURCE_STRING_ADDRESS_PATTERN = Pattern.compile("^((0[xX](?<indexGroupHex>[0-9a-fA-F]{1,8}))|(?<indexGroup>\\d{1,10}))" +
            "/((0[xX](?<indexOffsetHex>[0-9a-fA-F]{1,8}))|(?<indexOffset>\\d{1,10}))" +
            ":(?<adsDataType>STRING|WSTRING)\\((?<stringLength>\\d{1,3})\\)" +
            "(\\[(?<numberOfElements>\\d{1,10})])?");

    private final int stringLength;

    public DirectAdsStringTag(long indexGroup, long indexOffset, String adsDataTypeName, int stringLength, Integer numberOfElements) {
        super(indexGroup, indexOffset, adsDataTypeName, numberOfElements);
        this.stringLength = stringLength;
    }

    public static DirectAdsStringTag of(long indexGroup, long indexOffset, String adsDataTypeName, int stringLength, Integer numberOfElements) {
        return new DirectAdsStringTag(indexGroup, indexOffset, adsDataTypeName, stringLength, numberOfElements);
    }

    public static DirectAdsStringTag of(String address) {
        Matcher matcher = RESOURCE_STRING_ADDRESS_PATTERN.matcher(address);
        if (!matcher.matches()) {
            throw new PlcInvalidTagException(address, RESOURCE_STRING_ADDRESS_PATTERN, "{indexGroup}/{indexOffset}:{adsDataType}([numberOfElements])?");
        }

        String indexGroupStringHex = matcher.group("indexGroupHex");
        String indexGroupString = matcher.group("indexGroup");

        String indexOffsetStringHex = matcher.group("indexOffsetHex");
        String indexOffsetString = matcher.group("indexOffset");

        long indexGroup = parseUint32("indexGroup", indexGroupString, indexGroupStringHex);
        long indexOffset = parseUint32("indexOffset", indexOffsetString, indexOffsetStringHex);

        String adsDataTypeName = matcher.group("adsDataType");

        String stringLengthString = matcher.group("stringLength");
        int stringLength = stringLengthString != null ? Integer.parseInt(stringLengthString) : 0;

        String numberOfElementsString = matcher.group("numberOfElements");
        Integer numberOfElements = numberOfElementsString != null
            ? parseElementCount(numberOfElementsString) : null;

        return new DirectAdsStringTag(indexGroup, indexOffset, adsDataTypeName, stringLength, numberOfElements);
    }

    public static boolean matches(String address) {
        return RESOURCE_STRING_ADDRESS_PATTERN.matcher(address).matches();
    }

    @Override
    public String getAddressString() {
        String address = String.format("0x%d/%d:%s(%d)", getIndexGroup(), getIndexOffset(), getPlcDataType(), getStringLength());
        if(getNumberOfElements() != 1) {
            address += "[" + getNumberOfElements() + "]";
        }
        return address;
    }

    @Override
    public int getStringLength() {
        return stringLength;
    }

    @Override
    public String toString() {
        return "DirectAdsStringTag{" +
            "indexGroup=" + getIndexGroup() +
            ", indexOffset=" + getIndexOffset() +
            ", stringLength=" + stringLength +
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
        writeBuffer.writeUnsignedLong(32, getStringLength(), WithOption.WithName("stringLength"));

        writeBuffer.popContext(WithOption.WithName(getClass().getSimpleName()));
    }

}
