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
package org.apache.plc4x.java.knxnetip.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.knxnetip.ets.model.GroupAddress;

import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class KnxNetIpTag implements PlcTag {

    private static final String WILDCARD = "*";
    private static final Pattern KNX_GROUP_ADDRESS_1_LEVEL =
        Pattern.compile("^(?<mainGroup>(\\d{1,5}|\\*))");
    private static final Pattern KNX_GROUP_ADDRESS_2_LEVEL =
        Pattern.compile("^(?<mainGroup>(\\d{1,2}|\\*))/(?<subGroup>(\\d{1,4}|\\*))");
    private static final Pattern KNX_GROUP_ADDRESS_3_LEVEL =
        Pattern.compile("^(?<mainGroup>(\\d{1,2}|\\*))/(?<middleGroup>(\\d|\\*))/(?<subGroup>(\\d{1,3}|\\*))");

    private final int levels;
    private final String mainGroup;
    private final String middleGroup;
    private final String subGroup;

    public static boolean matches(String tagString) {
        return KNX_GROUP_ADDRESS_3_LEVEL.matcher(tagString).matches() ||
            KNX_GROUP_ADDRESS_2_LEVEL.matcher(tagString).matches() ||
            KNX_GROUP_ADDRESS_1_LEVEL.matcher(tagString).matches();
    }

    public static KnxNetIpTag of(String tagString) {
        Matcher matcher = KNX_GROUP_ADDRESS_1_LEVEL.matcher(tagString);
        if(matcher.matches()) {
            return new KnxNetIpTag(1, matcher.group("mainGroup"), null, null);
        }
        matcher = KNX_GROUP_ADDRESS_2_LEVEL.matcher(tagString);
        if(matcher.matches()) {
            return new KnxNetIpTag(2, matcher.group("mainGroup"), null, matcher.group("subGroup"));
        }
        matcher = KNX_GROUP_ADDRESS_3_LEVEL.matcher(tagString);
        if(matcher.matches()) {
            return new KnxNetIpTag(3, matcher.group("mainGroup"), matcher.group("middleGroup"), matcher.group("subGroup"));
        }
        throw new PlcInvalidTagException("Unable to parse address: " + tagString);
    }

    public KnxNetIpTag(int levels, String mainGroup, String middleGroup, String subGroup) {
        this.levels = levels;
        this.mainGroup = mainGroup;
        this.middleGroup = middleGroup;
        this.subGroup = subGroup;
    }

    public int getLevels() {
        return levels;
    }

    public String getMainGroup() {
        return mainGroup;
    }

    public String getMiddleGroup() {
        return middleGroup;
    }

    public String getSubGroup() {
        return subGroup;
    }


    @Override
    public String getAddressString() {
        String address = mainGroup;
        if(middleGroup != null) {
            address += "/" + middleGroup;
        }
        if(subGroup != null) {
            address += "/" + subGroup;
        }
        return address;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.NULL;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return PlcTag.super.getArrayInfo();
    }

    // As our fields can contain wildcards and complex matching logic,
    // do a check if a given GroupAddress is actually compatible with this field.
    public boolean matchesGroupAddress(GroupAddress groupAddress) {
        KnxNetIpTag otherAddress = KnxNetIpTag.of(groupAddress.getGroupAddress());
        // If the levels don't match the whole address can't match.
        if(otherAddress.getLevels() != getLevels()) {
            return false;
        }
        // NOTE: This case fallthrough is intentional :-)
        if(getLevels() == 3) {
            if (!WILDCARD.equals(getMiddleGroup()) && !getMiddleGroup().equals(otherAddress.getMiddleGroup())) {
                return false;
            }
        }
        if(getLevels() >= 2) {
            if (!WILDCARD.equals(getSubGroup()) && !getSubGroup().equals(otherAddress.getSubGroup())) {
                return false;
            }
        }
        return WILDCARD.equals(getMainGroup()) || getMainGroup().equals(otherAddress.getMainGroup());
    }

}
