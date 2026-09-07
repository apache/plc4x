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
package org.apache.plc4x.java.knxnetip.ets.filehandlers;

import net.lingala.zip4j.ZipFile;
import org.bouncycastle.crypto.digests.SHA256Digest;
import org.bouncycastle.crypto.generators.PKCS5S2ParametersGenerator;
import org.bouncycastle.crypto.params.KeyParameter;

import java.io.File;
import java.nio.charset.StandardCharsets;
import java.util.Base64;

/**
 * Password handling for ETS6.x project files (schema 21+). Derives the ZIP
 * password via PBKDF2-SHA256 with a fixed salt — schema versions 21, 22, 23,
 * ... all use the same derivation; the schema number changing does not
 * change the salt.
 */
public class Ets6FileHandler implements EtsFileHandler {

    private static final String SALT = "21.project.ets.knx.org";
    private static final int ITERATIONS = 65536;
    private static final int KEY_LENGTH_BITS = 256;

    @Override
    public ZipFile getProjectFiles(File archive, String password) {
        return new ZipFile(archive, getProcessedPassword(password).toCharArray());
    }

    protected String getProcessedPassword(String originalPassword) {
        PKCS5S2ParametersGenerator gen = new PKCS5S2ParametersGenerator(new SHA256Digest());
        gen.init(originalPassword.getBytes(StandardCharsets.UTF_16LE),
            SALT.getBytes(StandardCharsets.UTF_8), ITERATIONS);
        byte[] hashedPassword = ((KeyParameter) gen.generateDerivedParameters(KEY_LENGTH_BITS)).getKey();
        return Base64.getEncoder().encodeToString(hashedPassword);
    }

}
