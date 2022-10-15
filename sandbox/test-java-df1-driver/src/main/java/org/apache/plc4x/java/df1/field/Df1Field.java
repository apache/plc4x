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

import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.df1.types.DataType;

public class Df1Field implements PlcField {

    private final int address;
    private final int size;
    private final DataType dataType;
    private final addressType address_type;

    public Df1Field(int address, int size, DataType dataType, addressType address_type) {
        this.address = address;
        this.size = size;
        this.dataType = dataType;
        this.address_type = address_type;
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

    public addressType getAddress_type() { return address_type; }

    public static PlcField of(String fieldQuery) {
        String[] tmp = fieldQuery.split(":");
        return new Df1Field(Integer.parseInt(tmp[0]), DataType.valueOf(tmp[1]).getLength(),DataType.valueOf(tmp[1]), addressType.OFFSET);
    }

    public enum addressType {
        OFFSET,
        LOGICAL
    }
}
