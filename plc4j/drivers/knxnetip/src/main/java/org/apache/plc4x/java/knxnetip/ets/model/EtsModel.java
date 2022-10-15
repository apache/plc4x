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
package org.apache.plc4x.java.knxnetip.ets.model;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.apache.commons.lang3.builder.ToStringStyle;

import java.util.Map;

public class EtsModel {

    private final byte groupAddressType;
    private final Map<String, GroupAddress> groupAddresses;
    private final Map<String, String> topologyNames;

    public EtsModel(byte groupAddressType, Map<String, GroupAddress> groupAddresses, Map<String, String> topologyNames) {
        this.groupAddressType = groupAddressType;
        this.groupAddresses = groupAddresses;
        this.topologyNames = topologyNames;
    }

    public byte getGroupAddressType() {
        return groupAddressType;
    }

    public static String parseGroupAddress(byte groupAddressType, byte[] addressBytes) {
        int addressInt = (addressBytes[0] << 8) | (addressBytes[1] & 0xFF);
        return parseGroupAddress(groupAddressType, addressInt);
    }

    public static String parseGroupAddress(byte groupAddressType, int addressInt) {
        switch (groupAddressType) {
            case 1: {
                return Integer.toString(addressInt);
            }
            case 2: {
                int mainGroup = (addressInt & 0xF800) >> 11;
                int subGroup = (addressInt & 0x07FF);
                return mainGroup + "/" + subGroup;
            }
            case 3: {
                int mainGroup = (addressInt & 0xF800) >> 11;
                int middleGroup = (addressInt & 0x0700) >> 8;
                int subGroup = (addressInt & 0x00FF);
                return mainGroup + "/" + middleGroup + "/" + subGroup;
            }
        }
        return null;
    }

    public String parseGroupAddress(byte[] addressBytes) {
        int addressInt = (addressBytes[0] << 8) | (addressBytes[1] & 0xFF);
        return parseGroupAddress(getGroupAddressType(), addressInt);
    }

    public String parseGroupAddress(int addressInt) {
        return parseGroupAddress(getGroupAddressType(), addressInt);
    }

    public Map<String, GroupAddress> getGroupAddresses() {
        return groupAddresses;
    }

    public String getTopologyName(String addressPart) {
        return topologyNames.get(addressPart);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }

        if (!(o instanceof EtsModel)) {
            return false;
        }

        EtsModel etsModel = (EtsModel) o;

        return new EqualsBuilder()
            .append(getGroupAddressType(), etsModel.getGroupAddressType())
            .append(getGroupAddresses(), etsModel.getGroupAddresses())
            .isEquals();
    }

    @Override
    public int hashCode() {
        return new HashCodeBuilder(17, 37)
            .append(getGroupAddressType())
            .append(getGroupAddresses())
            .toHashCode();
    }

    @Override
    public String toString() {
        return new ToStringBuilder(this, ToStringStyle.SHORT_PREFIX_STYLE)
            .append("groupAddressNumLevels", groupAddressType)
            .append("groupAddresses", groupAddresses)
            .toString();
    }

}
