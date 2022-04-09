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
package org.apache.plc4x.java.firmata.readwrite.model;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.apache.plc4x.java.firmata.readwrite.field.FirmataField;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;

public class FirmataSubscriptionHandle extends DefaultPlcSubscriptionHandle {

    private final String name;
    private final FirmataField field;

    public FirmataSubscriptionHandle(PlcSubscriber plcSubscriber, String name, FirmataField field) {
        super(plcSubscriber);
        this.name = name;
        this.field = field;
    }

    public String getName() {
        return name;
    }

    public FirmataField getField() {
        return field;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }

        if (!(o instanceof FirmataSubscriptionHandle)) {
            return false;
        }

        FirmataSubscriptionHandle that = (FirmataSubscriptionHandle) o;

        return new EqualsBuilder()
            .append(getField(), that.getField())
            .isEquals();
    }

    @Override
    public int hashCode() {
        return new HashCodeBuilder(17, 37)
            .appendSuper(super.hashCode())
            .append(getField())
            .toHashCode();
    }

    @Override
    public String toString() {
        return new ToStringBuilder(this)
            .append("field", field)
            .toString();
    }

}
