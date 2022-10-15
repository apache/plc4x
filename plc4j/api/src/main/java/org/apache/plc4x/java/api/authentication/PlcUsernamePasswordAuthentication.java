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
package org.apache.plc4x.java.api.authentication;

import java.util.Objects;

public class PlcUsernamePasswordAuthentication implements PlcAuthentication {

    private final String username;
    private final String password;

    public PlcUsernamePasswordAuthentication(String username, String password) {
        Objects.requireNonNull(username, "User name must not be null");
        Objects.requireNonNull(password, "Password must not be null");
        this.username = username;
        this.password = password;
    }

    public String getUsername() {
        return username;
    }

    public String getPassword() {
        return password;
    }

    @Override
    public final boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof PlcUsernamePasswordAuthentication)) {
            return false;
        }
        PlcUsernamePasswordAuthentication that = (PlcUsernamePasswordAuthentication) o;
        return Objects.equals(username, that.username) &&
            Objects.equals(password, that.password);
    }

    @Override
    public final int hashCode() {
        return Objects.hash(username, password);
    }

    @Override
    public String toString() {
        return "PlcUsernamePasswordAuthentication{" +
            "username='" + username + '\'' +
            ", password='" + "*****************" + '\'' +
            '}';
    }

}
