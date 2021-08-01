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
package org.apache.plc4x.java.knxnetip.ets5.model;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.apache.commons.lang3.builder.ToStringStyle;

public class AddressType {

    private final String id;
    private final int mainType;
    private final int subType;
    private final String name;

    public AddressType(String id, int mainType, int subType, String name) {
        this.id = id;
        this.mainType = mainType;
        this.subType = subType;
        this.name = name;
    }

    public String getId() {
        return id;
    }

    public int getMainType() {
        return mainType;
    }

    public int getSubType() {
        return subType;
    }

    public String getName() {
        return name;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }

        if (!(o instanceof AddressType)) {
            return false;
        }

        AddressType that = (AddressType) o;

        return new EqualsBuilder()
            .append(getMainType(), that.getMainType())
            .append(getSubType(), that.getSubType())
            .append(getId(), that.getId())
            .append(getName(), that.getName())
            .isEquals();
    }

    @Override
    public int hashCode() {
        return new HashCodeBuilder(17, 37)
            .append(getId())
            .append(getMainType())
            .append(getSubType())
            .append(getName())
            .toHashCode();
    }

    @Override
    public String toString() {
        return new ToStringBuilder(this, ToStringStyle.SHORT_PREFIX_STYLE)
            .append("id", id)
            .append("mainType", mainType)
            .append("subType", subType)
            .append("name", name)
            .toString();
    }

}
