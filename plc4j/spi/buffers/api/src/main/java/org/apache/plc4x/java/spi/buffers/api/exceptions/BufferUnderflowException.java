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

package org.apache.plc4x.java.spi.buffers.api.exceptions;

/**
 * A read asked for more data than the buffer has left.
 * <p>
 * This is distinct from the other buffer failures because it says nothing is wrong with the data
 * that <em>is</em> there - only that it ended. An optional field can therefore treat it as the field
 * being absent and carry on, which is what makes a best-effort read of a message that may end early
 * possible. Every other failure means the bytes present could not be made sense of, and must not be
 * mistaken for an absent field.
 * <p>
 * It deliberately does not cover a read whose requested width is implausible or negative: that is a
 * malformed length rather than a short message, and callers size allocations from those widths.
 */
public class BufferUnderflowException extends BufferException {

    public BufferUnderflowException(String message) {
        super(message);
    }

}
