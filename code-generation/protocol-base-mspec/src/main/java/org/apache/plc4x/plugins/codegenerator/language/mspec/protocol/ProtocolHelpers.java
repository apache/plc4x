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
package org.apache.plc4x.plugins.codegenerator.language.mspec.protocol;

import org.apache.plc4x.plugins.codegenerator.protocol.Protocol;
import org.apache.plc4x.plugins.codegenerator.types.exceptions.GenerationException;

import java.io.InputStream;
import java.util.Objects;

public interface ProtocolHelpers extends Protocol {

    /**
     * Returns a mspec stream by using {@link Protocol#getName()}
     *
     * @return the {@link InputStream} of the referenced {@code mspecName}
     * @throws GenerationException if the mspec can't be found.
     */
    default InputStream getMspecStream() throws GenerationException {
        return getMspecStream(getName());
    }

    /**
     * Returns a mspec stream for a give name
     *
     * @param mspecName the name without the .mspec extension
     * @return the {@link InputStream} of the referenced {@code mspecName}
     * @throws GenerationException if the mspec can't be found.
     */
    default InputStream getMspecStream(String mspecName) throws GenerationException {
        Objects.requireNonNull(mspecName, "mspecName must be set");
        String versionSubPath = getVersion().map(version -> "/v" + version).orElse("");
        String packageName = getPackageName();
        String path = "/protocols/" + packageName + versionSubPath + "/" + mspecName + ".mspec";
        InputStream inputStream = getClass().getResourceAsStream(path);
        if (inputStream == null) {
            throw new GenerationException("Error loading " + mspecName + " schema for protocol '" + getName() + "' (path " + path + ")");
        }
        return inputStream;
    }

    /**
     * @return {@link Protocol#getName()} in sanitized form
     */
    default String getSanitizedName() {
        return Objects.requireNonNull(getName(), "protocol should return useful value at getName()")
            // Replace - with emptiness
            .replaceAll("-", "");
    }

    /**
     * @return the package name
     */
    default String getPackageName() {
        return getSanitizedName();
    }
}
