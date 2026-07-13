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
package org.apache.plc4x.java.opcua.tag;

import org.apache.plc4x.java.api.model.PlcQuery;

/**
 * A browse query for the OPC UA driver. The query string is the address of the node to
 * start browsing from (e.g. {@code ns=2;s=HelloWorld}); an empty query starts at the
 * standard {@code Objects} folder.
 */
public class OpcuaQuery implements PlcQuery {

    private final String queryString;

    public OpcuaQuery(String queryString) {
        this.queryString = queryString;
    }

    @Override
    public String getQueryString() {
        return queryString;
    }

    @Override
    public String toString() {
        return "OpcuaQuery{" + queryString + '}';
    }

}
