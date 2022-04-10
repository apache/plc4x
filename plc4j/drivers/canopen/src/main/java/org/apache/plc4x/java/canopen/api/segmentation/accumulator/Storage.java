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
package org.apache.plc4x.java.canopen.api.segmentation.accumulator;

/**
 * A storage which is called for each received segment.
 *
 * @param <T> Type of frame.
 * @param <R> Type of result.
 */
public interface Storage<T, R> {

    /**
     * Appends segmented frame.
     *
     * @param frame Segmented frame.
     */
    void append(T frame);

    /**
     * Gets accumulated size of stored data.
     *
     * @return Occupied memory in bytes.
     */
    long size();

    /**
     * Retrieves final result from segmented payload.
     *
     * @return Assembled result.
     */
    R get();

}
