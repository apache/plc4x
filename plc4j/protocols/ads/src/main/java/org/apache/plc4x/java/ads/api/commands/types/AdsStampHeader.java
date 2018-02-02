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
package org.apache.plc4x.java.ads.api.commands.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.generic.AMSTCPPaket;
import org.apache.plc4x.java.ads.api.util.ByteReadable;

public class AdsStampHeader implements ByteReadable {

    /**
     * 8 bytes	The timestamp is coded after the Windos FILETIME format. I.e. the value contains the number of the nano seconds, which passed since 1.1.1601. In addition, the local time change is not considered. Thus the time stamp is present as universal Coordinated time (UTC).
     */
    private final TimeStamp timeStamp;
    /**
     * 4 bytes	Number of elements of type AdsNotificationSample.
     */
    private final Samples samples;
    /**
     * n bytes	Array with elements of type AdsNotificationSample.
     */
    private final AdsNotificationSample adsNotificationSample;

    AdsStampHeader(TimeStamp timeStamp, Samples samples, AdsNotificationSample adsNotificationSample) {
        this.timeStamp = timeStamp;
        this.samples = samples;
        this.adsNotificationSample = adsNotificationSample;
    }

    public static AdsStampHeader of(TimeStamp timeStamp, Samples samples, AdsNotificationSample adsNotificationSample) {
        return new AdsStampHeader(timeStamp, samples, adsNotificationSample);
    }

    @Override
    public ByteBuf getByteBuf() {
        return AMSTCPPaket.buildByteBuff(timeStamp, samples, adsNotificationSample);
    }
}
