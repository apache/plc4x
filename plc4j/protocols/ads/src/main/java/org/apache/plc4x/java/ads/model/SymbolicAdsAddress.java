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

import org.apache.plc4x.java.api.model.Address;

import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class SymbolicAdsAddress implements Address {
    private static final Pattern SYMBOLIC_ADDRESS_PATTERN = Pattern.compile("^(?<symbolicAddress>.+)");

    private final String symbolicAddress;

    private SymbolicAdsAddress(String symbolicAddress) {
        this.symbolicAddress = Objects.requireNonNull(symbolicAddress);
    }

    public static SymbolicAdsAddress of(String address) {
        Matcher matcher = SYMBOLIC_ADDRESS_PATTERN.matcher(address);
        if (!matcher.matches()) {
            throw new IllegalArgumentException(
                "address " + address + " doesn't match '{address}' RAW:" + SYMBOLIC_ADDRESS_PATTERN);
        }
        String symbolicAddress = matcher.group("symbolicAddress");

        return new SymbolicAdsAddress(symbolicAddress);
    }

    public static boolean matches(String address) {
        return SYMBOLIC_ADDRESS_PATTERN.matcher(address).matches();
    }

    public String getSymbolicAddress() {
        return symbolicAddress;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof SymbolicAdsAddress)) return false;
        SymbolicAdsAddress that = (SymbolicAdsAddress) o;
        return Objects.equals(symbolicAddress, that.symbolicAddress);
    }

    @Override
    public int hashCode() {
        return Objects.hash(symbolicAddress);
    }

    @Override
    public String toString() {
        return "SymbolicAdsAddress{" +
            "symbolicAddress='" + symbolicAddress + '\'' +
            '}';
    }
}
