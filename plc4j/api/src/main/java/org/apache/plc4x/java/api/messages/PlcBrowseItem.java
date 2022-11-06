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
package org.apache.plc4x.java.api.messages;

import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;

import java.util.List;
import java.util.Map;

public interface PlcBrowseItem {

    /**
     * @return returns the address of this item
     */
    String getAddress();

    /**
     * @return returns a textual description of this item
     */
    String getName();

    /**
     * @return returns the data-type of this item
     */
    PlcValueType getPlcValueType();

    /**
     * @return returns the array info for this element
     * (this is usually null, but for lists, it contains the array sizes)
     */
    default List<PlcBrowseItemArrayInfo> getArrayInfo() {
        return null;
    }

    /**
     * @return returns 'true' if we can read this variable.
     */
    boolean isReadable();

    /**
     * @return returns 'true' if we can write to this variable.
     */
    boolean isWritable();

    /**
     * @return returns 'true' if we can subscribe this variable.
     */
    boolean isSubscribable();

    /**
     * @return returns any children this item might have
     */
    Map<String, PlcBrowseItem> getChildren();

    /**
     * @return returns a map of additional options the given protocol might provide.
     */
    Map<String, PlcValue> getOptions();

}
