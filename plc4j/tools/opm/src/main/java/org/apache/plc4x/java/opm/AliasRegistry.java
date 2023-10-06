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
package org.apache.plc4x.java.opm;

/**
 * This interface can be used to give "aliases" for tag names in {@link PlcEntity}s {@link PlcTag} strings.
 * These are then resolved.
 */
public interface AliasRegistry {

    /**
     * Checks if this registry can resolve this alias
     */
    boolean canResolve(String alias);

    /**
     * Checks if this registry can resolve this alias for the given connection.
     */
    boolean canResolve(String connection, String alias);

    /**
     * Resolves an alias to a valid PLC Tag Address
     */
    String resolve(String alias);

    /**
     * Resolves an alias to a valid PLC Field based on the connection.
     * This means that the same alias could be resolved to different Addresses for different connections.
     */
    String resolve(String connection, String alias);

}
