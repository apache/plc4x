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
package org.apache.plc4x.java.api.model;

public interface ArrayInfo {

    /**
     * @return Number of elements in total
     */
    int getSize();

    /**
     * The lower index of the selection, as it was written in the address. For a single element
     * such as [6] this is 6, and {@link #getUpperBound()} is 6 as well - a bare index selects
     * one element, not a range starting at zero.
     * @return Returns the index of lower bound of the array.
     */
    int getLowerBound();

    /**
     * The upper index of the selection, as it was written in the address. For the range [0..7]
     * this is 7 and {@link #getSize()} is 8, both bounds being inclusive.
     * @return Returns the index of upper bound of the array.
     */
    int getUpperBound();

    /**
     * The array's declared lower bound, as in PLCs not every array starts at 0. An address may
     * state it explicitly - [4..7;1] selects elements 4 to 7 of an array declared from 1 - so
     * that the bounds above can be written the way the PLC program declares them. The offset of
     * an element from the start of the array is its index minus this value.
     *
     * <p>Defaults to 0, which is correct for any array that does not declare otherwise.
     *
     * @return Returns the index the array is declared to start at.
     */
    default int getBase() {
        return 0;
    }

    /**
     * Whether the address wrote this dimension as a range rather than a single index.
     *
     * <p>The two mean different things to a caller: a single index selects one element and yields
     * a scalar, while a range yields an array - even a range spanning one element. Equal bounds
     * alone cannot tell them apart, so the written form has to be remembered.
     *
     * @return true when the dimension was written as a range.
     */
    default boolean isRange() {
        return getLowerBound() != getUpperBound();
    }

}
