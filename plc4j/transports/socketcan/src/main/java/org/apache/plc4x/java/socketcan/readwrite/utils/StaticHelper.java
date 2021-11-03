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
package org.apache.plc4x.java.socketcan.readwrite.utils;

public class StaticHelper {

    public static final int EFF_FLAG = 0b10000000_00000000_00000000_00000000;
    public static final int RTR_FLAG = 0b01000000_00000000_00000000_00000000;
    public static final int ERR_FLAG = 0b00100000_00000000_00000000_00000000;
    public static final int SFF_MASK = 0b00000000_00000000_00000111_11111111;
    public static final int EFF_MASK = 0b00011111_11111111_11111111_11111111;
    public static final int ERR_MASK = EFF_MASK;

    public static final int EXTENDED_FRAME_FORMAT_FLAG = 0x80000000;

    public static final int REMOTE_TRANSMISSION_FLAG = 0x40000000;

    public static final int ERROR_FRAME_FLAG = 0x20000000;

    public static final int STANDARD_FORMAT_IDENTIFIER_MASK = 0x7ff;

    public static final int EXTENDED_FORMAT_IDENTIFIER_MASK = 0x1fffffff;

    public static int readIdentifier(int identifier) {
        if ((isExtended(identifier))) {
            return identifier & EXTENDED_FORMAT_IDENTIFIER_MASK;
        }
        return identifier & STANDARD_FORMAT_IDENTIFIER_MASK;
    }

    public static boolean isExtended(int identifier) {
        return (identifier & EXTENDED_FRAME_FORMAT_FLAG) != 0;
    }

    public static boolean isRemote(int identifier) {
        return (identifier & REMOTE_TRANSMISSION_FLAG) != 0;
    }

    public static boolean isError(int identifier) {
        return (identifier & ERROR_FRAME_FLAG) != 0;
    }

}
