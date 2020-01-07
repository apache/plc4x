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

package org.apache.plc4x.java.s7.netty.model.payloads.items;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportSize;

/**
 *
 * @author cgarcia
 */
public class AssociatedValueItem {
    private final DataTransportErrorCode returnCode;
    private final DataTransportSize dataTransportSize;
    private final int Length;
    private final ByteBuf Data;   

    public AssociatedValueItem(DataTransportErrorCode returnCode, DataTransportSize dataTransportSize, int Length, ByteBuf Data) {
        this.returnCode = returnCode;
        this.dataTransportSize = dataTransportSize;
        this.Length = Length;
        this.Data = Data;
    }

    public DataTransportErrorCode getReturnCode() {
        return returnCode;
    }

    public DataTransportSize getDataTransportSize() {
        return dataTransportSize;
    }

    public int getLength() {
        return Length;
    }

    public ByteBuf getData() {
        return Data;
    }

    @Override
    public String toString() {
        return "AssociatedValueItem{" + "returnCode=" + returnCode 
                + ", dataTransportSize=" + dataTransportSize 
                + ", Length=" + Length 
                + ", Data=" + ByteBufUtil.hexDump(Data) 
                + '}';
    }

}
