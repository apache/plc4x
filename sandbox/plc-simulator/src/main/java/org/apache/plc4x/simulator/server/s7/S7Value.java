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

package org.apache.plc4x.simulator.server.s7;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

public class S7Value {

    private final S7Type type;
    private final Object value;

    S7Value(S7Type type, Object value) {
        this.type = type;
        this.value = value;
    }

    public S7Type getType() {
        return type;
    }

    public Object getValue() {
        return value;
    }

    public byte[] getData() {
        try {
            WriteBuffer writeBuffer;
            switch (type) {
                case INT:
                    writeBuffer = new WriteBuffer(2, false);
                    writeBuffer.writeShort(16, ((short) value));
                    break;
                case UINT:
                    writeBuffer = new WriteBuffer(2, false);
                    writeBuffer.writeUnsignedInt(16, ((int) value));
                    break;
                case DINT:
                    writeBuffer = new WriteBuffer(4, false);
                    writeBuffer.writeInt(32, ((int) value));
                    break;
                case LINT:
                    writeBuffer = new WriteBuffer(8, false);
                    writeBuffer.writeLong(64, ((long) value));
                    break;
                case REAL:
                    writeBuffer = new WriteBuffer(4, false);
                    writeBuffer.writeFloat(32, ((float) value));
                    break;
                default:
                    throw new NotImplementedException("Currently the type " + type + " is not implemented!");
            }
            return writeBuffer.getData();
        } catch (ParseException e) {
            throw new IllegalStateException("This should not happen!", e);
        }
    }

    @Override
    public String toString() {
        return "S7Value{" +
            "type=" + type +
            ", value=" + value +
            '}';
    }

}
