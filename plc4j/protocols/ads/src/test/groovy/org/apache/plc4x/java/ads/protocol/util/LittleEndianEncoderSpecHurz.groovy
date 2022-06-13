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
package org.apache.plc4x.java.ads.protocol.util

import org.apache.plc4x.java.ads.model.AdsDataType
import spock.lang.Specification
import spock.lang.Unroll

import static org.apache.plc4x.java.ads.model.AdsDataType.*

class LittleEndianEncoderSpecHurz extends Specification {
    @Unroll
    def "encode of #adsdt.name() using #values"(AdsDataType adsdt, def expectedBytes, def values) {
        when:
        def bytes = LittleEndianEncoder.encodeData(adsdt, *values)

        then:
        assert expectedBytes == bytes
        where:
        adsdt  | expectedBytes                                                                                                          | values
        SINT   | [0x01, 0x00, 0x01, 0x00] as byte[]                                                                                     | [true, false, true, false]

        SINT   | [0x12, 0x03, 0x05, 0x7f] as byte[]                                                                                     | [0x12, 0x03, 0x05, 0x7f] as byte[]
        INT    | [0x1, 0x00] as byte[]                                                                                                  | [1] as short[]
        INT    | [0x0e, 0x00, 0x50, 0x00] as byte[]                                                                                     | [14, 80] as short[]
        INT32  | [0x5a, 0x0a, 0x00, 0x00] as byte[]                                                                                     | [2650]
        INT32  | [0x5a, 0x0a, 0x00, 0x00, 0x50, 0x00, 0x00, 0x00] as byte[]                                                             | [2650, 80]
        REAL   | [0xc3, 0xf5, 0x48, 0x40] as byte[]                                                                                     | [3.14f]
        REAL   | [0xc3, 0xf5, 0x48, 0x40, 0x14, 0xae, 0x07, 0x40] as byte[]                                                             | [3.14f, 2.12f]
        LREAL  | [0x1F, 0x85, 0xEB, 0x51, 0xB8, 0x1E, 0x09, 0x40] as byte[]                                                             | [3.14d]
        LREAL  | [0x1F, 0x85, 0xEB, 0x51, 0xB8, 0x1E, 0x09, 0x40, 0xF6, 0x28, 0x5C, 0x8F, 0xC2, 0xF5, 0x00, 0x40] as byte[]             | [3.14d, 2.12d]
        STRING | [0x70, 0x6c, 0x63, 0x34, 0x78, 0x00] as byte[]                                                                         | ["plc4x"]
        STRING | [0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00] as byte[]                                     | ["HelloWorld!"]
        STRING | [0x70, 0x6c, 0x63, 0x34, 0x78, 0x00, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00] as byte[] | ["plc4x", "HelloWorld!"]
    }
}
