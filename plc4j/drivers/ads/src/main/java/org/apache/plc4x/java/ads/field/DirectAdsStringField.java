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
package org.apache.plc4x.java.ads.field;

import org.apache.plc4x.java.ads.readwrite.types.AdsDataType;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.w3c.dom.Document;
import org.w3c.dom.Element;

import java.nio.charset.StandardCharsets;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * ADS address witch is defined by {@code indexGroup/indexOffset}. These values can be either supplied as int or hex
 * representation.
 */
public class DirectAdsStringField extends DirectAdsField implements AdsStringField {

    private static final Pattern RESOURCE_STRING_ADDRESS_PATTERN = Pattern.compile("^((0[xX](?<indexGroupHex>[0-9a-fA-F]+))|(?<indexGroup>\\d+))/((0[xX](?<indexOffsetHex>[0-9a-fA-F]+))|(?<indexOffset>\\d+)):(?<adsDataType>STRING|WSTRING)\\((?<stringLength>\\d{1,3})\\)(\\[(?<numberOfElements>\\d+)])?");

    private final int stringLength;

    public DirectAdsStringField(long indexGroup, long indexOffset, AdsDataType adsDataType, int stringLength, Integer numberOfElements) {
        super(indexGroup, indexOffset, adsDataType, numberOfElements);
        this.stringLength = stringLength;
    }

    public static DirectAdsStringField of(long indexGroup, long indexOffset, AdsDataType adsDataType, int stringLength, Integer numberOfElements) {
        return new DirectAdsStringField(indexGroup, indexOffset, adsDataType, stringLength, numberOfElements);
    }

    public static DirectAdsStringField of(String address) {
        Matcher matcher = RESOURCE_STRING_ADDRESS_PATTERN.matcher(address);
        if (!matcher.matches()) {
            throw new PlcInvalidFieldException(address, RESOURCE_STRING_ADDRESS_PATTERN, "{indexGroup}/{indexOffset}:{adsDataType}([numberOfElements])?");
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

        String stringLengthString = matcher.group("stringLength");
        Integer stringLength = stringLengthString != null ? Integer.valueOf(stringLengthString) : null;

        String numberOfElementsString = matcher.group("numberOfElements");
        Integer numberOfElements = numberOfElementsString != null ? Integer.valueOf(numberOfElementsString) : null;

        return new DirectAdsStringField(indexGroup, indexOffset, adsDataType, stringLength, numberOfElements);
    }

    public static boolean matches(String address) {
        return RESOURCE_STRING_ADDRESS_PATTERN.matcher(address).matches();
    }

    @Override
    public int getStringLength() {
        return stringLength;
    }

    @Override
    public String toString() {
        return "DirectAdsStringField{" +
            "indexGroup=" + getIndexGroup() +
            ", indexOffset=" + getIndexOffset() +
            ", stringLength=" + stringLength +
            '}';
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws ParseException {
        writeBuffer.pushContext(getClass().getSimpleName());

        writeBuffer.writeUnsignedLong("indexGroup", 32, getIndexGroup());
        writeBuffer.writeUnsignedLong("indexOffset", 32, getIndexOffset());
        writeBuffer.writeUnsignedLong("numberOfElements", 32, getNumberOfElements());
        String plcDataType = getPlcDataType();
        writeBuffer.writeString("dataType", plcDataType.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), plcDataType);

        writeBuffer.popContext(getClass().getSimpleName());
    }

}
