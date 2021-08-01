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
package org.apache.plc4x.java.ads.api.commands;

import org.apache.plc4x.java.ads.api.commands.types.NotificationHandle;
import org.apache.plc4x.java.ads.api.commands.types.Result;
import org.apache.plc4x.java.ads.api.generic.AdsData;
import org.apache.plc4x.java.ads.api.generic.AmsHeader;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;

import static java.util.Objects.requireNonNull;

/**
 * A notification is created in an ADS device.
 */
@AdsCommandType(Command.ADS_ADD_DEVICE_NOTIFICATION)
public class AdsAddDeviceNotificationResponse extends AdsAbstractResponse {

    /**
     * 4 bytes	ADS error number
     */
    private final Result result;

    /**
     * 4 bytes	Handle of notification
     */
    private final NotificationHandle notificationHandle;

    private AdsAddDeviceNotificationResponse(AmsHeader amsHeader, Result result, NotificationHandle notificationHandle) {
        super(amsHeader);
        this.result = requireNonNull(result);
        this.notificationHandle = requireNonNull(notificationHandle);
    }

    private AdsAddDeviceNotificationResponse(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, Result result, NotificationHandle notificationHandle) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.result = requireNonNull(result);
        this.notificationHandle = requireNonNull(notificationHandle);
    }

    public static AdsAddDeviceNotificationResponse of(AmsHeader amsHeader, Result result, NotificationHandle notificationHandle) {
        return new AdsAddDeviceNotificationResponse(amsHeader, result, notificationHandle);
    }

    public static AdsAddDeviceNotificationResponse of(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, Result result, NotificationHandle notificationHandle) {
        return new AdsAddDeviceNotificationResponse(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, result, notificationHandle);
    }

    public Result getResult() {
        return result;
    }

    public NotificationHandle getNotificationHandle() {
        return notificationHandle;
    }

    @Override
    public AdsData getAdsData() {
        return buildADSData(result, notificationHandle);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof AdsAddDeviceNotificationResponse))
            return false;
        if (!super.equals(o))
            return false;

        AdsAddDeviceNotificationResponse that = (AdsAddDeviceNotificationResponse) o;

        if (!result.equals(that.result))
            return false;

        return notificationHandle.equals(that.notificationHandle);
    }

    @Override
    public int hashCode() {
        int result1 = super.hashCode();
        result1 = 31 * result1 + result.hashCode();
        result1 = 31 * result1 + notificationHandle.hashCode();
        return result1;
    }

    @Override
    public String toString() {
        return "AdsAddDeviceNotificationResponse{" +
            "result=" + result +
            ", notificationHandle=" + notificationHandle +
            "} " + super.toString();
    }
}
