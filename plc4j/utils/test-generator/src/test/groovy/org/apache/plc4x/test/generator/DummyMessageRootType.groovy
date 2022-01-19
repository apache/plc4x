/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.test.generator

import org.apache.plc4x.java.spi.generation.ParseException
import org.apache.plc4x.java.spi.generation.ReadBuffer
import org.apache.plc4x.java.spi.generation.SerializationException
import org.apache.plc4x.java.spi.generation.WriteBuffer

class DummyMessageRootType {

    def dummyOutput

    static DummyMessageRootType staticParse(ReadBuffer readBuffer) throws ParseException {
        String dummyOutput = "<someXmlWithSomeContent>\n"
        int i = 0;
        while (true) {
            i++
            try {
                dummyOutput += "  <byte$i>${readBuffer.readByte()}</byte$i>\n"
            } catch (ignore) {
                break
            }
        }
        dummyOutput += "</someXmlWithSomeContent>\n"
        return new DummyMessageRootType(dummyOutput: dummyOutput)
    }

    void serialize(WriteBuffer writeBuffer) throws SerializationException {
        // TODO: we need to use this
    }

    @Override
    String toString() {
        return dummyOutput
    }
}
