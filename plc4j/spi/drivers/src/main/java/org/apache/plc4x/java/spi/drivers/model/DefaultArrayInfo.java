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
package org.apache.plc4x.java.spi.drivers.model;

import org.apache.plc4x.java.api.model.ArrayInfo;

public class DefaultArrayInfo implements ArrayInfo {

    private final int lowerBound;
    private final int upperBound;
    private final int base;
    private final boolean range;

    /** An array declared from 0, which is every array that does not say otherwise. */
    public DefaultArrayInfo(int lowerBound, int upperBound) {
        this(lowerBound, upperBound, 0);
    }

    public DefaultArrayInfo(int lowerBound, int upperBound, int base) {
        this(lowerBound, upperBound, base, lowerBound != upperBound);
    }

    /**
     * @param range whether the address wrote this dimension as a range. A one-element range is
     *              still a range - see {@link ArrayInfo#isRange()}.
     */
    public DefaultArrayInfo(int lowerBound, int upperBound, int base, boolean range) {
        this.lowerBound = lowerBound;
        this.upperBound = upperBound;
        this.base = base;
        this.range = range;
    }

    @Override
    public int getSize() {
        return upperBound - lowerBound + 1;
    }

    @Override
    public int getLowerBound() {
        return lowerBound;
    }

    @Override
    public int getUpperBound() {
        return upperBound;
    }

    @Override
    public int getBase() {
        return base;
    }

    @Override
    public boolean isRange() {
        return range;
    }

    /** The offset of the first selected element from the start of the array. */
    public int getOffset() {
        return lowerBound - base;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof DefaultArrayInfo other)) {
            return false;
        }
        return lowerBound == other.lowerBound && upperBound == other.upperBound
            && base == other.base && range == other.range;
    }

    @Override
    public int hashCode() {
        return java.util.Objects.hash(lowerBound, upperBound, base, range);
    }

    @Override
    public String toString() {
        return "DefaultArrayInfo{lowerBound=" + lowerBound
            + ", upperBound=" + upperBound
            + ", base=" + base
            + ", range=" + range + '}';
    }

}
