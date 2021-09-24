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
package org.apache.plc4x.java.spi.generation.io;

import com.github.jinahya.bit.io.ArrayByteInput;
import com.github.jinahya.bit.io.DefaultBitInput;

import java.io.IOException;

/**
 * Modified version that exposes the position.
 */
public class MyDefaultBitInput extends DefaultBitInput<ArrayByteInput> {

    public MyDefaultBitInput(ArrayByteInput delegate) {
        super(delegate);
    }

    public long getPos() {
        // TODO: we should report bits as we would loose this information on a reset
        return delegate.getIndex();
    }

    public void reset(int pos) {
        try {
            align(1);
        } catch (IOException e) {
            e.printStackTrace();
        }
        delegate.setIndex(pos);
    }
}
