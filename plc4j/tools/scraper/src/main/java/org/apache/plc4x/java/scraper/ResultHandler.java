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
package org.apache.plc4x.java.scraper;

import java.util.Map;

/**
 * Callback interface to handle results of one run of a {@link ScraperTask}.
 */
@FunctionalInterface
public interface ResultHandler {

    /**
     * Callback handler.
     * @param job name of the job (from config)
     * @param alias alias of the connection (<b>not</b> connection String)
     * @param results Results in the form alias to result value
     */
    void handle(String job, String alias, Map<String, Object> results);

}
