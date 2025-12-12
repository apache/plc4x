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
package org.apache.plc4x.java.spi.buffers.bytebased.encoding;

import org.apache.plc4x.java.spi.buffers.api.WithOption;

import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;

public class EncodingUTF8 extends BaseStringEncoding {

    public static final String NAME = "UTF8";

    public static WithOption optionEncodingUTF8() {
        return WithOption.WithEncoding(NAME);
    }

    @Override
    public String getName() {
        return NAME;
    }

    @Override
    protected int getBitsPerCharacter() {
        return 8;
    }

    @Override
    protected Charset getCharset() {
        return StandardCharsets.UTF_8;
    }

}
