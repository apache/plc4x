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
package org.apache.plc4x.java.profinet.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ProfinetTag implements PlcTag {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?<slot>\\d+).(?<subSlot>\\d+).(?<direction>INPUT|OUTPUT)(.(?<index>\\d+))?:(?<dataType>[a-zA-Z_]+)(\\[(?<numElements>\\d+)])?");
    private final int slot;
    private final int subSlot;
    private final Direction direction;
    private final int index;
    private final PlcValueType dataType;
    private final int numElements;

    public ProfinetTag(int slot, int subSlot, Direction direction, int index, PlcValueType dataType, int numElements) {
        this.slot = slot;
        this.subSlot = subSlot;
        this.direction = direction;
        this.index = index;
        this.dataType = dataType;
        this.numElements = numElements;
        if (this.numElements <= 0) {
            throw new IllegalArgumentException("numElements must be greater than zero. Was " + this.numElements);
        }
    }

    public static ProfinetTag of(String addressString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (!matcher.matches()) {
            throw new PlcInvalidTagException(addressString, ADDRESS_PATTERN);
        }

        int slot = Integer.parseInt(matcher.group("slot"));
        int subSlot = Integer.parseInt(matcher.group("subSlot"));
        Direction direction = Direction.valueOf(matcher.group("direction"));
        int index = Integer.parseInt(matcher.group("index"));
        PlcValueType dataType = PlcValueType.valueOf(matcher.group("dataType"));
        int numElements = (matcher.group("numElements") != null) ? Integer.parseInt(matcher.group("numElements")) : 1;

        return new ProfinetTag(slot, subSlot, direction, index, dataType, numElements);
    }

    public int getSlot() {
        return slot;
    }

    public int getSubSlot() {
        return subSlot;
    }

    public Direction getDirection() {
        return direction;
    }

    public int getIndex() {
        return index;
    }

    public PlcValueType getDataType() {
        return dataType;
    }

    public int getNumElements() {
        return numElements;
    }

    @Override
    public String getAddressString() {
        return String.format("%d.%d.%s.%d:%s%s", slot, subSlot, direction, index, dataType, (numElements > 1) ? "[" + numElements + "]" : "");
    }

    @Override
    public PlcValueType getPlcValueType() {
        return dataType;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return PlcTag.super.getArrayInfo();
    }

    public static enum Direction {
        INPUT,
        OUTPUT
    }

}
