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
package org.apache.plc4x.java.ads.model.commands;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.model.generic.ADSData;
import org.apache.plc4x.java.ads.model.generic.AMSHeader;
import org.apache.plc4x.java.ads.model.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.model.generic.AMSTCPPaket;
import org.apache.plc4x.java.ads.model.util.ByteReadable;
import org.apache.plc4x.java.ads.model.util.ByteValue;

/**
 * Data will carry forward independently from an ADS device to a Client
 * <p>
 * The data which are transfered at the Device Notification  are multiple nested into one another. The Notification Stream contains an array with elements of type AdsStampHeader. This array again contains elements of type AdsNotificationSample.
 */
public class ADSDeviceNotificationRequest extends AMSTCPPaket {

    /**
     * 4 bytes	Size of data in byte.
     */
    final Length length;
    /**
     * 4 bytes	Number of elements of type AdsStampHeader.
     */
    final Stamps stamps;
    /**
     * n bytes	Array with elements of type AdsStampHeader.
     */
    final AdsStampHeader adsStampHeader;

    public ADSDeviceNotificationRequest(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, Length length, Stamps stamps, AdsStampHeader adsStampHeader) {
        super(amstcpHeader, amsHeader);
        this.length = length;
        this.stamps = stamps;
        this.adsStampHeader = adsStampHeader;
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(length, stamps, adsStampHeader);
    }

    class Length extends ByteValue {
        public Length(byte... value) {
            super(value);
            assertLength(4);
        }
    }

    class Stamps extends ByteValue {
        public Stamps(byte... value) {
            super(value);
            assertLength(4);
        }
    }

    class AdsStampHeader implements ByteReadable {

        /**
         * 8 bytes	The timestamp is coded after the Windos FILETIME format. I.e. the value contains the number of the nano seconds, which passed since 1.1.1601. In addition, the local time change is not considered. Thus the time stamp is present as universal Coordinated time (UTC).
         */
        final TimeStamp timeStamp;
        /**
         * 4 bytes	Number of elements of type AdsNotificationSample.
         */
        final Samples samples;
        /**
         * n bytes	Array with elements of type AdsNotificationSample.
         */
        final AdsNotificationSample adsNotificationSample;

        public AdsStampHeader(TimeStamp timeStamp, Samples samples, AdsNotificationSample adsNotificationSample) {
            this.timeStamp = timeStamp;
            this.samples = samples;
            this.adsNotificationSample = adsNotificationSample;
        }

        class TimeStamp extends ByteValue {

            public TimeStamp(byte... value) {
                super(value);
                assertLength(8);
            }
        }

        class Samples extends ByteValue {

            public Samples(byte... value) {
                super(value);
                assertLength(4);
            }
        }

        class AdsNotificationSample implements ByteReadable {

            /**
             * 4 Bytes	Handle of notification.
             */
            final NotificationHandle notificationHandle;
            /**
             * 4 Bytes	Size of data range in bytes.
             */
            final SampleSize sampleSize;
            /**
             * n Bytes	Data
             */
            final Data data;

            public AdsNotificationSample(NotificationHandle notificationHandle, SampleSize sampleSize, Data data) {
                this.notificationHandle = notificationHandle;
                this.sampleSize = sampleSize;
                this.data = data;
            }

            @Override
            public byte[] getBytes() {
                return getByteBuf().array();
            }

            @Override
            public ByteBuf getByteBuf() {
                return buildByteBuff(notificationHandle, sampleSize, data);
            }

            class NotificationHandle extends ByteValue {

                public NotificationHandle(byte... value) {
                    super(value);
                    assertLength(4);
                }
            }

            /**
             * Notice: If your handle becomes invalid, one notification without data will be send once as advice.
             */
            class InvalidationNotificationHandle extends NotificationHandle {
                public InvalidationNotificationHandle() {
                }
            }

            class SampleSize extends ByteValue {

                public SampleSize(byte... value) {
                    super(value);
                    assertLength(4);
                }
            }

            class Data extends ByteValue {

                public Data(byte... value) {
                    super(value);
                }
            }
        }

        @Override
        public byte[] getBytes() {
            return getByteBuf().array();
        }

        @Override
        public ByteBuf getByteBuf() {
            return buildByteBuff(timeStamp, samples, adsNotificationSample);
        }
    }
}
