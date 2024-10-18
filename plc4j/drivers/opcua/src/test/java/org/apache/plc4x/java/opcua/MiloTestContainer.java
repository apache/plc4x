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

        waitingFor(Wait.forLogMessage("Server started\\s*", 1));
        addExposedPort(12686);
    }

    private static ImageFromDockerfile inlineImage() {
        Path absolutePath = Paths.get(".").toAbsolutePath();
        logger.info("Building milo server image from {}", absolutePath);
        return new ImageFromDockerfile("plc4x-milo-test", false)
            .withBuildImageCmdModifier(cmd -> cmd.withNoCache(true))
            .withDockerfile(absolutePath.resolve("Dockerfile.test"));
    }

}
