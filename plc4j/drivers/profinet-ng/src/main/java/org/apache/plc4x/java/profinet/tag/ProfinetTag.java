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
package org.apache.plc4x.java.profinet.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ProfinetTag implements PlcTag {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?<address>[\\w\\-. ]+)(:(?<datatype>[a-zA-Z_]+)){1}(\\[(?<quantity>\\d+)])?");
    private final String address;
    private final int quantity;
    private final PlcValueType dataType;

    protected ProfinetTag(String address, Integer quantity, PlcValueType dataType) {
        this.address = address;
        this.quantity = (quantity != null) ? quantity : 1;
        if (this.quantity <= 0) {
            throw new IllegalArgumentException("quantity must be greater than zero. Was " + this.quantity);
        }
        this.dataType = dataType;
    }

    public static ProfinetTag of(String addressString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (!matcher.matches()) {
            throw new PlcInvalidTagException(addressString, ADDRESS_PATTERN);
        }

        String quantity = matcher.group("quantity") == null ? "1" :  matcher.group("quantity");
        PlcValueType plcValueType = PlcValueType.valueOf(matcher.group("datatype"));

        return new ProfinetTag(matcher.group("address"), Integer.parseInt(quantity), plcValueType);
    }

    @Override
    public String getAddressString() {
        return address;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return dataType;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return PlcTag.super.getArrayInfo();
    }
}
