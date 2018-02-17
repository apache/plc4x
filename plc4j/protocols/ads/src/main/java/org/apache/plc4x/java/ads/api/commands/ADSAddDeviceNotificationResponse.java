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
import org.apache.plc4x.java.ads.api.commands.types.Result;
import org.apache.plc4x.java.ads.api.generic.ADSData;
import org.apache.plc4x.java.ads.api.generic.AMSHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;

import static java.util.Objects.requireNonNull;

/**
 * A notification is created in an ADS device.
 */
@ADSCommandType(Command.ADS_ADD_DEVICE_NOTIFICATION)
public class ADSAddDeviceNotificationResponse extends ADSAbstractResponse {

    /**
     * 4 bytes	ADS error number
     */
    private final Result result;

    /**
     * 4 bytes	Handle of notification
     */
    private final NotificationHandle notificationHandle;

    private ADSAddDeviceNotificationResponse(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, Result result, NotificationHandle notificationHandle) {
        super(amstcpHeader, amsHeader);
        this.result = requireNonNull(result);
        this.notificationHandle = requireNonNull(notificationHandle);
    }

    private ADSAddDeviceNotificationResponse(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, Result result, NotificationHandle notificationHandle) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.result = requireNonNull(result);
        this.notificationHandle = requireNonNull(notificationHandle);
    }

    public static ADSAddDeviceNotificationResponse of(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, Result result, NotificationHandle notificationHandle) {
        return new ADSAddDeviceNotificationResponse(amstcpHeader, amsHeader, result, notificationHandle);
    }

    public static ADSAddDeviceNotificationResponse of(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, Result result, NotificationHandle notificationHandle) {
        return new ADSAddDeviceNotificationResponse(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, result, notificationHandle);
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(result, notificationHandle);
    }

    @Override
    public String toString() {
        return "ADSAddDeviceNotificationResponse{" +
            "result=" + result +
            ", notificationHandle=" + notificationHandle +
            "} " + super.toString();
    }
}
