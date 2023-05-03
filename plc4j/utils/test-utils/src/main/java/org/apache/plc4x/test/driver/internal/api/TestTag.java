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
package org.apache.plc4x.test.driver.internal.api;

import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

public class TestTag implements Serializable {

    private final String name;
    private final String address;

    public TestTag(String name, String address) {
        this.name = name;
        this.address = address;
    }

    public static TestTag staticParse(ReadBuffer readBuffer, Object... args) throws ParseException {
        readBuffer.pullContext("TestTag");
        String name = readBuffer.readString("name", 64 * 8); // TODO: where to get the bitlength from
        String address = readBuffer.readString("address", 64 * 8); // TODO: where to get the bitlength from
        readBuffer.closeContext("TestTag");
        return new TestTag(name, address);
    }

    public String getName() {
        return name;
    }

    public String getAddress() {
        return address;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("TestTag");
        writeBuffer.writeString("name", 64 * 8, name);
        writeBuffer.writeString("address", 64 * 8, name);
        writeBuffer.popContext("TestTag");
    }
}
