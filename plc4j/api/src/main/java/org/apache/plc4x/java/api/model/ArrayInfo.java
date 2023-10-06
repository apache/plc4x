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

public interface ArrayInfo {

    /**
     * @return Number of elements in total
     */
    int getSize();

    /**
     * As in PLCs not every array starts at 0, we need to be flexible with this.
     * In the default usage scenario of a simple array [6] this index will be 0 by default.
     * @return Returns the index of lower bound of the array.
     */
    int getLowerBound();

    /**
     * As in PLCs not every array starts at 0, we need to be flexible with this.
     * In the default usage scenario of a simple array [6] this index will be match the array size.
     * @return Returns the index of upper bound of the array.
     */
    int getUpperBound();

}
