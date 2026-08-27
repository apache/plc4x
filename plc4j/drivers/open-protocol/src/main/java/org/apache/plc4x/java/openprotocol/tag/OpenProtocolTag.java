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
package org.apache.plc4x.java.openprotocol.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.List;

/**
 * Placeholder for an Open Protocol tag. This driver has no tag addressing yet: there is no
 * address syntax to parse and nothing for a tag to carry.
 *
 * <p>Every entry point says so rather than handing back something empty.
 * {@link OpenProtocolTagHandler#parseTag(String)} already threw; {@link #of(String)} returned
 * {@code null} for the same input, which the driver passed on to the caller as a tag - so the
 * failure surfaced later, somewhere else, as a {@code NullPointerException}.</p>
 */
public class OpenProtocolTag implements PlcTag {

    private static final String NOT_IMPLEMENTED =
        "The Open Protocol driver does not support tag addressing yet";

    /**
     * @throws PlcInvalidTagException always - see the class comment. Every other driver's
     *         {@code of()} throws this for an address it cannot parse, and this driver cannot
     *         parse any.
     */
    public static OpenProtocolTag of(String addressString) {
        throw new PlcInvalidTagException(NOT_IMPLEMENTED + ": '" + addressString + "'");
    }

    /**
     * @throws UnsupportedOperationException always - a tag of this type carries no address to
     *         render, and reporting an empty one would read as an address that is simply blank.
     */
    @Override
    public String getAddressString() {
        throw new UnsupportedOperationException(NOT_IMPLEMENTED);
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcTag.super.getPlcValueType();
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return PlcTag.super.getArrayInfo();
    }

}
