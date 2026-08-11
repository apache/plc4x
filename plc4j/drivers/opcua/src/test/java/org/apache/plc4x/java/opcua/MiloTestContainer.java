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

import org.testcontainers.containers.GenericContainer;
import org.testcontainers.containers.wait.strategy.Wait;
import org.testcontainers.images.builder.ImageFromDockerfile;

public class MiloTestContainer extends GenericContainer<MiloTestContainer> {

    private final static String SERVER_RESOURCES = "opcua/server";

    private final static String[] SERVER_SOURCES = {
        "org/eclipse/milo/examples/server/EventNotifierTask.java",
        "org/eclipse/milo/examples/server/Plc4xTestNamespace.java",
        "org/eclipse/milo/examples/server/Plc4xTestStruct.java",
        "org/eclipse/milo/examples/server/TestMiloServer.java"
    };

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
        // The build context is assembled from the test resources: the Dockerfile plus the
        // sources of the test server, which are compiled inside the image against the Milo
        // uber jar it downloads (see the Dockerfile for why Milo is not a Maven dependency).
        //
        // Keeping the named image around reuses the Docker layer cache across runs; the
        // context only changes when the Dockerfile or those sources change, so a cached
        // build stays correct while turning a multi-minute rebuild into a near-instant one.
        ImageFromDockerfile image = new ImageFromDockerfile("plc4x-milo-test", false)
            .withFileFromClasspath("Dockerfile", SERVER_RESOURCES + "/Dockerfile");
        for (String source : SERVER_SOURCES) {
            image.withFileFromClasspath("src/" + source, SERVER_RESOURCES + "/src/" + source);
        }
        return image;
    }

}
