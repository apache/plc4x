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
package org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms;

import org.apache.plc4x.plugins.codegenerator.types.terms.HexadecimalLiteral;

import java.util.Objects;

public class DefaultHexadecimalLiteral implements HexadecimalLiteral {

    private final String hexString;

    public DefaultHexadecimalLiteral(String hexString) {
        this.hexString = hexString;
    }

    public String getHexString() {
        return hexString;
    }

    @Override
    public String stringRepresentation() {
        return hexString;
    }

    @Override
    public String toString() {
        return "DefaultHexadecimalLiteral{" +
                "hexString=" + hexString +
                '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DefaultHexadecimalLiteral that = (DefaultHexadecimalLiteral) o;
        return hexString.equals(that.hexString);
    }

    @Override
    public int hashCode() {
        return Objects.hash(hexString);
    }
}
