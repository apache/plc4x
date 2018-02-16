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
package org.apache.plc4x.java.ads.api.generic.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.UnsignedIntLEByteValue;

public class Invoke extends UnsignedIntLEByteValue {

    public static final int NUM_BYTES = UnsignedIntLEByteValue.NUM_BYTES;

    public static final Invoke NONE = of(0);

    protected Invoke(byte... values) {
        super(values);
    }

    protected Invoke(long value) {
        super(value);
    }

    protected Invoke(ByteBuf byteBuf) {
        super(byteBuf);
    }

    public static Invoke of(byte... values) {
        return new Invoke(values);
    }

    public static Invoke of(long errorCode) {
        checkUnsignedBounds(errorCode, NUM_BYTES);
        return new Invoke(errorCode);
    }

    public static Invoke of(ByteBuf byteBuf) {
        return new Invoke(byteBuf);
    }
}
