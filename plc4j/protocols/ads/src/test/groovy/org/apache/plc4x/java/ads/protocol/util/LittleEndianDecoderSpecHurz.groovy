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
package org.apache.plc4x.java.ads.protocol.util

import org.apache.plc4x.java.ads.model.AdsDataType
import org.apache.plc4x.java.api.exceptions.PlcProtocolException
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException
import spock.lang.Ignore
import spock.lang.Specification
import spock.lang.Unroll

import static org.apache.plc4x.java.ads.model.AdsDataType.*

class LittleEndianDecoderSpecHurz extends Specification {
    @Unroll
    def "decode of #adsdt.name() using get#retrievalType.simpleName [#expectedValues]"(AdsDataType adsdt, def retrievalType, def expectedValues, def rawData) {
        when:
        def plcValue = LittleEndianDecoder.decodeData(adsdt, rawData)

        and:
        def data = []
        (0..plcValue.numberOfValues - 1).forEach({ index ->
            data << plcValue."get${retrievalType.simpleName}"(index)
        })

        then:
        assert data == expectedValues
        where:
        adsdt  | retrievalType | expectedValues                                     | rawData
        BIT    | Boolean       | [true, false]                                      | [0x1, 0x0] as byte[]
        BIT    | Byte          | [1, 0]                                             | [0x1, 0x0] as byte[]
        BIT    | Short         | [1]                                                | [0x1] as byte[]
        BIT    | Short         | [1]                                                | [0x1] as byte[]
        BIT    | Short         | [1, 0]                                             | [0x1, 0x0] as byte[]

        INT    | Byte          | [1]                                                | [0x1, 0x0] as byte[]
        INT    | Short         | [1]                                                | [0x1, 0x0] as byte[]
        INT    | Short         | [256]                                              | [0x0, 0x1] as byte[]
        INT    | Short         | [256, 256]                                         | [0x0, 0x1, 0x0, 0x1] as byte[]
        INT    | Integer       | [1]                                                | [0x1, 0x0] as byte[]
        INT    | Integer       | [256]                                              | [0x0, 0x1] as byte[]
        INT    | Integer       | [256, 256]                                         | [0x0, 0x1, 0x0, 0x1] as byte[]
        INT    | Float         | [1.0f]                                             | [0x1, 0x0] as byte[]
        INT    | Float         | [256f]                                             | [0x0, 0x1] as byte[]
        INT    | Float         | [256f, 256f]                                       | [0x0, 0x1, 0x0, 0x1] as byte[]
        INT    | Double        | [1d]                                               | [0x1, 0x0] as byte[]
        INT    | Double        | [256d]                                             | [0x0, 0x1] as byte[]
        INT    | Double        | [256d, 256d]                                       | [0x0, 0x1, 0x0, 0x1] as byte[]

        INT32  | Integer       | [16777216]                                         | [0x0, 0x0, 0x0, 0x1] as byte[]
        INT32  | Integer       | [16777216, 16777216]                               | [0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1] as byte[]

        /*REAL   | Float         | [1.4E-45f]                                         | [0x1, 0x0, 0x0, 0x0] as byte[]
        REAL   | Float         | [2.3509887E-38f]                                   | [0x0, 0x0, 0x0, 0x1] as byte[]
        REAL   | Float         | [2.3509887E-38f, 2.3509887E-38f]                   | [0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1] as byte[]

        LREAL  | Double        | [4.9E-324]                                         | [0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0] as byte[]
        LREAL  | Double        | [7.2911220195563975E-304]                          | [0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1] as byte[]
        LREAL  | Double        | [7.2911220195563975E-304, 7.2911220195563975E-304] | [0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1] as byte[]*/

        STRING | String        | ["plc4x"]                                          | [0x70, 0x6c, 0x63, 0x34, 0x78, 0x0] as byte[]
        STRING | String        | ["plc4xplc4x"]                                     | [0x70, 0x6c, 0x63, 0x34, 0x78, 0x70, 0x6c, 0x63, 0x34, 0x78, 0x0] as byte[]
        STRING | String        | ["plc4x", "plc4x"]                                 | [0x70, 0x6c, 0x63, 0x34, 0x78, 0x0, 0x70, 0x6c, 0x63, 0x34, 0x78, 0x0] as byte[]
    }

    @Ignore("Needs finishing")
    def "failure test with string"() {
        when:
        LittleEndianDecoder.decodeData(STRING, [0x01] as byte[])

        then:
        thrown PlcProtocolException
    }

    def "failure test with unsupported data type"() {
        when:
        LittleEndianDecoder.decodeData(UNKNOWN, new byte[10])

        then:
        thrown PlcUnsupportedDataTypeException
    }
}
