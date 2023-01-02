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
package org.apache.plc4x.java.s7.readwrite.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.model.PlcQuery;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;

public class S7PlcTagHandler implements PlcTagHandler {

    @Override
    public PlcTag parseTag(String tagAddress) {
        if (S7Tag.matches(tagAddress)) {
            return S7Tag.of(tagAddress);
        } else if (S7SubscriptionTag.matches(tagAddress)){
            return S7SubscriptionTag.of(tagAddress);
        } else if (S7SslTag.matches(tagAddress)){
            return S7SslTag.of(tagAddress);
        }
        throw new PlcInvalidTagException(tagAddress);
    }

    @Override
    public PlcQuery parseQuery(String query) {
        throw new UnsupportedOperationException("This driver doesn't support browsing");
    }

}
