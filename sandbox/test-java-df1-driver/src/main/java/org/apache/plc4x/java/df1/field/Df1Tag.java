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
package org.apache.plc4x.java.df1.field;

import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.df1.types.DataType;

import java.util.Collections;
import java.util.List;

public class Df1Tag implements PlcTag {

    private final int address;
    private final int size;
    private final DataType dataType;
    private final AddressType addressType;

    public Df1Tag(int address, int size, DataType dataType, AddressType addressType) {
        this.address = address;
        this.size = size;
        this.dataType = dataType;
        this.addressType = addressType;
    }

    @Override
    public String getAddressString() {
        return String.format("%d:%s", address, addressType.toString());
    }

    @Override
    public PlcValueType getPlcValueType() {
        return dataType.getPlcValueType();
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        // TODO: This might need some work...
        return Collections.emptyList();
    }

    public int getAddress() {
        return address;
    }

    public int getSize() {
        return size;
    }

    public DataType getDataType() {
        return dataType;
    }

    public AddressType getAddressType() { return addressType; }

    public static PlcTag of(String fieldQuery) {
        String[] tmp = fieldQuery.split(":");
        return new Df1Tag(Integer.parseInt(tmp[0]), DataType.valueOf(tmp[1]).getLength(),DataType.valueOf(tmp[1]), AddressType.OFFSET);
    }

    public enum AddressType {
        OFFSET,
        LOGICAL
    }
}
