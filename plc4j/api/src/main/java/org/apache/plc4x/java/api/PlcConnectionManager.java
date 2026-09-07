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
package org.apache.plc4x.java.api;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;

/**
 * A {@link PlcConnectionFactory} that keeps the connections it hands out, and therefore has a
 * lifecycle of its own.
 * <p>
 * A connection cache is the typical example: it holds on to the real connections so it can hand
 * them out again, which means someone has to release them eventually. Closing the manager is what
 * releases them, and having that on the interface is what lets a caller dispose of a manager
 * without knowing which implementation it was handed, be it in a try-with-resources block or
 * through a framework's disposal callback.
 * <p>
 * Implementations that merely create connections and leave them to their callers implement
 * {@link PlcConnectionFactory} instead - there is nothing for them to close.
 */
public interface PlcConnectionManager extends PlcConnectionFactory, AutoCloseable {

    /**
     * Closes every connection this manager still holds and releases the resources it acquired,
     * leaving it unusable afterward.
     * <p>
     * Closing a manager that is already closed does nothing.
     *
     * @throws PlcConnectionException if the resources could not be released.
     */
    @Override
    void close() throws PlcConnectionException;

}
