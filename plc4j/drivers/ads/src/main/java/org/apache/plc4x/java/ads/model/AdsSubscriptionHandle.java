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
package org.apache.plc4x.java.ads.model;

import org.apache.plc4x.java.ads.readwrite.AdsDataTypeTableEntry;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;

import java.util.Collections;
import java.util.Objects;
import java.util.function.Consumer;

public class AdsSubscriptionHandle implements PlcSubscriptionHandle {

    private final PlcSubscriber plcSubscriber;
    private final String tagName;
    private final AdsDataTypeTableEntry adsDataType;
    private final Long notificationHandle;

    public AdsSubscriptionHandle(PlcSubscriber plcSubscriber, String tagName, AdsDataTypeTableEntry adsDataType, Long notificationHandle) {
        this.plcSubscriber = plcSubscriber;
        this.tagName = tagName;
        this.adsDataType = adsDataType;
        this.notificationHandle = notificationHandle;
    }

    public String getTagName() {
        return tagName;
    }

    public AdsDataTypeTableEntry getAdsDataType() {
        return adsDataType;
    }

    public Long getNotificationHandle() {
        return notificationHandle;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        return plcSubscriber.registerConsumer(consumer, Collections.singletonList(this));
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof AdsSubscriptionHandle that)) {
            return false;
        }
        return Objects.equals(tagName, that.tagName) &&
            adsDataType == that.adsDataType &&
            Objects.equals(notificationHandle, that.notificationHandle);
    }

    @Override
    public int hashCode() {
        return Objects.hash(tagName, adsDataType, notificationHandle);
    }

    @Override
    public String toString() {
        return "AdsSubscriptionHandle{" +
            "tagName='" + tagName + '\'' +
            ", adsDataType=" + adsDataType +
            ", notificationHandle=" + notificationHandle +
            '}';
    }

}
