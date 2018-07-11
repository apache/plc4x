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
package org.apache.plc4x.java.api.messages.items;

import org.apache.plc4x.java.api.model.SubscriptionHandle;

import java.util.Objects;

public class UnsubscriptionRequestItem {

    private SubscriptionHandle subscriptionHandle;

    public UnsubscriptionRequestItem(SubscriptionHandle subscriptionHandle) {
        Objects.requireNonNull(subscriptionHandle, "SubscriptionHandle must not be null");
        this.subscriptionHandle = subscriptionHandle;
    }

    public SubscriptionHandle getSubscriptionHandle() {
        return subscriptionHandle;
    }

    @Override
    public String toString() {
        return "UnsubscriptionRequestItem{" +
            "subscriptionHandle=" + subscriptionHandle +
            '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof UnsubscriptionRequestItem)) {
            return false;
        }
        UnsubscriptionRequestItem that = (UnsubscriptionRequestItem) o;
        return Objects.equals(subscriptionHandle, that.subscriptionHandle);
    }

    @Override
    public int hashCode() {
        return Objects.hash(subscriptionHandle);
    }

}
