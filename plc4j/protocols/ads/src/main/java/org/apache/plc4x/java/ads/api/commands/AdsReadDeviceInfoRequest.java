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
package org.apache.plc4x.java.ads.api.commands;

import org.apache.plc4x.java.ads.api.generic.AdsData;
import org.apache.plc4x.java.ads.api.generic.AmsHeader;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;

/**
 * Reads the name and the version number of the ADS device.
 * <p>
 * No additional data required
 */
@AdsCommandType(Command.ADS_READ_DEVICE_INFO)
public class AdsReadDeviceInfoRequest extends AdsAbstractRequest {

    private AdsReadDeviceInfoRequest(AmsHeader amsHeader) {
        super(amsHeader);
    }

    private AdsReadDeviceInfoRequest(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
    }

    public static AdsReadDeviceInfoRequest of(AmsHeader amsHeader) {
        return new AdsReadDeviceInfoRequest(amsHeader);
    }

    public static AdsReadDeviceInfoRequest of(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId) {
        return new AdsReadDeviceInfoRequest(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
    }

    @Override
    public AdsData getAdsData() {
        return AdsData.EMPTY;
    }

}
