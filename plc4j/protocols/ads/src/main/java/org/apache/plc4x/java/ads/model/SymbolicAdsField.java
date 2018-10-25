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

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;

import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * ADS address witch is defined by symbolic name (e.g. {@code Main.items[0]}).
 */
public class SymbolicAdsField implements AdsField {

    private static final Pattern SYMBOLIC_ADDRESS_PATTERN = Pattern.compile("^(?<symbolicAddress>.+):(?<adsDataType>\\w+)(\\[(?<numberOfElements>\\d)])?");

    private final String symbolicAddress;

    private final AdsDataType adsDataType;

    private final int numberOfElements;

    private SymbolicAdsField(String symbolicAddress, AdsDataType adsDataType, Integer numberOfElements) {
        this.symbolicAddress = Objects.requireNonNull(symbolicAddress);
        this.adsDataType = Objects.requireNonNull(adsDataType);
        this.numberOfElements = numberOfElements != null ? numberOfElements : 1;
        if (this.numberOfElements <= 0) {
            throw new IllegalArgumentException("numberOfElements must be greater then zero. Was " + this.numberOfElements);
        }
    }

    public static SymbolicAdsField of(String address) {
        Matcher matcher = SYMBOLIC_ADDRESS_PATTERN.matcher(address);
        if (!matcher.matches()) {
            throw new PlcInvalidFieldException(address, SYMBOLIC_ADDRESS_PATTERN, "{address}");
        }
        String symbolicAddress = matcher.group("symbolicAddress");

        String adsDataTypeString = matcher.group("adsDataType");
        AdsDataType adsDataType = AdsDataType.valueOf(adsDataTypeString);

        String numberOfElementsString = matcher.group("numberOfElements");
        Integer numberOfElements = numberOfElementsString != null ? Integer.valueOf(numberOfElementsString) : null;

        return new SymbolicAdsField(symbolicAddress, adsDataType, numberOfElements);
    }

    public static boolean matches(String address) {
        return SYMBOLIC_ADDRESS_PATTERN.matcher(address).matches();
    }

    public String getSymbolicField() {
        return symbolicAddress;
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
        if (!(o instanceof SymbolicAdsField)) {
            return false;
        }
        SymbolicAdsField that = (SymbolicAdsField) o;
        return Objects.equals(symbolicAddress, that.symbolicAddress);
    }

    @Override
    public int hashCode() {
        return Objects.hash(symbolicAddress);
    }

    @Override
    public String toString() {
        return "SymbolicAdsField{" +
            "symbolicAddress='" + symbolicAddress + '\'' +
            '}';
    }
}
