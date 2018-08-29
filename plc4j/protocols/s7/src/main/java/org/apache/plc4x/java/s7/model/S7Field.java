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
package org.apache.plc4x.java.s7.model;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.apache.plc4x.java.s7.types.S7DataType;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class S7Field implements PlcField {

    // %I0.1:BOOL           <-- Digital Input
    // %IW64:REAL           <-- Analog Input
    // %Q0.4:BOOL           <-- Digital Output
    // %DB1.DBX38.1:BOOL    <-- Memory block DB1

    private static final Pattern ADDRESS_PATTERN =
        Pattern.compile("^%(?<memoryArea>.)(?<transferSizeCode>.?)(?<byteOffset>\\d{1,4})(.(?<bitOffset>[0-7]))?:(?<dataType>.+)");
    private static final Pattern DATA_BLOCK_ADDRESS_PATTERN =
        Pattern.compile("^%DB(?<blockNumber>\\d{1,4}).DB(?<transferSizeCode>.?)(?<byteOffset>\\d{1,4})(.(?<bitOffset>[0-7]))?:(?<dataType>.+)");

    public static boolean matches(String fieldString) {
        return DATA_BLOCK_ADDRESS_PATTERN.matcher(fieldString).matches() ||
            ADDRESS_PATTERN.matcher(fieldString).matches();
    }

    public static S7Field of(String fieldString) throws PlcInvalidFieldException {
        Matcher matcher = DATA_BLOCK_ADDRESS_PATTERN.matcher(fieldString);
        if(matcher.matches()) {
            S7DataType dataType = S7DataType.valueOf(matcher.group("dataType"));
            String transferSizeCode = matcher.group("transferSizeCode");
            if(!transferSizeCode.isEmpty() && !dataType.getSizeCode().equals(transferSizeCode)) {
                throw new PlcInvalidFieldException("Transfer size code '" + transferSizeCode +
                    "' doesn't match specified data type '" + dataType.name() + "'");
            }
            MemoryArea memoryArea = MemoryArea.INSTANCE_DATA_BLOCKS;
            short byteOffset = Short.parseShort(matcher.group("byteOffset"));
            short bitOffset = 0;
            if(matcher.group("bitOffset") != null) {
                bitOffset = Short.parseShort(matcher.group("bitOffset"));
            }
            return new S7Field(dataType, memoryArea, byteOffset, bitOffset);
        } else {
            matcher = ADDRESS_PATTERN.matcher(fieldString);
            if (matcher.matches()) {
                S7DataType dataType = S7DataType.valueOf(matcher.group("dataType"));
                String transferSizeCode = matcher.group("transferSizeCode");
                if(!transferSizeCode.isEmpty() && !dataType.getSizeCode().equals(transferSizeCode)) {
                    throw new PlcInvalidFieldException("Transfer size code '" + transferSizeCode +
                        "' doesn't match specified data type '" + dataType.name() + "'");
                }
                MemoryArea memoryArea = MemoryArea.valueOf(matcher.group("memoryArea"));
                short byteOffset = Short.parseShort(matcher.group("byteOffset"));
                short bitOffset = 0;
                if(matcher.group("bitOffset") != null) {
                    bitOffset = Short.parseShort(matcher.group("bitOffset"));
                }
                return new S7Field(dataType, memoryArea, byteOffset, bitOffset);
            }
        }
        throw new PlcInvalidFieldException("Unable to parse address: " + fieldString);
    }

    private final S7DataType dataType;
    private final MemoryArea memoryArea;
    private final short byteOffset;
    private final short bitOffset;

    public S7Field(S7DataType dataType, MemoryArea memoryArea, short byteOffset, short bitOffset) {
        this.dataType = dataType;
        this.memoryArea = memoryArea;
        this.byteOffset = byteOffset;
        this.bitOffset = bitOffset;
    }

    public S7DataType getDataType() {
        return dataType;
    }

    public MemoryArea getMemoryArea() {
        return memoryArea;
    }

    public short getByteOffset() {
        return byteOffset;
    }

    public short getBitOffset() {
        return bitOffset;
    }

}
