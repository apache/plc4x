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
package org.apache.plc4x.java.utils.connectionpool;

import org.apache.plc4x.java.api.authentication.PlcAuthentication;

import java.util.Objects;

public abstract class PoolKey {
    protected final String url;
    protected final PlcAuthentication plcAuthentication;

    public PoolKey(String url, PlcAuthentication plcAuthentication) {
        this.url = url;
        this.plcAuthentication = plcAuthentication;
    }

    public String getUrl() {
        return url;
    }

    public PlcAuthentication getPlcAuthentication() {
        return plcAuthentication;
    }

    /**
     * @return the part of the url that should be pooled.
     */
    public abstract String getPoolableKey();

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof PoolKey)) {
            return false;
        }
        PoolKey poolKey = (PoolKey) o;
        return Objects.equals(getPoolableKey(), poolKey.getPoolableKey()) &&
            Objects.equals(plcAuthentication, poolKey.plcAuthentication);
    }

    @Override
    public int hashCode() {
        return Objects.hash(getPoolableKey(), plcAuthentication);
    }

    @Override
    public String toString() {
        return "PoolKey{" +
            "url='" + url + '\'' +
            (plcAuthentication != PooledPlcDriverManager.noPlcAuthentication ? ", plcAuthentication=" + plcAuthentication : "") +
            '}';
    }
}