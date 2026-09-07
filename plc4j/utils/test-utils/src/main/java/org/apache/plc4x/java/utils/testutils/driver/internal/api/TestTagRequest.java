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


import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

public class TestTagRequest implements TestRequest {

    private final TestTag[] tags;

    public TestTagRequest(TestTag[] tags) {
        this.tags = tags;
    }

    public static TestTagRequest staticParse(ReadBuffer readBuffer, Object... args) throws BufferException {
        readBuffer.pushContext(WithOption.WithName("TestTagRequest"));

        int numberOfTags = readBuffer.readUnsignedInt(32, WithOption.WithName("numberOfTags"));
        readBuffer.pushContext(WithOption.WithName("tags"), WithOption.WithRenderAsList(true));
        TestTag[] tags = new TestTag[numberOfTags];
        for (int i = 0; i < numberOfTags; i++) {
            tags[i] = TestTag.staticParse(readBuffer, args);
        }
        readBuffer.popContext();

        readBuffer.popContext();
        return new TestTagRequest(tags);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        writeBuffer.pushContext(WithOption.WithName("TestTagRequest"));

        writeBuffer.writeUnsignedInt(32, tags.length, WithOption.WithName("numberOfTags"));
        writeBuffer.pushContext(WithOption.WithName("tags"), WithOption.WithRenderAsList(true));
        for (TestTag testTag : tags) {
            testTag.serialize(writeBuffer);
        }
        writeBuffer.popContext();

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

    public TestTag[] getTags() {
        return tags;
    }

}
