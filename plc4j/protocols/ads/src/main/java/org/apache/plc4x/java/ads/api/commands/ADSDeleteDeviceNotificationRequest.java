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

/**
 * One before defined notification is deleted in an ADS device.
 */
@ADSCommandType(Command.ADS_Delete_Device_Notification)
public class ADSDeleteDeviceNotificationRequest extends ADSAbstractRequest {

    /**
     * 4 bytes	Handle of notification
     */
    private final NotificationHandle notificationHandle;

    public ADSDeleteDeviceNotificationRequest(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, NotificationHandle notificationHandle) {
        super(amstcpHeader, amsHeader);
        this.notificationHandle = notificationHandle;
    }

    public ADSDeleteDeviceNotificationRequest(AMSHeader amsHeader, NotificationHandle notificationHandle) {
        super(amsHeader);
        this.notificationHandle = notificationHandle;
    }

    public ADSDeleteDeviceNotificationRequest(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, NotificationHandle notificationHandle) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.notificationHandle = notificationHandle;
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(notificationHandle);
    }

    @Override
    public String toString() {
        return "ADSDeleteDeviceNotificationRequest{" +
            "notificationHandle=" + notificationHandle +
            "} " + super.toString();
    }
}
