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
package org.apache.plc4x.java.spi.codegen.io;

import org.apache.plc4x.java.spi.generation.*;

public class DataReaderSimpleByte implements DataReaderSimple<Byte> {

    private final ReadBuffer readBuffer;

    public DataReaderSimpleByte(ReadBuffer readBuffer) {
        this.readBuffer = readBuffer;
    }

    @Override
    public Byte read(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException  {
        if(bitLength != 8) {
            throw new ParseException("Byte fields only support bitLength of 8");
        }
        return readBuffer.readByte(logicalName, readerArgs);
    }

    public int getPos() {
        return readBuffer.getPos();
    }

    public void setPos(int position) {
        readBuffer.reset(position);
    }
}
