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
package org.apache.plc4x.java.spi.codegen;

import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.WithReaderWriterArgs;

public interface WithOption extends WithReaderWriterArgs {

    static WithOption WithByteOrder(ByteOrder byteOrder) {
        return (withOptionByteOrder) () -> byteOrder;
    }

    static WithReaderWriterArgs WithEncoding(String encoding) {
        return WithReaderWriterArgs.WithEncoding(encoding);
    }

    static WithOption WithSerializationContext(String context) {
        return (withOptionSerializationContext) () -> context;
    }

    static WithReaderWriterArgs WithNullBytesHex(String nullBytesHex) {
        return WithReaderWriterArgs.WithNullBytesHex(nullBytesHex);
    }

}

interface withOptionByteOrder extends WithOption {
    ByteOrder byteOrder();
}

interface withOptionSerializationContext extends WithOption {
    String context();
}
