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
package org.apache.plc4x.java.modbus.model;

import org.apache.plc4x.java.api.exceptions.PlcInvalidAddressException;

import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class MaskWriteRegisterModbusAddress extends ModbusAddress {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("maskwrite:" + ModbusAddress.ADDRESS_PATTERN + "/" + "(?<andMask>\\d+)/(?<orMask>\\d+)");

    private final int andMask;
    private final int orMask;

    protected MaskWriteRegisterModbusAddress(int address, int andMask, int orMask) {
        super(address);
        this.andMask = andMask;
        this.orMask = orMask;
    }

    public static MaskWriteRegisterModbusAddress of(String addressString) throws PlcInvalidAddressException {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (!matcher.matches()) {
            throw new PlcInvalidAddressException(addressString, ADDRESS_PATTERN);
        }
        int address = Integer.parseInt(matcher.group("address"));
        int andMask = Integer.parseInt(matcher.group("andMask"));
        int orMask = Integer.parseInt(matcher.group("orMask"));
        return new MaskWriteRegisterModbusAddress(address, andMask, orMask);
    }

    public int getAndMask() {
        return andMask;
    }

    public int getOrMask() {
        return orMask;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof MaskWriteRegisterModbusAddress)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }
        MaskWriteRegisterModbusAddress that = (MaskWriteRegisterModbusAddress) o;
        return andMask == that.andMask &&
            orMask == that.orMask;
    }

    @Override
    public int hashCode() {

        return Objects.hash(super.hashCode(), andMask, orMask);
    }

    @Override
    public String toString() {
        return "MaskWriteRegisterModbusAddress{" +
            "andMask=" + andMask +
            ", orMask=" + orMask +
            "} " + super.toString();
    }
}
