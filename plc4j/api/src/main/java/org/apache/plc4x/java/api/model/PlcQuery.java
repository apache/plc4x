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
package org.apache.plc4x.java.api.model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonTypeInfo;

/**
 * Base type for all query types.
 * Typically, every driver provides an implementation of this interface in order
 * to be able to describe the queries of a browse request. As this is completely tied to
 * the implemented protocol, this base interface makes absolutely no assumption to
 * any information it should provide.
 *
 * In order to stay platform and protocol independent every driver connection implementation
 * provides a prepareQuery(String) method that is able to parse a string representation of
 * a resource into it's individual query type. Manually constructing PlcQuery objects
 * manually makes the solution less independent from the protocol, but might be faster.
 */
@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public interface PlcQuery {

    @JsonIgnore
    String getQueryString();

}
