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
package org.apache.plc4x.java.slmp.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.slmp.SlmpDataType;
import org.apache.plc4x.java.slmp.readwrite.SlmpDeviceCode;
import org.apache.plc4x.java.spi.drivers.model.DefaultArrayInfo;

import java.io.Serializable;
import java.util.Collections;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * A single SLMP word-device tag. v0 supports the word devices D (decimal addr),
 * W (hex addr) and R (decimal addr); bit devices are rejected. Address grammar:
 * {@code <device><address>[:<datatype>][\[<quantity>\]]}, e.g. {@code D350},
 * {@code R200:REAL[4]}, {@code W1A:WORD[10]} (W also accepts {@code W0x1A}).
 */
public class SlmpTag implements PlcTag, Serializable {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile(
        "^(?<device>[A-Za-z]+)(?<hexPrefix>0[xX])?(?<address>[0-9A-Fa-f]+)" +
        "(:(?<datatype>[A-Za-z_]+))?(\\[(?<quantity>\\d+)])?$");

    /** Conservative single-frame word ceiling for 3E binary Batch Read/Write (not the exact device max). */
    static final int MAX_POINTS = 960;

    /** Device addresses are serialized as an unsigned 24-bit field in the 3E frame. */
    static final int MAX_DEVICE_NUMBER = 0xFFFFFF;

    private final SlmpDeviceCode deviceCode;
    private final int deviceNumber;
    private final SlmpDataType dataType;
    private final int quantity;

    public SlmpTag(SlmpDeviceCode deviceCode, int deviceNumber, SlmpDataType dataType, int quantity) {
        this.deviceCode = deviceCode;
        this.deviceNumber = deviceNumber;
        this.dataType = dataType;
        this.quantity = quantity;
    }

    public static SlmpTag of(String addressString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (!matcher.matches()) {
            throw new PlcInvalidTagException("Unable to parse SLMP address: " + addressString);
        }
        String deviceToken = matcher.group("device").toUpperCase();
        SlmpDeviceCode device;
        int radix;
        switch (deviceToken) {
            case "D":
                device = SlmpDeviceCode.D;
                radix = 10;
                break;
            case "R":
                device = SlmpDeviceCode.R;
                radix = 10;
                break;
            case "W":
                device = SlmpDeviceCode.W;
                radix = 16;
                break;
            default:
                throw new PlcInvalidTagException(
                    "device '" + deviceToken + "' not supported in this version (word devices D/W/R only)");
        }

        boolean hasHexPrefix = matcher.group("hexPrefix") != null;
        if (hasHexPrefix && radix != 16) {
            throw new PlcInvalidTagException("0x prefix is only valid for hex devices (W): " + addressString);
        }
        int deviceNumber;
        try {
            deviceNumber = Integer.parseInt(matcher.group("address"), radix);
        } catch (NumberFormatException e) {
            throw new PlcInvalidTagException("Invalid " + (radix == 16 ? "hex" : "decimal")
                + " device number in: " + addressString);
        }
        if (deviceNumber > MAX_DEVICE_NUMBER) {
            throw new PlcInvalidTagException("device number " + deviceNumber
                + " exceeds the 24-bit SLMP device-address range [0.." + MAX_DEVICE_NUMBER + "]: " + addressString);
        }

        String datatypeToken = matcher.group("datatype");
        SlmpDataType dataType;
        if (datatypeToken == null) {
            dataType = SlmpDataType.WORD;
        } else {
            try {
                dataType = SlmpDataType.valueOf(datatypeToken.toUpperCase());
            } catch (IllegalArgumentException e) {
                throw new PlcInvalidTagException("Unsupported SLMP data type '" + datatypeToken
                    + "' (supported: WORD, INT, UINT, DINT, UDINT, REAL)");
            }
        }

        String quantityToken = matcher.group("quantity");
        int quantity;
        if (quantityToken == null) {
            quantity = 1;
        } else {
            try {
                quantity = Integer.parseInt(quantityToken);
            } catch (NumberFormatException e) {
                throw new PlcInvalidTagException("quantity out of range in: " + addressString);
            }
        }
        if (quantity < 1) {
            throw new PlcInvalidTagException("quantity must be >= 1 in: " + addressString);
        }

        long numberOfPoints = (long) quantity * dataType.getWordsPerElement();
        if (numberOfPoints > MAX_POINTS) {
            throw new PlcInvalidTagException("requested " + numberOfPoints + " words exceeds the v0 single-frame "
                + "Batch Read/Write ceiling of " + MAX_POINTS + " (no optimizer to split): " + addressString);
        }
        return new SlmpTag(device, deviceNumber, dataType, quantity);
    }

    public SlmpDeviceCode getDeviceCode() {
        return deviceCode;
    }

    public int getDeviceNumber() {
        return deviceNumber;
    }

    public SlmpDataType getDataType() {
        return dataType;
    }

    public int getQuantity() {
        return quantity;
    }

    /** Number of 16-bit words to transfer (Batch Read/Write, word units). */
    public int getNumberOfPoints() {
        return quantity * dataType.getWordsPerElement();
    }

    @Override
    public String getAddressString() {
        String addr = (deviceCode == SlmpDeviceCode.W)
            ? "0x" + Integer.toHexString(deviceNumber).toUpperCase()
            : Integer.toString(deviceNumber);
        StringBuilder sb = new StringBuilder(deviceCode.name()).append(addr);
        if (dataType != SlmpDataType.WORD || quantity != 1) {
            sb.append(':').append(dataType.name());
        }
        if (quantity != 1) {
            sb.append('[').append(quantity).append(']');
        }
        return sb.toString();
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.valueOf(dataType.name());
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        if (quantity > 1) {
            return Collections.singletonList(new DefaultArrayInfo(0, quantity - 1));
        }
        return Collections.emptyList();
    }

    @Override
    public String toString() {
        return "SlmpTag{" + getAddressString() + '}';
    }
}
