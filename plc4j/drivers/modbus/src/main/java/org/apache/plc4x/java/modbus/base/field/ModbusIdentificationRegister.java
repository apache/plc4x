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
package org.apache.plc4x.java.modbus.base.field;

import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.modbus.readwrite.ModbusDeviceInformationLevel;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

public class ModbusIdentificationRegister extends ModbusField {

    public static final String HEX = "0[xX][0-9a-fA-F]+";
    public static final String LEVEL = "(?<level>\\d+|" + HEX + "|(?:BASIC|REGULAR|EXTENDED|INDIVIDUAL))";
    public static final String OBJECT_ID = "(?<objectId>\\d+|" + HEX + ")";
    public static final Pattern ADDRESS_PATTERN = Pattern.compile("identification:" + LEVEL + ":" + OBJECT_ID);

    private final ModbusDeviceInformationLevel level;
    private final short objectId;

    protected ModbusIdentificationRegister(ModbusDeviceInformationLevel level, short objectId) {
        this.level = level;
        this.objectId = objectId;
    }

    public static boolean matches(String addressString) {
        return ADDRESS_PATTERN.matcher(addressString).matches();
    }

    public static Matcher getMatcher(String addressString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (matcher.matches()) {
          return matcher;
        }

        throw new PlcInvalidFieldException(addressString, ADDRESS_PATTERN);
    }

    public static ModbusIdentificationRegister of(String addressString) {
        Matcher matcher = getMatcher(addressString);
        String levelGroup = matcher.group("level");
        String objectidGroup = matcher.group("objectId");

        ModbusDeviceInformationLevel level = parseNumber(levelGroup)
            .map(Integer::byteValue)
            .map(ModbusDeviceInformationLevel::enumForValue)
            .orElseGet(() -> ModbusDeviceInformationLevel.valueOf(levelGroup.toUpperCase()));

        int objectId = parseNumber(objectidGroup)
            .orElseThrow(() -> new IllegalArgumentException("Invalid field definition detected, unknown object id"));

        return new ModbusIdentificationRegister(level, (short) objectId);
    }

    private static Optional<Integer> parseNumber(String value) {
        if (value.matches("\\d+")) {
            return Optional.of(Integer.parseInt(value));
        } else if (value.matches(HEX)) {
            return Optional.of(Integer.parseInt(value.substring(2), 16));
        }
        return Optional.empty();
    }

    public ModbusDeviceInformationLevel getLevel() {
        return level;
    }

    public short getObjectId() {
        return objectId;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());
        writeBuffer.writeShort("level", 8, level.getValue());
        writeBuffer.writeShort("objectId", 8, objectId);
        writeBuffer.popContext(getClass().getSimpleName());
    }
}
