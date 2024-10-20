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

import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.value.PlcValue;

import java.util.List;
import java.util.Map;

public interface PlcBrowseItem {

    /**
     * @return returns the tag
     */
    PlcTag getTag();

    /**
     * @return returns a textual description of this item
     */
    String getName();

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
     * @return returns 'true' if we can publish this variable.
     */
    boolean isPublishable();

    boolean isArray();

    /**
     * @return list of elements providing information about the array dimensions of this item.
     */
    List<ArrayInfo> getArrayInformation();

    /**
     * @return returns any children this item might have
     */
    Map<String, PlcBrowseItem> getChildren();

    /**
     * @return returns a map of additional options the given protocol might provide.
     */
    Map<String, PlcValue> getOptions();

    /**
     * Sometimes it would be beneficial for clients to have the array information resolved to
     * individual elements, allowing to treat array items as children. As not all drivers form addresses
     * the same way, this option allows the driver to override the structure.
     *
     * @return if a browse item has array information, resolve this information
     *  translating the array elements to child elements.
     */
    //PlcBrowseItem resolveArrayItems();

}
