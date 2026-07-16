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

package org.apache.plc4x.java.tools.eventpump.config;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * Configuration for a PLC connection.
 */
public class ConnectionConfiguration {

    @JsonProperty("id")
    private String id;

    @JsonProperty("url")
    private String url;

    @JsonProperty("username")
    private String username;

    @JsonProperty("password")
    private String password;

    /**
     * Get the connection ID.
     *
     * @return The connection ID
     */
    public String getId() {
        return id;
    }

    /**
     * Set the connection ID.
     *
     * @param id The connection ID
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Get the connection URL.
     *
     * @return The connection URL
     */
    public String getUrl() {
        return url;
    }

    /**
     * Set the connection URL.
     *
     * @param url The connection URL
     */
    public void setUrl(String url) {
        this.url = url;
    }

    /**
     * Get the username for authentication.
     *
     * @return The username, or null if not using authentication
     */
    public String getUsername() {
        return username;
    }

    /**
     * Set the username for authentication.
     *
     * @param username The username
     */
    public void setUsername(String username) {
        this.username = username;
    }

    /**
     * Get the password for authentication.
     *
     * @return The password, or null if not using authentication
     */
    public String getPassword() {
        return password;
    }

    /**
     * Set the password for authentication.
     *
     * @param password The password
     */
    public void setPassword(String password) {
        this.password = password;
    }
}
