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
package org.apache.plc4x.java.ads.api.commands.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.ByteReadable;
import org.apache.plc4x.java.ads.api.util.LengthSupplier;

import static java.util.Objects.requireNonNull;

public class AdsNotificationSample implements ByteReadable {

    /**
     * 4 Bytes	Handle of notification.
     */
    private final NotificationHandle notificationHandle;
    /**
     * 4 Bytes	Size of data range in bytes.
     */
    private final SampleSize sampleSize;
    /**
     * n Bytes	Data
     */
    private final Data data;

    private final transient LengthSupplier lengthSupplier;

    private AdsNotificationSample(NotificationHandle notificationHandle, Data data) {
        this.notificationHandle = requireNonNull(notificationHandle);
        this.sampleSize = null;
        this.data = requireNonNull(data);
        lengthSupplier = data;
    }

    private AdsNotificationSample(NotificationHandle notificationHandle, SampleSize sampleSize, Data data) {
        this.notificationHandle = requireNonNull(notificationHandle);
        this.sampleSize = requireNonNull(sampleSize);
        this.data = requireNonNull(data);
        lengthSupplier = null;
    }

    public static AdsNotificationSample of(NotificationHandle notificationHandle, Data data) {
        return new AdsNotificationSample(notificationHandle, data);
    }

    public static AdsNotificationSample of(NotificationHandle notificationHandle, SampleSize sampleSize, Data data) {
        return new AdsNotificationSample(notificationHandle, sampleSize, data);
    }

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(notificationHandle, getSampleSize(), data);
    }

    public NotificationHandle getNotificationHandle() {
        return notificationHandle;
    }

    public SampleSize getSampleSize() {
        return lengthSupplier == null ? sampleSize : SampleSize.of(lengthSupplier);
    }

    public Data getData() {
        return data;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof AdsNotificationSample))
            return false;

        AdsNotificationSample that = (AdsNotificationSample) o;

        if (!notificationHandle.equals(that.notificationHandle))
            return false;
        if (!data.equals(that.data))
            return false;

        return getSampleSize().equals(that.getSampleSize());
    }

    @Override
    public int hashCode() {
        int result = notificationHandle.hashCode();
        result = 31 * result + data.hashCode();
        result = 31 * result + getSampleSize().hashCode();
        return result;
    }

    @Override
    public String toString() {
        return "AdsNotificationSample{" +
            "notificationHandle=" + notificationHandle +
            ", sampleSize=" + getSampleSize() +
            ", data=" + data +
            '}';
    }
}
