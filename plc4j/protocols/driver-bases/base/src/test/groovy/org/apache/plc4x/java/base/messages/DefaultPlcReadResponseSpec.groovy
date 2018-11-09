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

package org.apache.plc4x.java.base.messages

import org.apache.commons.lang3.tuple.ImmutablePair
import org.apache.plc4x.java.api.exceptions.PlcFieldRangeException
import org.apache.plc4x.java.api.types.PlcResponseCode
import org.apache.plc4x.java.base.messages.items.*
import spock.lang.Specification
import spock.lang.Unroll

import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime

class DefaultPlcReadResponseSpec extends Specification {

    @Unroll
    def "Using the different types of getters should work '#dataType'"(def dataType, def fieldType, def fieldValues) {
        setup:
        InternalPlcReadRequest request = Mock(InternalPlcReadRequest)
        request.getFieldNames() >> ['foo']

        when:
        DefaultPlcReadResponse SUT = new DefaultPlcReadResponse(request,
            ["foo": new ImmutablePair<>(PlcResponseCode.OK, fieldType.newInstance(*fieldValues))])

        and:
        def fieldNames = SUT.getFieldNames()
        // The time fields have a "Local" added to their type names, so we must cut that off.
        def name = dataType.simpleName.replaceAll("Local", "")
        def valid = SUT."isValid${name}"("foo")
        def numberOfValues = SUT.getNumberOfValues("foo")
        def responseCode = SUT.getResponseCode("foo")
        def firstValue = SUT."get${name}"("foo")
        def secondValue = SUT."get${name}"("foo", 1)
        def allValues = SUT."getAll${name}s"("foo")
        SUT.getObject("foo")
        SUT.getAllObjects("foo")
        SUT."get${name}"("foo", 2)

        then:
        assert fieldNames == ['foo'] as LinkedHashSet
        assert valid
        assert numberOfValues == fieldValues.length
        assert responseCode == PlcResponseCode.OK
        assert fieldValues[0] == firstValue
        assert fieldValues[1] == secondValue
        assert fieldValues == allValues
        thrown PlcFieldRangeException

        where:
        dataType      | fieldType                     | fieldValues
        Boolean       | DefaultBooleanFieldItem       | [true, false] as boolean[]
        Byte          | DefaultByteFieldItem          | [42, 23] as byte[]
        Short         | DefaultShortFieldItem         | [42, 23] as short[]
        Integer       | DefaultIntegerFieldItem       | [42, 23] as int[]
        Long          | DefaultLongFieldItem          | [42, 23] as long[]
        BigInteger    | DefaultBigDecimalFieldItem    | [42, 23] as BigInteger[]
        Float         | DefaultFloatFieldItem         | [42, 23] as float[]
        Double        | DefaultDoubleFieldItem        | [42, 23] as double[]
        BigDecimal    | DefaultBigDecimalFieldItem    | [42, 23] as BigDecimal[]
        String        | DefaultStringFieldItem        | ["foo", "bar"] as String[]
        LocalTime     | DefaultLocalTimeFieldItem     | [LocalTime.NOON, LocalTime.MIDNIGHT] as LocalTime[]
        LocalDate     | DefaultLocalDateFieldItem     | [LocalDate.MIN, LocalDate.MAX] as LocalDate[]
        LocalDateTime | DefaultLocalDateTimeFieldItem | [LocalDateTime.MIN, LocalDateTime.MAX] as LocalDateTime[]
    }

}
