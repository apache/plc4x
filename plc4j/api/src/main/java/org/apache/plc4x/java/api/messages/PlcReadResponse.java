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
package org.apache.plc4x.java.api.messages;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Collection;

/**
 * Response to a {@link PlcReadRequest}.
 */
public interface PlcReadResponse extends PlcFieldResponse {

    @Override
    PlcReadRequest getRequest();

    int getNumberOfValues(String name);

    Object getObject(String name);

    Object getObject(String name, int index);

    Collection<Object> getAllObjects(String name);

    boolean isValidBoolean(String name);

    boolean isValidBoolean(String name, int index);

    Boolean getBoolean(String name);

    Boolean getBoolean(String name, int index);

    Collection<Boolean> getAllBooleans(String name);

    boolean isValidByte(String name);

    boolean isValidByte(String name, int index);

    Byte getByte(String name);

    Byte getByte(String name, int index);

    Collection<Byte> getAllBytes(String name);

    boolean isValidShort(String name);

    boolean isValidShort(String name, int index);

    Short getShort(String name);

    Short getShort(String name, int index);

    Collection<Short> getAllShorts(String name);

    boolean isValidInteger(String name);

    boolean isValidInteger(String name, int index);

    Integer getInteger(String name);

    Integer getInteger(String name, int index);

    Collection<Integer> getAllIntegers(String name);

    boolean isValidBigInteger(String name);

    boolean isValidBigInteger(String name, int index);

    BigInteger getBigInteger(String name);

    BigInteger getBigInteger(String name, int index);

    Collection<BigInteger> getAllBigIntegers(String name);

    boolean isValidLong(String name);

    boolean isValidLong(String name, int index);

    Long getLong(String name);

    Long getLong(String name, int index);

    Collection<Long> getAllLongs(String name);

    boolean isValidFloat(String name);

    boolean isValidFloat(String name, int index);

    Float getFloat(String name);

    Float getFloat(String name, int index);

    Collection<Float> getAllFloats(String name);

    boolean isValidDouble(String name);

    boolean isValidDouble(String name, int index);

    Double getDouble(String name);

    Double getDouble(String name, int index);

    Collection<Double> getAllDoubles(String name);

    boolean isValidBigDecimal(String name);

    boolean isValidBigDecimal(String name, int index);

    BigDecimal getBigDecimal(String name);

    BigDecimal getBigDecimal(String name, int index);

    Collection<BigDecimal> getAllBigDecimals(String name);

    boolean isValidString(String name);

    boolean isValidString(String name, int index);

    String getString(String name);

    String getString(String name, int index);

    Collection<String> getAllStrings(String name);

    boolean isValidTime(String name);

    boolean isValidTime(String name, int index);

    LocalTime getTime(String name);

    LocalTime getTime(String name, int index);

    Collection<LocalTime> getAllTimes(String name);

    boolean isValidDate(String name);

    boolean isValidDate(String name, int index);

    LocalDate getDate(String name);

    LocalDate getDate(String name, int index);

    Collection<LocalDate> getAllDates(String name);

    boolean isValidDateTime(String name);

    boolean isValidDateTime(String name, int index);

    LocalDateTime getDateTime(String name);

    LocalDateTime getDateTime(String name, int index);

    Collection<LocalDateTime> getAllDateTimes(String name);

    boolean isValidByteArray(String name);

    boolean isValidByteArray(String name, int index);

    Byte[] getByteArray(String name);

    Byte[] getByteArray(String name, int index);

    Collection<Byte[]> getAllByteArrays(String name);

}
