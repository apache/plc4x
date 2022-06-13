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
package org.apache.plc4x.java.bacnetip.field;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class BacNetIpField implements PlcField {

    private static final Pattern ADDRESS_PATTERN =
        Pattern.compile("^(?<deviceIdentifier>(\\d|\\*))/(?<objectType>(\\d|\\*))/(?<objectInstance>(\\d|\\*))");

    public static final int INT_WILDCARD = -1;
    public static final long LONG_WILDCARD = -1;

    private final long deviceIdentifier;
    private final int objectType;
    private final long objectInstance;

    public static boolean matches(String fieldString) {
        return ADDRESS_PATTERN.matcher(fieldString).matches();
    }

    public static BacNetIpField of(String fieldString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(fieldString);
        if(matcher.matches()) {
            long deviceIdentifier = matcher.group("deviceIdentifier").equals("*") ?
                LONG_WILDCARD : Long.parseLong(matcher.group("deviceIdentifier"));
            int objectType = matcher.group("objectType").equals("*") ?
                INT_WILDCARD : Integer.parseInt(matcher.group("objectType"));
            long objectInstance = matcher.group("objectInstance").equals("*") ?
                LONG_WILDCARD : Long.parseLong(matcher.group("objectInstance"));
            return new BacNetIpField(deviceIdentifier, objectType, objectInstance);
        }
        throw new PlcInvalidFieldException("Unable to parse address: " + fieldString);
    }

    public BacNetIpField(long deviceIdentifier, int objectType, long objectInstance) {
        this.deviceIdentifier = deviceIdentifier;
        this.objectType = objectType;
        this.objectInstance = objectInstance;
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

    public boolean matches(BacNetIpField otherField) {
        return ((deviceIdentifier == LONG_WILDCARD) || (deviceIdentifier == otherField.deviceIdentifier)) &&
            ((objectType == INT_WILDCARD) || (objectType == otherField.objectType)) &&
            ((objectInstance == LONG_WILDCARD) || (objectInstance == otherField.objectInstance));
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }

        if (!(o instanceof BacNetIpField)) {
            return false;
        }

        BacNetIpField that = (BacNetIpField) o;

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
