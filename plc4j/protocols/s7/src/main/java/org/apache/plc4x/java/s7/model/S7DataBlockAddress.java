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

import org.apache.plc4x.java.api.exceptions.PlcInvalidAddressException;
import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class S7DataBlockAddress extends S7Address {

    private static final Pattern ADDRESS_PATTERN =
        Pattern.compile("^DATA_BLOCKS/(?<blockNumber>\\d{1,4})/(?<byteOffset>\\d{1,4})");

    public static boolean matches(String addressString) {
        return ADDRESS_PATTERN.matcher(addressString).matches();
    }

    public static S7DataBlockAddress of(String addressString) throws PlcInvalidAddressException {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (!matcher.matches()) {
            throw new PlcInvalidAddressException(addressString, ADDRESS_PATTERN);
        }
        int datablockNumber = Integer.parseInt(matcher.group("blockNumber"));
        int datablockByteOffset = Integer.parseInt(matcher.group("byteOffset"));
        return new S7DataBlockAddress((short) datablockNumber, (short) datablockByteOffset);
    }

    private final short dataBlockNumber;

    public S7DataBlockAddress(short dataBlockNumber, short byteOffset) {
        super(MemoryArea.DATA_BLOCKS, byteOffset);
        this.dataBlockNumber = dataBlockNumber;
    }

    public short getDataBlockNumber() {
        return dataBlockNumber;
    }

}
