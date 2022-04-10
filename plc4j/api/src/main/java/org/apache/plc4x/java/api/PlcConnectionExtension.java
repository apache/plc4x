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
package org.apache.plc4x.java.api;

import org.apache.commons.lang3.NotImplementedException;

import java.util.concurrent.Future;

/**
 * Suggestion for new API
 */
@Experimental
public interface PlcConnectionExtension {

    /**
     * <code>
     *     "%DB400.xxx"
     * </code>
     * <code>
     *     ["%DB400.xxx", ""]
     * </code>
     * <code>
     *     { "item1": "%DB400.xxx", "item2": "xx" }
     * </code>
     * Prepared Statement
     * <code>
     *     { "item1": ?, "item2": ? }
     * </code>
     *
     * @param s
     * @return
     */
    @Experimental
    default Future<NewPlcResponse> query(String s, Object... args) {
        throw new NotImplementedException("Not implemented");
    }

    /**
     * <code>
     *     define { "item1": "%DB400.xxx", "item2": "xx" } AS my_pymél_struct
     * </code>
     * <code>
     *     define %DN4ßß" AS "my_structure"
     * </code>
     * @param pql
     * @return
     */
    @Experimental
    default Future<Boolean> execute(String pql) {
        throw new NotImplementedException("Not implemented");
    }

    /**
     * Planned successor for
     * {@link org.apache.plc4x.java.api.messages.PlcResponse}
     * {@link org.apache.plc4x.java.api.messages.PlcReadResponse}
     * {@link org.apache.plc4x.java.api.messages.PlcWriteResponse}
     * {@link org.apache.plc4x.java.api.messages.PlcSubscriptionResponse}
     * {@link org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse}
     */
    @Experimental
    interface NewPlcResponse {

    }
}
