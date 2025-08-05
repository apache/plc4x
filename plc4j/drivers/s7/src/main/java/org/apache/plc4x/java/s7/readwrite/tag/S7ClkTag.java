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
package org.apache.plc4x.java.s7.readwrite.tag;

import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.s7.readwrite.DateAndTime;

import java.time.LocalDateTime;
import java.util.concurrent.TimeUnit;
import java.util.regex.Pattern;

/**
 * It seems reading CLK and CLKF are convenience ways of reading the PLCs Clock Data
 * - CLK seems to read the clock data via structured value (year, month, day, hour, ...)
 * - CLKF seems to read a formatted or extended format.
 */
public class S7ClkTag implements PlcTag {

    private static final Pattern CLK_ADDRESS_PATTERN =
        Pattern.compile("^CLK\\b");

    private static final Pattern CLKF_ADDRESS_PATTERN =
        Pattern.compile("^CLKF\\b");

    private final String address;
    private DateAndTime dat;

    public S7ClkTag(String address) {
        this.address = address;
        LocalDateTime ldt = LocalDateTime.now();
        this.dat = new DateAndTime(
            ((short) (ldt.getYear() - 2000)),
            (short) ldt.getMonthValue(),
            (short) ldt.getDayOfMonth(),
            (short) ldt.getHour(),
            (short) ldt.getMinute(),
            (short) ldt.getSecond(),
            (short) TimeUnit.NANOSECONDS.toMillis(ldt.getNano()),
            (byte) ldt.getDayOfWeek().plus(1).getValue()
        );
    }

    @Override
    public String getAddressString() {
        return address;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.DATE_AND_TIME;
    }

    public DateAndTime getDateAndTime() {
        return this.dat;
    }

    public void setDateAndTime(LocalDateTime ldt) {
        this.dat = new DateAndTime(
            ((short) (ldt.getYear() - 2000)),
            (short) ldt.getMonthValue(),
            (short) ldt.getDayOfMonth(),
            (short) ldt.getHour(),
            (short) ldt.getMinute(),
            (short) ldt.getSecond(),
            (short) TimeUnit.NANOSECONDS.toMillis(ldt.getNano()),
            (byte) ldt.getDayOfWeek().plus(1).getValue()
        );
    }

    public static boolean matches(String address) {
        return CLK_ADDRESS_PATTERN.matcher(address).matches() ||
                CLKF_ADDRESS_PATTERN.matcher(address).matches();
    }

    public static S7ClkTag of(String address) {
        return new S7ClkTag(address);

    }

}
