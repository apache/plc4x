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
package org.apache.plc4x.java.spi.generation.io;

import com.github.jinahya.bit.io.ArrayByteInput;
import com.github.jinahya.bit.io.DefaultBitInput;
import org.apache.plc4x.java.spi.utils.hex.Hex;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;

/**
 * Modified version that exposes the position.
 */
public class MyDefaultBitInput extends DefaultBitInput<ArrayByteInput> {

    private static final Logger LOGGER = LoggerFactory.getLogger(MyDefaultBitInput.class);

    public MyDefaultBitInput(ArrayByteInput delegate) {
        super(delegate);
    }

    public long getPos() {
        // TODO: we should report bits as we would loose this information on a reset
        return delegate.getIndex();
    }

    public void reset(int pos) {
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("Before\n{}", Hex.dump(delegate.getSource(), Hex.DefaultWidth, delegate.getIndex()));
        }
        try {
            long align = align(1);
            LOGGER.debug("aligned {}", align);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        delegate.setIndex(pos);
        LOGGER.debug("set to index {}", pos);
        if (LOGGER.isDebugEnabled()) {
            LOGGER.debug("After\n{}", Hex.dump(delegate.getSource(), Hex.DefaultWidth, pos));
        }
    }
}
