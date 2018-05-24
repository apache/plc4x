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
package org.apache.plc4x.java.api.messages;

import org.apache.plc4x.java.api.model.Address;

import java.util.Date;
import java.util.List;
import java.util.Objects;

public class PlcNotification<T> {

    protected final Date timeStamp;

    protected final Address address;

    protected final List<T> values;

    public PlcNotification(Date timeStamp, Address address, List<T> values) {
        this.timeStamp = timeStamp;
        this.address = address;
        this.values = values;
    }

    public Date getTimeStamp() {
        return timeStamp;
    }

    public Address getAddress() {
        return address;
    }

    public List<T> getValues() {
        return values;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof PlcNotification)) {
            return false;
        }
        PlcNotification<?> that = (PlcNotification<?>) o;
        return Objects.equals(timeStamp, that.timeStamp) &&
            Objects.equals(address, that.address) &&
            Objects.equals(values, that.values);
    }

    @Override
    public int hashCode() {
        return Objects.hash(timeStamp, address, values);
    }

    @Override
    public String toString() {
        return "PlcNotification{" +
            "timeStamp=" + timeStamp +
            ", address=" + address +
            ", values=" + values +
            '}';
    }
}
