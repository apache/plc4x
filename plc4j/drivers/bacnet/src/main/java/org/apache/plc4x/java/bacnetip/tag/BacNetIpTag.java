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
package org.apache.plc4x.java.bacnetip.tag;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class BacNetIpTag implements PlcTag {

    private static final Pattern ADDRESS_PATTERN =
        Pattern.compile("TODO: finish me... see golang implementation");

    public static final int INT_WILDCARD = -1;
    public static final long LONG_WILDCARD = -1;

    private final long deviceIdentifier;
    private final int objectType;
    private final long objectInstance;

    public static boolean matches(String tagString) {
        return ADDRESS_PATTERN.matcher(tagString).matches();
    }

    public static BacNetIpTag of(String tagString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(tagString);
        if (matcher.matches()) {
            long deviceIdentifier = matcher.group("deviceIdentifier").equals("*") ?
                LONG_WILDCARD : Long.parseLong(matcher.group("deviceIdentifier"));
            int objectType = matcher.group("objectType").equals("*") ?
                INT_WILDCARD : Integer.parseInt(matcher.group("objectType"));
            long objectInstance = matcher.group("objectInstance").equals("*") ?
                LONG_WILDCARD : Long.parseLong(matcher.group("objectInstance"));
            return new BacNetIpTag(deviceIdentifier, objectType, objectInstance);
        }
        throw new PlcInvalidTagException("Unable to parse address: " + tagString);
    }

    public BacNetIpTag(long deviceIdentifier, int objectType, long objectInstance) {
        this.deviceIdentifier = deviceIdentifier;
        this.objectType = objectType;
        this.objectInstance = objectInstance;
    }

    @Override
    public String getAddressString() {
        return "TODO: finish me... see golang implementation";
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcTag.super.getPlcValueType();
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return PlcTag.super.getArrayInfo();
    }

    public long getDeviceIdentifier() {
        return deviceIdentifier;
    }

    public int getObjectType() {
        return objectType;
    }

    public long getObjectInstance() {
        return objectInstance;
    }

    public boolean matches(BacNetIpTag otherTag) {
        return ((deviceIdentifier == LONG_WILDCARD) || (deviceIdentifier == otherTag.deviceIdentifier)) &&
            ((objectType == INT_WILDCARD) || (objectType == otherTag.objectType)) &&
            ((objectInstance == LONG_WILDCARD) || (objectInstance == otherTag.objectInstance));
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }

        if (!(o instanceof BacNetIpTag)) {
            return false;
        }

        BacNetIpTag that = (BacNetIpTag) o;

        return new EqualsBuilder()
            .append(getDeviceIdentifier(), that.getDeviceIdentifier())
            .append(getObjectType(), that.getObjectType())
            .append(getObjectInstance(), that.getObjectInstance())
            .isEquals();
    }

    @Override
    public int hashCode() {
        return new HashCodeBuilder(17, 37)
            .append(getDeviceIdentifier())
            .append(getObjectType())
            .append(getObjectInstance())
            .toHashCode();
    }

    @Override
    public String toString() {
        return new ToStringBuilder(this)
            .append("deviceIdentifier", deviceIdentifier)
            .append("objectType", objectType)
            .append("objectInstance", objectInstance)
            .toString();
    }

}
