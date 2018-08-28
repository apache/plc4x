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

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;

/**
 * Response to a {@link PlcReadRequest}.
 */
public interface PlcReadResponse extends PlcFieldResponse<PlcReadRequest> {

    int getNumValues(String name);

    boolean isRaw(String name);
    byte[] getRaw(String name);
    byte[] getRaw(String name, int index);

    boolean isObject(String name);
    Object getObject(String name);
    Object getObject(String name, int index);

    boolean isBoolean(String name);
    Boolean getBoolean(String name);
    Boolean getBoolean(String name, int index);

    boolean isByte(String name);
    Byte getByte(String name);
    Byte getByte(String name, int index);

    boolean isShort(String name);
    Short getShort(String name);
    Short getShort(String name, int index);

    boolean isInteger(String name);
    Integer getInteger(String name);
    Integer getInteger(String name, int index);

    boolean isLong(String name);
    Long getLong(String name);
    Long getLong(String name, int index);

    boolean isFloat(String name);
    Float getFloat(String name);
    Float getFloat(String name, int index);

    boolean isDouble(String name);
    Double getDouble(String name);
    Double getDouble(String name, int index);

    boolean isString(String name);
    String getString(String name);
    String getString(String name, int index);

    boolean isTime(String name);
    LocalTime getTime(String name);
    LocalTime getTime(String name, int index);

    boolean isDate(String name);
    LocalDate getDate(String name);
    LocalDate getDate(String name, int index);

    boolean isDateTime(String name);
    LocalDateTime getDateTime(String name);
    LocalDateTime getDateTime(String name, int index);

}
