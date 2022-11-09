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
package org.apache.plc4x.java.canopen.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.canopen.readwrite.CANOpenService;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.nio.charset.StandardCharsets;
import java.util.Collections;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class CANOpenNMTTag extends CANOpenTag implements CANOpenSubscriptionTag {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("NMT|NMT:" + NODE_PATTERN);

    public CANOpenNMTTag(int node) {
        super(node);
    }

    @Override
    public CANOpenService getService() {
        return CANOpenService.NMT;
    }

    public boolean isWildcard() {
        return getNodeId() == 0;
    }

    @Override
    public String getAddressString() {
        return "NMT:" + getNodeId();
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.NULL;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return Collections.emptyList();
    }

    public static boolean matches(String addressString) {
        return ADDRESS_PATTERN.matcher(addressString).matches();
    }

    public static Matcher getMatcher(String addressString) throws PlcInvalidTagException {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (matcher.matches()) {
            return matcher;
        }

        throw new PlcInvalidTagException(addressString, ADDRESS_PATTERN);
    }

    public static CANOpenNMTTag of(String addressString) {
        Matcher matcher = getMatcher(addressString);
        int nodeId = matcher.group("nodeId") == null ? 0 : Integer.parseInt(matcher.group("nodeId"));

        return new CANOpenNMTTag(nodeId);
    }


    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());

        String serviceName = getService().name();
        writeBuffer.writeString("service", serviceName.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), serviceName);
        writeBuffer.writeInt("node",64, getNodeId());

        writeBuffer.popContext(getClass().getSimpleName());
    }
}
