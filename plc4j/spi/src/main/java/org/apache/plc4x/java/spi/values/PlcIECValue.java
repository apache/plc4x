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

package org.apache.plc4x.java.spi.values;

import com.fasterxml.jackson.annotation.JsonIgnore;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.w3c.dom.Document;
import org.w3c.dom.Element;

import java.nio.charset.StandardCharsets;

public abstract class PlcIECValue<T> extends PlcValueAdapter {

    T value;
    boolean isNullable;

    PlcIECValue() {
        this.value = null;
        this.isNullable = true;
    }

    @Override
    @JsonIgnore
    public Object getObject() {
        return value;
    }

    @Override
    @JsonIgnore
    public int getLength() {
        return 1;
    }

    @Override
    @JsonIgnore
    public boolean isSimple() {
        return true;
    }

    @Override
    @JsonIgnore
    public boolean isNullable() {
        return isNullable;
    }

    @Override
    @JsonIgnore
    public boolean isNull() {
        return isNullable && value == null;
    }

    /**
     * Convenience method to simplify accessing items with a list syntax.
     * @param i item number
     * @return if i == 0 returns itself, otherwise throws an exception.
     */
    @Override
    @JsonIgnore
    public PlcValue getIndex(int i) {
        if(i == 0) {
            return this;
        }
        return super.getIndex(i);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws ParseException {
        String valueString = value.toString();
        writeBuffer.writeString(getClass().getSimpleName(), valueString.getBytes(StandardCharsets.UTF_8).length*8,StandardCharsets.UTF_8.name(),valueString);
    }

}
