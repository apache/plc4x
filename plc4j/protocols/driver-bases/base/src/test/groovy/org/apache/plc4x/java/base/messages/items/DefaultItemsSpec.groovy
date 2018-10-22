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
package org.apache.plc4x.java.base.messages.items


import org.apache.plc4x.java.api.exceptions.PlcIncompatibleDatatypeException
import spock.lang.Specification
import spock.lang.Unroll

import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime

class DefaultItemsSpec extends Specification {

    @Unroll
    def "The '#fieldItemType.simpleName'.isValidXYZ methods should respect the correct boundaries of the type XYZ and "(
        Class<? extends BaseDefaultFieldItem> fieldItemType, Object value, Boolean isValidBoolean, Boolean isValidByte, Boolean isValidShort, Boolean isValidInteger,
        Boolean isValidLong, Boolean isValidBigInteger, Boolean isValidFloat, Boolean isValidDouble, Boolean isValidBigDecimal,
        Boolean isValidString, Boolean isValidTime, Boolean isValidDate, Boolean isValidDateTime, Boolean isValidByteArray) {

        setup:
        BaseDefaultFieldItem fieldItem = fieldItemType.newInstance(value)

        expect:
        fieldItem.getObject(0) == value

        assertItem(fieldItem, "Boolean", isValidBoolean)
        assertItem(fieldItem, "Byte", isValidByte)
        assertItem(fieldItem, "Short", isValidShort)
        assertItem(fieldItem, "Integer", isValidInteger)
        assertItem(fieldItem, "Long", isValidLong)
        assertItem(fieldItem, "BigInteger", isValidBigInteger)
        assertItem(fieldItem, "Float", isValidFloat)
        assertItem(fieldItem, "Double", isValidDouble)
        assertItem(fieldItem, "BigDecimal", isValidBigDecimal)
        assertItem(fieldItem, "String", isValidString)
        assertItem(fieldItem, "Time", isValidTime)
        assertItem(fieldItem, "Date", isValidDate)
        assertItem(fieldItem, "DateTime", isValidDateTime)
        assertItem(fieldItem, "ByteArray", isValidByteArray)

        where:
        fieldItemType              | value                                      || isValidBoolean | isValidByte | isValidShort | isValidInteger | isValidLong | isValidBigInteger | isValidFloat | isValidDouble | isValidBigDecimal | isValidString | isValidTime | isValidDate | isValidDateTime | isValidByteArray
        DefaultBooleanFieldItem    | null                                       || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | false       | false       | false           | false
        DefaultBooleanFieldItem    | true                                       || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBooleanFieldItem    | false                                      || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
                                                                                
//      fieldItemType              | value                                      || isValidBoolean | isValidByte | isValidShort | isValidInteger | isValidLong | isValidBigInteger | isValidFloat | isValidDouble | isValidBigDecimal | isValidString | isValidTime | isValidDate | isValidDateTime | isValidByteArray
        DefaultByteFieldItem       | null                                       || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | false       | false       | false           | false
        DefaultByteFieldItem       | (byte) 0                                   || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultByteFieldItem       | (byte) 42                                  || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultByteFieldItem       | Byte.MAX_VALUE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultByteFieldItem       | Byte.MIN_VALUE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false

//      fieldItemType              | value                                      || isValidBoolean | isValidByte | isValidShort | isValidInteger | isValidLong | isValidBigInteger | isValidFloat | isValidDouble | isValidBigDecimal | isValidString | isValidTime | isValidDate | isValidDateTime | isValidByteArray
        DefaultShortFieldItem      | null                                       || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | false       | false       | false           | false
        DefaultShortFieldItem      | (short) 0                                  || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultShortFieldItem      | (short) 42                                 || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultShortFieldItem      | Byte.MIN_VALUE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultShortFieldItem      | Byte.MAX_VALUE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultShortFieldItem      | Short.MAX_VALUE                            || true           | false       | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultShortFieldItem      | Short.MIN_VALUE                            || true           | false       | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false

//      fieldItemType              | value                                      || isValidBoolean | isValidByte | isValidShort | isValidInteger | isValidLong | isValidBigInteger | isValidFloat | isValidDouble | isValidBigDecimal | isValidString | isValidTime | isValidDate | isValidDateTime | isValidByteArray
        DefaultIntegerFieldItem    | null                                       || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | false       | false       | false           | false
        DefaultIntegerFieldItem    | (int) 0                                    || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultIntegerFieldItem    | (int) 42                                   || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultIntegerFieldItem    | Byte.MIN_VALUE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultIntegerFieldItem    | Byte.MAX_VALUE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultIntegerFieldItem    | Short.MAX_VALUE                            || true           | false       | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultIntegerFieldItem    | Short.MIN_VALUE                            || true           | false       | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultIntegerFieldItem    | Integer.MAX_VALUE                          || true           | false       | false        | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultIntegerFieldItem    | Integer.MIN_VALUE                          || true           | false       | false        | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
                                                                                
//      fieldItemType              | value                                      || isValidBoolean | isValidByte | isValidShort | isValidInteger | isValidLong | isValidBigInteger | isValidFloat | isValidDouble | isValidBigDecimal | isValidString | isValidTime | isValidDate | isValidDateTime | isValidByteArray
        DefaultLongFieldItem       | null                                       || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | false       | false       | false           | false
        DefaultLongFieldItem       | (int) 0                                    || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultLongFieldItem       | (int) 42                                   || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultLongFieldItem       | Byte.MIN_VALUE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultLongFieldItem       | Byte.MAX_VALUE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultLongFieldItem       | Short.MAX_VALUE                            || true           | false       | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultLongFieldItem       | Short.MIN_VALUE                            || true           | false       | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultLongFieldItem       | Integer.MAX_VALUE                          || true           | false       | false        | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultLongFieldItem       | Integer.MIN_VALUE                          || true           | false       | false        | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultLongFieldItem       | Long.MAX_VALUE                             || true           | false       | false        | false          | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultLongFieldItem       | Long.MIN_VALUE                             || true           | false       | false        | false          | true        | true              | true         | true          | true              | false         | false       | false       | false           | false

//      fieldItemType              | value                                      || isValidBoolean | isValidByte | isValidShort | isValidInteger | isValidLong | isValidBigInteger | isValidFloat | isValidDouble | isValidBigDecimal | isValidString | isValidTime | isValidDate | isValidDateTime | isValidByteArray
        DefaultBigIntegerFieldItem | null                                       || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | false       | false       | false           | false
        DefaultBigIntegerFieldItem | (int) 0                                    || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigIntegerFieldItem | (int) 42                                   || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigIntegerFieldItem | Byte.MIN_VALUE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigIntegerFieldItem | Byte.MAX_VALUE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigIntegerFieldItem | Short.MIN_VALUE                            || true           | false       | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigIntegerFieldItem | Short.MAX_VALUE                            || true           | false       | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigIntegerFieldItem | Integer.MIN_VALUE                          || true           | false       | false        | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigIntegerFieldItem | Integer.MAX_VALUE                          || true           | false       | false        | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigIntegerFieldItem | Long.MIN_VALUE                             || true           | false       | false        | false          | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigIntegerFieldItem | Long.MAX_VALUE                             || true           | false       | false        | false          | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigIntegerFieldItem | BigInteger.valueOf(Long.MIN_VALUE).multiply(2) || true           | false       | false        | false          | false       | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigIntegerFieldItem | BigInteger.valueOf(Long.MAX_VALUE).multiply(2) || true           | false       | false        | false          | false       | true              | true         | true          | true              | false         | false       | false       | false           | false

//      fieldItemType              | value                                      || isValidBoolean | isValidByte | isValidShort | isValidInteger | isValidLong | isValidBigInteger | isValidFloat | isValidDouble | isValidBigDecimal | isValidString | isValidTime | isValidDate | isValidDateTime | isValidByteArray
        DefaultFloatFieldItem      | null                                       || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | false       | false       | false           | false
        DefaultFloatFieldItem      | (int) 0                                    || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultFloatFieldItem      | (int) 42                                   || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultFloatFieldItem      | Byte.MIN_VALUE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultFloatFieldItem      | Byte.MAX_VALUE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultFloatFieldItem      | Short.MIN_VALUE                            || true           | false       | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultFloatFieldItem      | Short.MAX_VALUE                            || true           | false       | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultFloatFieldItem      | Integer.MIN_VALUE                          || true           | false       | false        | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
// In this case the conversion to Float results in a rounding error
//        DefaultFloatFieldItem      | Integer.MAX_VALUE                          || true           | false       | false        | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultFloatFieldItem      | Long.MIN_VALUE                             || true           | false       | false        | false          | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultFloatFieldItem      | Long.MAX_VALUE                             || true           | false       | false        | false          | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultFloatFieldItem      | 1.2345678912345f                           || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultFloatFieldItem      | Float.MIN_VALUE                            || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultFloatFieldItem      | -Float.MAX_VALUE                           || true           | false       | false        | false          | false       | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultFloatFieldItem      | Float.MAX_VALUE                            || true           | false       | false        | false          | false       | true              | true         | true          | true              | false         | false       | false       | false           | false

//      fieldItemType              | value                                      || isValidBoolean | isValidByte | isValidShort | isValidInteger | isValidLong | isValidBigInteger | isValidFloat | isValidDouble | isValidBigDecimal | isValidString | isValidTime | isValidDate | isValidDateTime | isValidByteArray
        DefaultDoubleFieldItem     | null                                       || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | (int) 0                                    || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | (int) 42                                   || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | Byte.MIN_VALUE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | Byte.MAX_VALUE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | Short.MIN_VALUE                            || true           | false       | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | Short.MAX_VALUE                            || true           | false       | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | Integer.MIN_VALUE                          || true           | false       | false        | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | Integer.MAX_VALUE                          || true           | false       | false        | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | Long.MIN_VALUE                             || true           | false       | false        | false          | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | Long.MAX_VALUE                             || true           | false       | false        | false          | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | 1.23456789123456f                          || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | Float.MIN_VALUE                            || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | -Float.MAX_VALUE                           || true           | false       | false        | false          | false       | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | Float.MAX_VALUE                            || true           | false       | false        | false          | false       | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | 1.23456789012345d                          || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | Double.MIN_VALUE                           || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | -Double.MAX_VALUE                          || true           | false       | false        | false          | false       | true              | false        | true          | true              | false         | false       | false       | false           | false
        DefaultDoubleFieldItem     | Double.MAX_VALUE                           || true           | false       | false        | false          | false       | true              | false        | true          | true              | false         | false       | false       | false           | false

//      fieldItemType              | value                                      || isValidBoolean | isValidByte | isValidShort | isValidInteger | isValidLong | isValidBigInteger | isValidFloat | isValidDouble | isValidBigDecimal | isValidString | isValidTime | isValidDate | isValidDateTime | isValidByteArray
        DefaultBigDecimalFieldItem | null                                       || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.ZERO                            || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.ONE                             || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(42)                     || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(Byte.MIN_VALUE)         || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(Byte.MAX_VALUE)         || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(Short.MIN_VALUE)        || true           | false       | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(Short.MAX_VALUE)        || true           | false       | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(Integer.MIN_VALUE)      || true           | false       | false        | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(Integer.MAX_VALUE)      || true           | false       | false        | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(Long.MIN_VALUE)         || true           | false       | false        | false          | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(Long.MAX_VALUE)         || true           | false       | false        | false          | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(1.23456789123456f)      || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(Float.MIN_VALUE)        || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(-Float.MAX_VALUE)       || true           | false       | false        | false          | false       | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(Float.MAX_VALUE)        || true           | false       | false        | false          | false       | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(1.23456789012345d)      || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(Double.MIN_VALUE)       || true           | true        | true         | true           | true        | true              | true         | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(-Double.MAX_VALUE)      || true           | false       | false        | false          | false       | true              | false        | true          | true              | false         | false       | false       | false           | false
        DefaultBigDecimalFieldItem | BigDecimal.valueOf(Double.MAX_VALUE)       || true           | false       | false        | false          | false       | true              | false        | true          | true              | false         | false       | false       | false           | false

//      fieldItemType              | value                                      || isValidBoolean | isValidByte | isValidShort | isValidInteger | isValidLong | isValidBigInteger | isValidFloat | isValidDouble | isValidBigDecimal | isValidString | isValidTime | isValidDate | isValidDateTime | isValidByteArray
        DefaultLocalTimeFieldItem  | null                                       || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | false       | false       | false           | false
        DefaultLocalTimeFieldItem  | LocalTime.now()                            || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | true        | false       | false           | false

//      fieldItemType              | value                                      || isValidBoolean | isValidByte | isValidShort | isValidInteger | isValidLong | isValidBigInteger | isValidFloat | isValidDouble | isValidBigDecimal | isValidString | isValidTime | isValidDate | isValidDateTime | isValidByteArray
        DefaultLocalDateFieldItem  | null                                       || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | false       | false       | false           | false
        DefaultLocalDateFieldItem  | LocalDate.now()                            || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | false       | true        | false           | false

//      fieldItemType              | value                                      || isValidBoolean | isValidByte | isValidShort | isValidInteger | isValidLong | isValidBigInteger | isValidFloat | isValidDouble | isValidBigDecimal | isValidString | isValidTime | isValidDate | isValidDateTime | isValidByteArray
        DefaultLocalDateTimeFieldItem | null                                    || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | false       | false       | false           | false
        DefaultLocalDateTimeFieldItem | LocalDateTime.now()                     || false          | false       | false        | false          | false       | false             | false        | false         | false             | false         | true        | true        | true            | false
    }

    Boolean assertItem(BaseDefaultFieldItem fieldItem, String type, Boolean expectedToBeValid) {
        assert fieldItem."isValid$type"(0) == expectedToBeValid
        if (expectedToBeValid) {
            assert fieldItem."get$type"(0) != null
        } else {
            assert getExecutionException({fieldItem."get$type"(0)}) instanceof PlcIncompatibleDatatypeException
        }
        return true
    }

    Exception getExecutionException(Closure c) {
        try{
            c.call()
            return null
        }catch(Exception e){
            return e
        }
    }

}


