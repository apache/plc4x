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
package org.apache.plc4x.java.eip.readwrite.utils;

import org.apache.plc4x.java.eip.readwrite.IntegerEncoding;
import org.apache.plc4x.java.eip.readwrite.PathSegment;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;

import java.util.logging.Level;
import java.util.logging.Logger;



public class StaticHelper {

    /**
     * Tries to parse another PathSegment, if this works it returns false, if it doesn't it returns ture.
     * @param io readBuffer
     * @param order byte order
     * @return false if there's another PathSegment, true if not.
     */
    public static boolean noMorePathSegments(ReadBuffer io, IntegerEncoding order) {
        int initialPosition = io.getPos();
        try {
            PathSegment ps = PathSegment.staticParse(io, order);
            return false;
        } catch (Exception e) {
            return true;
        } finally {
            io.reset(initialPosition);
        }
    }

}
