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
package org.apache.plc4x.java.iec608705104.readwrite.model;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.apache.plc4x.java.iec608705104.readwrite.tag.Iec608705104Tag;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;

public class Iec608705104SubscriptionHandle extends DefaultPlcSubscriptionHandle {

    private final Iec608705104Tag tag;

    public Iec608705104SubscriptionHandle(PlcSubscriber plcSubscriber, Iec608705104Tag tag) {
        super(plcSubscriber);
        this.tag = tag;
    }

    public Iec608705104Tag getTag() {
        return tag;
    }

    /*public boolean matches(Iec608705104Tag address) {
        return tag.matchesGroupAddress(address);
    }*/

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }

        if (!(o instanceof Iec608705104SubscriptionHandle)) {
            return false;
        }

        Iec608705104SubscriptionHandle that = (Iec608705104SubscriptionHandle) o;

        return new EqualsBuilder()
            .append(getTag(), that.getTag())
            .isEquals();
    }

    @Override
    public int hashCode() {
        return new HashCodeBuilder(17, 37)
            .appendSuper(super.hashCode())
            .append(getTag())
            .toHashCode();
    }

    @Override
    public String toString() {
        return new ToStringBuilder(this)
            .append("tag", tag)
            .toString();
    }

}
