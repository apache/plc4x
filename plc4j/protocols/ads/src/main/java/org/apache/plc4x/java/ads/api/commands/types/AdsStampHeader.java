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

import java.util.List;

import static java.util.Objects.requireNonNull;

public class AdsStampHeader implements ByteReadable {

    /**
     * 8 bytes	The timestamp is coded after the Windows FILETIME format. I.e. the value contains the number of the nano seconds, which passed since 1.1.1601. In addition, the local time change is not considered. Thus the time stamp is present as universal Coordinated time (UTC).
     */
    private final TimeStamp timeStamp;
    /**
     * 4 bytes	Number of elements of type AdsNotificationSample.
     */
    private final Samples samples;
    /**
     * n bytes	Array with elements of type AdsNotificationSample.
     */
    private final List<AdsNotificationSample> adsNotificationSamples;

    private AdsStampHeader(TimeStamp timeStamp, Samples samples, List<AdsNotificationSample> adsNotificationSamples) {
        this.timeStamp = requireNonNull(timeStamp);
        this.samples = requireNonNull(samples);
        this.adsNotificationSamples = requireNonNull(adsNotificationSamples);
    }

    private AdsStampHeader(TimeStamp timeStamp, List<AdsNotificationSample> adsNotificationSamples) {
        this.timeStamp = requireNonNull(timeStamp);
        this.adsNotificationSamples = requireNonNull(adsNotificationSamples);
        this.samples = Samples.of(adsNotificationSamples.size());
    }

    public static AdsStampHeader of(TimeStamp timeStamp, Samples samples, List<AdsNotificationSample> adsNotificationSamples) {
        return new AdsStampHeader(timeStamp, samples, adsNotificationSamples);
    }

    public static AdsStampHeader of(TimeStamp timeStamp, List<AdsNotificationSample> adsNotificationSamples) {
        return new AdsStampHeader(timeStamp, adsNotificationSamples);
    }

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(timeStamp, samples, () -> buildByteBuff(adsNotificationSamples.toArray(new ByteReadable[adsNotificationSamples.size()])));
    }

    public TimeStamp getTimeStamp() {
        return timeStamp;
    }

    public Samples getSamples() {
        return samples;
    }

    public List<AdsNotificationSample> getAdsNotificationSamples() {
        return adsNotificationSamples;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof AdsStampHeader))
            return false;

        AdsStampHeader that = (AdsStampHeader) o;

        if (!timeStamp.equals(that.timeStamp))
            return false;
        if (!samples.equals(that.samples))
            return false;
        
        return adsNotificationSamples.equals(that.adsNotificationSamples);
    }

    @Override
    public int hashCode() {
        int result = timeStamp.hashCode();
        result = 31 * result + samples.hashCode();
        result = 31 * result + adsNotificationSamples.hashCode();
        return result;
    }

    @Override
    public String toString() {
        return "AdsStampHeader{" +
            "timeStamp=" + timeStamp +
            ", samples=" + samples +
            ", adsNotificationSamples=" + adsNotificationSamples +
            '}';
    }
}
