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

package org.apache.plc4x.java.api.value;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;

public class PlcDateTime extends PlcSimpleValue<LocalDateTime> {

    public PlcDateTime(LocalDateTime value) {
        super(value, true);
    }

    @Override
    public boolean isString() {
        return true;
    }

    @Override
    public String getString() {
        return value.toString();
    }

    @Override
    public boolean isTime() {
        return true;
    }

    @Override
    public LocalTime getTime() {
        return value.toLocalTime();
    }

    @Override
    public boolean isDate() {
        return true;
    }

    @Override
    public LocalDate getDate() {
        return value.toLocalDate();
    }

    @Override
    public boolean isDateTime() {
        return true;
    }

    @Override
    public LocalDateTime getDateTime() {
        return value;
    }

    @Override
    public String toString() {
        return String.valueOf(value);
    }

}
