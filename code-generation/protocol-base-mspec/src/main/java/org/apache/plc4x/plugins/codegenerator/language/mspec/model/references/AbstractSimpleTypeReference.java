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
package org.apache.plc4x.plugins.codegenerator.language.mspec.model.references;

import org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference;

import java.util.Objects;

public abstract class AbstractSimpleTypeReference implements SimpleTypeReference {

    protected final SimpleBaseType baseType;
    protected final int sizeInBits;

    public AbstractSimpleTypeReference(SimpleBaseType baseType, int sizeInBits) {
        this.baseType = Objects.requireNonNull(baseType);
        this.sizeInBits = sizeInBits;
    }

    @Override
    public SimpleBaseType getBaseType() {
        return baseType;
    }

    @Override
    public int getSizeInBits() {
        return sizeInBits;
    }

    @Override
    public String toString() {
        return "AbstractSimpleTypeReference{" +
                "baseType=" + baseType +
                ", sizeInBits=" + sizeInBits +
                '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof AbstractSimpleTypeReference)) return false;
        AbstractSimpleTypeReference that = (AbstractSimpleTypeReference) o;
        return getSizeInBits() == that.getSizeInBits() && getBaseType() == that.getBaseType();
    }

    @Override
    public int hashCode() {
        return Objects.hash(getBaseType(), getSizeInBits());
    }
}
