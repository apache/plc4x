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
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;

import java.util.Objects;

public class AdsSubscriptionHandle extends DefaultPlcSubscriptionHandle {

    private final String plcFieldName;

    private final AdsDataTypeTableEntry adsDataType;

    private final Long notificationHandle;

    public AdsSubscriptionHandle(PlcSubscriber plcSubscriber, String plcFieldName, AdsDataTypeTableEntry adsDataType, Long notificationHandle) {
        super(plcSubscriber);
        this.plcFieldName = plcFieldName;
        this.adsDataType = adsDataType;
        this.notificationHandle = notificationHandle;
    }

    public String getPlcFieldName() {
        return plcFieldName;
    }

    public AdsDataTypeTableEntry getAdsDataType() {
        return adsDataType;
    }

    public Long getNotificationHandle() {
        return notificationHandle;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof AdsSubscriptionHandle)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }
        AdsSubscriptionHandle that = (AdsSubscriptionHandle) o;
        return Objects.equals(plcFieldName, that.plcFieldName) &&
            adsDataType == that.adsDataType &&
            Objects.equals(notificationHandle, that.notificationHandle);
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), plcFieldName, adsDataType, notificationHandle);
    }

    @Override
    public String toString() {
        return "AdsSubscriptionHandle{" +
            "plcFieldName='" + plcFieldName + '\'' +
            ", adsDataType=" + adsDataType +
            ", notificationHandle=" + notificationHandle +
            "} " + super.toString();
    }

}
