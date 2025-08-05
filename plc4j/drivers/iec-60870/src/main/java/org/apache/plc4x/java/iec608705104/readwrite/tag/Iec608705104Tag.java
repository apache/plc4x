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

package org.apache.plc4x.java.iec608705104.readwrite.tag;

import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.List;

public class Iec608705104Tag implements PlcTag {

    /*private static final Pattern IEC_60870_5_104_ADDRESS =
        Pattern.compile("^(?<adsuAddrSingleValueGroup>(\\d|\\*))|((?<adsuAddrDoubleValueGroup1>(\\d|\\*))/(?<adsuAddrDoubleValueGroup2>(\\d|\\*)))/((?<objectAddressSingleValueGroup(\\d|\\*))|((?<objectAddressTrippleValueGroup1>(\\d|\\*))\\.(?<objectAddressTrippleValueGroup2>(\\d|\\*))\\.(?<objectAddressTrippleValueGroup3>(\\d|\\*))))");
     */

    private final int adsuAddress;
    private final int objectAddress;

    /*public static boolean matches(String tagString) {
        return IEC_60870_5_104_ADDRESS.matcher(tagString).matches();
    }*/

    public Iec608705104Tag(int adsuAddress, int objectAddress) {
        this.adsuAddress = adsuAddress;
        this.objectAddress = objectAddress;
    }

    public int getAdsuAddress() {
        return adsuAddress;
    }

    public int getObjectAddress() {
        return objectAddress;
    }

    @Override
    public String getAddressString() {
        return null;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.NULL;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return PlcTag.super.getArrayInfo();
    }

    @Override
    public String toString() {
        return "Iec608705104Tag{" +
            "adsuAddress=" + adsuAddress +
            ", objectAddress=" + objectAddress +
            '}';
    }

}
