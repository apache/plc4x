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

package org.apache.plc4x.java.opcua;

import java.nio.file.Path;
import java.nio.file.Paths;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.testcontainers.containers.GenericContainer;
import org.testcontainers.containers.wait.strategy.Wait;
import org.testcontainers.images.builder.ImageFromDockerfile;

public class MiloTestContainer extends GenericContainer<MiloTestContainer> {

    private final static Logger logger = LoggerFactory.getLogger(MiloTestContainer.class);

    private final static ImageFromDockerfile IMAGE = inlineImage();

    public MiloTestContainer() {
        super(IMAGE);

        waitingFor(Wait.forLogMessage("Server started\\s*", 1))
            // Uncomment below to debug Milo server
            //.withStartupTimeout(Duration.ofMinutes(10))
        ;
        // Fixed 12686 -> 12686 mapping. The Milo server advertises its endpoints on
        // localhost:12686 (its internal bind port), so the host port MUST equal 12686 —
        // otherwise the client connects to a random mapped port, is redirected to the
        // (unreachable) advertised localhost:12686 endpoint, and the handshake hangs.
        addFixedExposedPort(12686, 12686);

        // Uncomment below to enable server debug
        //withEnv("JAVA_TOOL_OPTIONS", "-agentlib:jdwp=transport=dt_socket,address=*:8000,server=y,suspend=y");
        //addFixedExposedPort(8000, 8000);
    }

    private static ImageFromDockerfile inlineImage() {
        Path absolutePath = Paths.get(".").toAbsolutePath();
        logger.info("Building milo server image from {}", absolutePath);
        // Reuse the Docker layer cache across runs; the build context only changes when
        // the milo jar or the compiled TestMiloServer classes change, so a cached build
        // stays correct while turning a multi-minute rebuild into a near-instant one.
        return new ImageFromDockerfile("plc4x-milo-test", false)
            .withDockerfile(absolutePath.resolve("Dockerfile.test"));
    }

}
