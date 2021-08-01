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
package org.apache.plc4x.java.ads.api.tcp.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.ByteValue;

import java.nio.charset.Charset;

import static java.util.Objects.requireNonNull;

public class UserData extends ByteValue {

    public static final UserData EMPTY = UserData.of();

    private UserData(byte... values) {
        super(values);
    }

    public UserData(ByteBuf byteBuf) {
        super(byteBuf);
    }

    public static UserData of(byte... values) {
        return new UserData(values);
    }

    public static UserData of(String value) {
        requireNonNull(value);
        return new UserData(value.getBytes());
    }

    public static UserData of(String value, Charset charset) {
        requireNonNull(value);
        return new UserData(value.getBytes(charset));
    }

    public static UserData of(ByteBuf byteBuf) {
        return new UserData(byteBuf);
    }

    @Override
    public String toString() {
        // TODO: maybe we could find a way to implement this to string
        return "UserData{" + value.length + "bytes} ";
    }
}