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
package org.apache.plc4x.java.knxnetip.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class KnxNetIpField implements PlcField {

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

    public static boolean matches(String fieldString) {
        return KNX_GROUP_ADDRESS_3_LEVEL.matcher(fieldString).matches() ||
            KNX_GROUP_ADDRESS_2_LEVEL.matcher(fieldString).matches() ||
            KNX_GROUP_ADDRESS_1_LEVEL.matcher(fieldString).matches();
    }

    public static KnxNetIpField of(String fieldString) {
        Matcher matcher = KNX_GROUP_ADDRESS_3_LEVEL.matcher(fieldString);
        if(matcher.matches()) {
            return new KnxNetIpField(3, matcher.group("mainGroup"), null, null);
        }
        matcher = KNX_GROUP_ADDRESS_2_LEVEL.matcher(fieldString);
        if(matcher.matches()) {
            return new KnxNetIpField(2, matcher.group("mainGroup"), null, matcher.group("subGroup"));
        }
        matcher = KNX_GROUP_ADDRESS_1_LEVEL.matcher(fieldString);
        if(matcher.matches()) {
            return new KnxNetIpField(1, matcher.group("mainGroup"), matcher.group("middleGroup"), matcher.group("subGroup"));
        }
        throw new PlcInvalidFieldException("Unable to parse address: " + fieldString);
    }

    public KnxNetIpField(int levels, String mainGroup, String middleGroup, String subGroup) {
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

}
