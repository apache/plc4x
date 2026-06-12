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
package org.apache.plc4x.java.knxnetip.tag;

import org.apache.plc4x.java.api.model.PlcQuery;

/**
 * Trivial {@link PlcQuery} that just carries the user's address pattern string
 * straight through to the connection. KNX query patterns are simple textual
 * group-address globs ({@code *}, {@code 1/*}, {@code 1/2/*}, {@code 1/2/3});
 * compiling them to a regex happens lazily inside the connection's browse path.
 */
public class KnxNetIpQuery implements PlcQuery {

    private final String queryString;

    public KnxNetIpQuery(String queryString) {
        this.queryString = queryString;
    }

    @Override
    public String getQueryString() {
        return queryString;
    }

}
