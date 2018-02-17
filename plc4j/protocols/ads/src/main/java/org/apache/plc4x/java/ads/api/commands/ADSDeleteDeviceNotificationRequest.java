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
package org.apache.plc4x.java.ads.api.commands;

import org.apache.plc4x.java.ads.api.commands.types.NotificationHandle;
import org.apache.plc4x.java.ads.api.generic.ADSData;
import org.apache.plc4x.java.ads.api.generic.AMSHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;

import static java.util.Objects.requireNonNull;

/**
 * One before defined notification is deleted in an ADS device.
 */
@ADSCommandType(Command.ADS_DELETE_DEVICE_NOTIFICATION)
public class ADSDeleteDeviceNotificationRequest extends ADSAbstractRequest {

    /**
     * 4 bytes	Handle of notification
     */
    private final NotificationHandle notificationHandle;

    private ADSDeleteDeviceNotificationRequest(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, NotificationHandle notificationHandle) {
        super(amstcpHeader, amsHeader);
        this.notificationHandle = requireNonNull(notificationHandle);
    }

    private ADSDeleteDeviceNotificationRequest(AMSHeader amsHeader, NotificationHandle notificationHandle) {
        super(amsHeader);
        this.notificationHandle = requireNonNull(notificationHandle);
    }

    private ADSDeleteDeviceNotificationRequest(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, NotificationHandle notificationHandle) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.notificationHandle = requireNonNull(notificationHandle);
    }

    public static ADSDeleteDeviceNotificationRequest of(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, NotificationHandle notificationHandle) {
        return new ADSDeleteDeviceNotificationRequest(amstcpHeader, amsHeader, notificationHandle);
    }

    public static ADSDeleteDeviceNotificationRequest of(AMSHeader amsHeader, NotificationHandle notificationHandle) {
        return new ADSDeleteDeviceNotificationRequest(amsHeader, notificationHandle);
    }

    public static ADSDeleteDeviceNotificationRequest of(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, NotificationHandle notificationHandle) {
        return new ADSDeleteDeviceNotificationRequest(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, notificationHandle);
    }

    public NotificationHandle getNotificationHandle() {
        return notificationHandle;
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(notificationHandle);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof ADSDeleteDeviceNotificationRequest)) return false;
        if (!super.equals(o)) return false;

        ADSDeleteDeviceNotificationRequest that = (ADSDeleteDeviceNotificationRequest) o;

        return notificationHandle.equals(that.notificationHandle);
    }

    @Override
    public int hashCode() {
        int result = super.hashCode();
        result = 31 * result + notificationHandle.hashCode();
        return result;
    }

    @Override
    public String toString() {
        return "ADSDeleteDeviceNotificationRequest{" +
            "notificationHandle=" + notificationHandle +
            "} " + super.toString();
    }
}
