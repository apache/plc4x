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

public class TestReadRequest implements Serializable {

    private final TestTagRequest testTagRequest;

    public TestReadRequest(TestTagRequest testTagRequest) {
        this.testTagRequest = testTagRequest;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("TestReadRequest");
        testTagRequest.serialize(writeBuffer);
        writeBuffer.popContext("TestReadRequest");
    }

    public static TestReadRequest staticParse(ReadBuffer readBuffer, Object... args) throws ParseException {
        readBuffer.pullContext("TestReadRequest");
        TestTagRequest testTagRequest = TestTagRequest.staticParse(readBuffer, args);
        readBuffer.closeContext("TestReadRequest");
        return new TestReadRequest(testTagRequest);
    }
}
