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
import org.apache.plc4x.java.spi.drivers.model.AddressConstraints;
import org.apache.plc4x.java.spi.drivers.model.ArrayNotationParser;
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
        "^(?<device>[A-Za-z]+)(?<hexPrefix>0[xX])?(?<address>[0-9A-Fa-f]+)"
            + ArrayNotationParser.ARRAY_GROUP
            + "(:(?<datatype>[A-Za-z_]+))?$");

    /** Conservative single-frame word ceiling for 3E binary Batch Read/Write (not the exact device max). */
    static final int MAX_POINTS = 960;

    /** Device addresses are serialized as an unsigned 24-bit field in the 3E frame. */
    static final int MAX_DEVICE_NUMBER = 0xFFFFFF;

    private final SlmpDeviceCode deviceCode;
    private final int deviceNumber;
    private final SlmpDataType dataType;
    private final int quantity;

    /**
     * Whether the address wrote the selection as a range. A one-element range is still a range -
     * {@code D100[4]} yields a scalar and {@code D100[4..4]} a list of one - and the count alone
     * cannot say which was written, so the parser's answer is carried here.
     */
    private final boolean explicitRange;

    public SlmpTag(SlmpDeviceCode deviceCode, int deviceNumber, SlmpDataType dataType, int quantity) {
        this(deviceCode, deviceNumber, dataType, quantity, quantity > 1);
    }

    public SlmpTag(SlmpDeviceCode deviceCode, int deviceNumber, SlmpDataType dataType, int quantity,
                   boolean explicitRange) {
        this.deviceCode = deviceCode;
        this.deviceNumber = deviceNumber;
        this.dataType = dataType;
        this.quantity = quantity;
        this.explicitRange = explicitRange;
    }

    public static SlmpTag of(String addressString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (!matcher.matches()) {
            throw ArrayNotationParser.invalidAddress(addressString,
                "{device}{address}[selection]:{TYPE} - for example D100[0..3]:INT");
        }
        String deviceToken = matcher.group("device").toUpperCase();
        SlmpDeviceCode device;
        int radix = switch (deviceToken) {
            case "D" -> {
                device = SlmpDeviceCode.D;
                yield 10;
            }
            case "R" -> {
                device = SlmpDeviceCode.R;
                yield 10;
            }
            case "W" -> {
                device = SlmpDeviceCode.W;
                yield 16;
            }
            default -> throw new PlcInvalidTagException(
                "device '" + deviceToken + "' not supported in this version (word devices D/W/R only)");
        };

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

        // The selection sits between the address and the type; its offset moves the device
        // number, and its size is how many devices are read.
        String arrayToken = matcher.group("array");
        int quantity = 1;
        boolean explicitRange = false;
        if (arrayToken != null) {
            ArrayInfo dimension = ArrayNotationParser
                .parse(arrayToken, addressString, AddressConstraints.SINGLE_DIMENSION).getFirst();
            // The offset counts elements; a device number counts 16-bit words. They coincide only
            // for a one-word type, which is why D100[4]:INT looked right while D100[4]:DINT read
            // four words short of its target. This is the same scale getWordsPerElement() applies
            // to the point count below, so the offset and the length cannot disagree.
            deviceNumber += (dimension.getLowerBound() - dimension.getBase()) * dataType.getWordsPerElement();
            quantity = dimension.getSize();
            explicitRange = dimension.isRange();
        }
        if (quantity < 1) {
            throw new PlcInvalidTagException("quantity must be >= 1 in: " + addressString);
        }

        long numberOfPoints = (long) quantity * dataType.getWordsPerElement();
        if (numberOfPoints > MAX_POINTS) {
            throw new PlcInvalidTagException("requested " + numberOfPoints + " words exceeds the v0 single-frame "
                + "Batch Read/Write ceiling of " + MAX_POINTS + " (no optimizer to split): " + addressString);
        }
        return new SlmpTag(device, deviceNumber, dataType, quantity, explicitRange);
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
        sb.append(ArrayNotationParser.render(getArrayInfo()));
        if (dataType != SlmpDataType.WORD || quantity != 1) {
            sb.append(':').append(dataType.name());
        }
        return sb.toString();
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.valueOf(dataType.name());
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        // A range is an array even when it spans one element, so the flag decides the shape and
        // the count only sizes it. Deriving the shape from the count alone reported D100[4..4] as
        // a scalar, contradicting the notation's own rule and plc4go's SLMP tag, which carries
        // the same flag.
        if (explicitRange) {
            return Collections.singletonList(new DefaultArrayInfo(0, quantity - 1, 0, true));
        }
        return Collections.emptyList();
    }

    /** Whether the address wrote a range, as opposed to selecting a single element. */
    public boolean isExplicitRange() {
        return explicitRange;
    }

    @Override
    public String toString() {
        return "SlmpTag{" + getAddressString() + '}';
    }
}
