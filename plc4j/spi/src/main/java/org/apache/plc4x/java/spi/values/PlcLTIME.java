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
package org.apache.plc4x.java.spi.values;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.w3c.dom.Element;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.temporal.ChronoUnit;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcLTIME extends PlcSimpleValue<Duration> {

    public static PlcLTIME of(Object value) {
        if (value instanceof Duration) {
            return new PlcLTIME((Duration) value);
        } else if(value instanceof Integer) {
            return new PlcLTIME(Duration.of((long) value, ChronoUnit.MILLIS));
        } else if(value instanceof Long) {
            return new PlcLTIME(Duration.of((long) value, ChronoUnit.NANOS));
        }
        throw new PlcRuntimeException("Invalid value type");
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcLTIME(@JsonProperty("value") Duration value) {
        super(value, true);
    }

    public PlcLTIME(@JsonProperty("value") Integer value) {
        super(Duration.of((long) value, ChronoUnit.NANOS), true);
    }

    public PlcLTIME(@JsonProperty("value") Long value) {
        super(Duration.of(value, ChronoUnit.NANOS), true);
    }

    public PlcLTIME(@JsonProperty("value") BigInteger value) {
        super(Duration.of(value.longValue(), ChronoUnit.NANOS), true);
    }

    @Override
    @JsonIgnore
    public boolean isString() {
        return true;
    }

    @Override
    public boolean isDuration() {
        return true;
    }

    @Override
    public Duration getDuration() {
        return value;
    }

    @Override
    @JsonIgnore
    public String getString() {
        return value.toString();
    }

    @Override
    @JsonIgnore
    public String toString() {
        return String.valueOf(value);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws ParseException {
        String valueString = value.toString();
        writeBuffer.writeString(getClass().getSimpleName(), valueString.getBytes(StandardCharsets.UTF_8).length*8,StandardCharsets.UTF_8.name(),valueString);
    }

}
