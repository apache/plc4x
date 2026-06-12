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
package org.apache.plc4x.java.utils.testutils.driver.internal.api;

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

public class TestTag implements Message {

    private final String name;
    private final String address;

    public TestTag(String name, String address) {
        this.name = name;
        this.address = address;
    }

    public static TestTag staticParse(ReadBuffer readBuffer, Object... args) throws BufferException {
        readBuffer.pushContext(WithOption.WithName("TestTag"));
        String name = readBuffer.readString(64 * 8, WithOption.WithName("name")); // TODO: where to get the bitlength from
        String address = readBuffer.readString(64 * 8, WithOption.WithName("address")); // TODO: where to get the bitlength from
        readBuffer.popContext();
        return new TestTag(name, address);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        writeBuffer.pushContext(WithOption.WithName("TestTag"));
        writeBuffer.writeString(64 * 8, name, WithOption.WithName("name"));
        writeBuffer.writeString(64 * 8, name, WithOption.WithName("address"));
        writeBuffer.popContext();
    }

    @Override
    public int getLengthInBytes() {
        return getLengthInBits() / 8;
    }

    @Override
    public int getLengthInBits() {
        return 0;
    }

    public String getName() {
        return name;
    }

    public String getAddress() {
        return address;
    }

}
