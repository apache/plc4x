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

package org.apache.plc4x.java.tools.ui.model;

import jakarta.persistence.*;

import java.util.Map;

@Entity
public class Device {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Integer id;

    private String name;

    private String protocolCode;
    private String getTransportCode;

    private String getTransportUrl;
    @ElementCollection
    @MapKeyColumn(name="name")
    @Column(name="value")
    @CollectionTable(name="device_options", joinColumns=@JoinColumn(name="device_id"))
    private Map<String, String> getOptions;
    @ElementCollection
    @MapKeyColumn(name="name")
    @Column(name="value")
    @CollectionTable(name="device_attributes", joinColumns=@JoinColumn(name="device_id"))
    private Map<String, String> getAttributes;

    public Device() {
    }

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getProtocolCode() {
        return protocolCode;
    }

    public void setProtocolCode(String protocolCode) {
        this.protocolCode = protocolCode;
    }

    public String getGetTransportCode() {
        return getTransportCode;
    }

    public void setGetTransportCode(String getTransportCode) {
        this.getTransportCode = getTransportCode;
    }

    public String getGetTransportUrl() {
        return getTransportUrl;
    }

    public void setGetTransportUrl(String getTransportUrl) {
        this.getTransportUrl = getTransportUrl;
    }

    public Map<String, String> getGetOptions() {
        return getOptions;
    }

    public void setGetOptions(Map<String, String> getOptions) {
        this.getOptions = getOptions;
    }

    public Map<String, String> getGetAttributes() {
        return getAttributes;
    }

    public void setGetAttributes(Map<String, String> getAttributes) {
        this.getAttributes = getAttributes;
    }

}
