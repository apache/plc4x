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
import org.apache.plc4x.java.ads.api.generic.AdsData;
import org.apache.plc4x.java.ads.api.generic.AmsHeader;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;

import static java.util.Objects.requireNonNull;

/**
 * One before defined notification is deleted in an ADS device.
 */
@AdsCommandType(Command.ADS_DELETE_DEVICE_NOTIFICATION)
public class AdsDeleteDeviceNotificationRequest extends AdsAbstractRequest {

    /**
     * 4 bytes	Handle of notification
     */
    private final NotificationHandle notificationHandle;

    private AdsDeleteDeviceNotificationRequest(AmsHeader amsHeader, NotificationHandle notificationHandle) {
        super(amsHeader);
        this.notificationHandle = requireNonNull(notificationHandle);
    }

    private AdsDeleteDeviceNotificationRequest(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, NotificationHandle notificationHandle) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.notificationHandle = requireNonNull(notificationHandle);
    }

    public static AdsDeleteDeviceNotificationRequest of(AmsHeader amsHeader, NotificationHandle notificationHandle) {
        return new AdsDeleteDeviceNotificationRequest(amsHeader, notificationHandle);
    }

    public static AdsDeleteDeviceNotificationRequest of(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, NotificationHandle notificationHandle) {
        return new AdsDeleteDeviceNotificationRequest(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, notificationHandle);
    }

    public NotificationHandle getNotificationHandle() {
        return notificationHandle;
    }

    @Override
    public AdsData getAdsData() {
        return buildADSData(notificationHandle);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof AdsDeleteDeviceNotificationRequest))
            return false;
        if (!super.equals(o))
            return false;

        AdsDeleteDeviceNotificationRequest that = (AdsDeleteDeviceNotificationRequest) o;

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
        return "AdsDeleteDeviceNotificationRequest{" +
            "notificationHandle=" + notificationHandle +
            "} " + super.toString();
    }
}
