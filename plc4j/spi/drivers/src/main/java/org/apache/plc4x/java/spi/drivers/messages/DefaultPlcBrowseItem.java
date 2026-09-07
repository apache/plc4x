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
package org.apache.plc4x.java.spi.drivers.messages;

import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.value.PlcValue;

import java.util.Collections;
import java.util.EnumSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

public class DefaultPlcBrowseItem implements PlcBrowseItem {

    private final PlcTag tag;
    private final String name;
    private final boolean readable;
    private final boolean writable;
    private final Set<PlcSubscriptionType> supportedSubscriptionTypes;
    private final boolean publishable;
    private final List<ArrayInfo> arrayInformation;
    private final Map<String, PlcBrowseItem> children;
    private final Map<String, PlcValue> options;

    public DefaultPlcBrowseItem(PlcTag tag, String name, boolean readable, boolean writable,
                                Set<PlcSubscriptionType> supportedSubscriptionTypes, boolean publishable,
                                List<ArrayInfo> arrayInformation,
                                Map<String, PlcBrowseItem> children,
                                Map<String, PlcValue> options) {
        this.tag = tag;
        this.name = name;
        this.readable = readable;
        this.writable = writable;
        this.supportedSubscriptionTypes = (supportedSubscriptionTypes == null || supportedSubscriptionTypes.isEmpty())
            ? Collections.emptySet()
            : Collections.unmodifiableSet(EnumSet.copyOf(supportedSubscriptionTypes));
        this.publishable = publishable;
        this.arrayInformation = arrayInformation;
        this.children = children;
        this.options = options;
    }

    @Override
    public PlcTag getTag() {
        return tag;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public boolean isReadable() {
        return readable;
    }

    @Override
    public boolean isWritable() {
        return writable;
    }

    @Override
    public Set<PlcSubscriptionType> getSupportedSubscriptionTypes() {
        return supportedSubscriptionTypes;
    }

    @Override
    public boolean isPublishable() {
        return publishable;
    }

    @Override
    public boolean isArray() {
        return arrayInformation != null && !arrayInformation.isEmpty();
    }

    @Override
    public List<ArrayInfo> getArrayInformation() {
        return arrayInformation;
    }

    @Override
    public Map<String, PlcBrowseItem> getChildren() {
        return children;
    }

    @Override
    public Map<String, PlcValue> getOptions() {
        return options;
    }

}
