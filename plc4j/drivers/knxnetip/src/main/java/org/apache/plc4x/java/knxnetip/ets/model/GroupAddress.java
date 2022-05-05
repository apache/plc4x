/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
import org.apache.plc4x.java.knxnetip.readwrite.KnxDatapointType;

public class GroupAddress {

    private final String groupAddress;
    private final String name;
    private final KnxDatapointType type;
    private final Function function;

    public GroupAddress(String groupAddress, String name, KnxDatapointType type, Function function) {
        this.groupAddress = groupAddress;
        this.name = name;
        this.type = type;
        this.function = function;
    }

    public String getGroupAddress() {
        return groupAddress;
    }

    public String getName() {
        return name;
    }

    public KnxDatapointType getType() {
        return type;
    }

    public Function getFunction() {
        return function;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }

        if (!(o instanceof GroupAddress)) {
            return false;
        }

        GroupAddress that = (GroupAddress) o;

        return new EqualsBuilder()
            .append(getGroupAddress(), that.getGroupAddress())
            .append(getName(), that.getName())
            .append(getType(), that.getType())
            .append(getFunction(), that.getFunction())
            .isEquals();
    }

    @Override
    public int hashCode() {
        return new HashCodeBuilder(17, 37)
            .append(getGroupAddress())
            .append(getName())
            .append(getType())
            .append(getFunction())
            .toHashCode();
    }

    @Override
    public String toString() {
        return new ToStringBuilder(this, ToStringStyle.SHORT_PREFIX_STYLE)
            .append("groupAddress", groupAddress)
            .append("name", name)
            .append("type", type)
            .append("function", function)
            .toString();
    }

}
