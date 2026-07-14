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
package org.apache.plc4x.java.slmp;

import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.slmp.tag.SlmpTag;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/** Maps a single 3E response (endCode + raw words) to a per-tag {@link PlcResponseItem}. */
final class SlmpResponseMapper {

    private static final Logger LOGGER = LoggerFactory.getLogger(SlmpResponseMapper.class);

    private SlmpResponseMapper() {
    }

    static PlcResponseItem<PlcValue> mapTag(SlmpTag tag, int endCode, byte[] responseData) {
        if (endCode != 0x0000) {
            LOGGER.warn("SLMP device returned endCode {} for {}", String.format("0x%04X", endCode), tag);
            return new DefaultPlcResponseItem<>(PlcResponseCode.REMOTE_ERROR, null);
        }
        PlcValue value = tag.getDataType().decode(responseData, tag.getQuantity());
        if (value == null) {
            LOGGER.warn("SLMP response too short for {} ({} bytes)", tag,
                responseData == null ? 0 : responseData.length);
            return new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_DATA, null);
        }
        return new DefaultPlcResponseItem<>(PlcResponseCode.OK, value);
    }
}
