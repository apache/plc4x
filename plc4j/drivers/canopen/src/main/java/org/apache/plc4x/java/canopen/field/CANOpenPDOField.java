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
package org.apache.plc4x.java.canopen.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.canopen.readwrite.CANOpenDataType;
import org.apache.plc4x.java.canopen.readwrite.CANOpenService;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.nio.charset.StandardCharsets;
import java.util.Collections;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class CANOpenPDOField extends CANOpenField implements CANOpenSubscriptionField {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?<pdo>(?:RECEIVE|TRANSMIT)_PDO_[1-4]):" + NODE_PATTERN + ":(?<canDataType>\\w+)(\\[(?<numberOfElements>\\d)])?");
    private final CANOpenService service;
    private final CANOpenDataType canOpenDataType;

    public CANOpenPDOField(int node, CANOpenService service, CANOpenDataType canOpenDataType) {
        super(node);
        this.service = service;
        this.canOpenDataType = canOpenDataType;
    }

    public CANOpenDataType getCanOpenDataType() {
        return canOpenDataType;
    }

    @Override
    public String getAddressString() {
        // Number of Elements not implemented.
        return service.name() + ":" + canOpenDataType.name();
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.valueOf(canOpenDataType.getPlcValueName());
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        // Number of elements not implemented.
        return Collections.emptyList();
    }

    public static boolean matches(String addressString) {
        return ADDRESS_PATTERN.matcher(addressString).matches();
    }

    public static Matcher getMatcher(String addressString) throws PlcInvalidFieldException {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (matcher.matches()) {
            return matcher;
        }

        throw new PlcInvalidFieldException(addressString, ADDRESS_PATTERN);
    }

    public static CANOpenPDOField of(String addressString) {
        Matcher matcher = getMatcher(addressString);
        int nodeId = Integer.parseInt(matcher.group("nodeId"));

        String pdo = matcher.group("pdo");
        CANOpenService service = CANOpenService.valueOf(pdo);

        String canDataTypeString = matcher.group("canDataType");
        CANOpenDataType canOpenDataType = CANOpenDataType.valueOf(canDataTypeString);

        return new CANOpenPDOField(nodeId, service, canOpenDataType);
    }

    public CANOpenService getService() {
        return service;
    }

    @Override
    public boolean isWildcard() {
        return false;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());

        String serviceName = getService().name();
        writeBuffer.writeString("service", serviceName.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), serviceName);
        writeBuffer.writeInt("node",64, getNodeId());
        String dataTypeName = getCanOpenDataType().name();
        writeBuffer.writeString("dataType", dataTypeName.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), dataTypeName);

        writeBuffer.popContext(getClass().getSimpleName());
    }
}
