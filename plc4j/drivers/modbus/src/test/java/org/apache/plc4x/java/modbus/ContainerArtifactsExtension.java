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
package org.apache.plc4x.java.modbus;

import org.junit.jupiter.api.extension.*;
import org.testcontainers.containers.GenericContainer;
import org.testcontainers.containers.output.OutputFrame;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.List;

public class ContainerArtifactsExtension implements BeforeEachCallback, AfterEachCallback, ParameterResolver {

    private final List<GenericContainer<?>> containers;
    private Path runDir;

    public ContainerArtifactsExtension(GenericContainer<?>... containers) {
        this.containers = List.of(containers);
    }

    @Override
    public void beforeEach(ExtensionContext context) throws Exception {
        // Create a unique folder per test run
        String testName = context.getDisplayName().replaceAll("[^a-zA-Z0-9._-]", "_");
        String stamp = LocalDateTime.now().format(DateTimeFormatter.ofPattern("yyyyMMdd_HHmmss"));
        runDir = Paths.get("target", "testcontainers", testName + "_" + stamp);
        Files.createDirectories(runDir);
    }

    @Override
    public void afterEach(ExtensionContext context) throws Exception {
        // Dump logs and any extra files you want to keep
        for (GenericContainer<?> c : containers) {
            dumpLogs(c);
            // Example: copy arbitrary files out of the container if they exist
            // safeCopy(c, "/app/logs/app.log", runDir.resolve(c.getDockerImageName().replace("/", "_") + "_app.log"));
            // safeCopy(c, "/var/log/simulator.log", runDir.resolve("simulator.log"));
        }
    }

    private void dumpLogs(GenericContainer<?> c) throws IOException {
        String prettyName = c.getDockerImageName().replace("/", "_");
        String idShort = c.getContainerId() != null ? c.getContainerId().substring(0, 12) : "unknown";
        Path all = runDir.resolve(prettyName + "_" + idShort + ".log");
        Path out = runDir.resolve(prettyName + "_" + idShort + "_stdout.log");
        Path err = runDir.resolve(prettyName + "_" + idShort + "_stderr.log");

        // Combined logs
        Files.writeString(all, c.getLogs(), StandardCharsets.UTF_8);

        // Split logs (optional)
        String stdout = c.getLogs(OutputFrame.OutputType.STDOUT);
        String stderr = c.getLogs(OutputFrame.OutputType.STDERR);
        Files.writeString(out, stdout, StandardCharsets.UTF_8);
        Files.writeString(err, stderr, StandardCharsets.UTF_8);
    }

    /** Copy a file from the container if present; ignore errors. */
    @SuppressWarnings("unused")
    private void safeCopy(GenericContainer<?> c, String containerPath, Path dst) {
        try {
            Files.createDirectories(dst.getParent());
            c.copyFileFromContainer(containerPath, dst.toString());
        } catch (Exception ignored) {
            // File didn't exist, or copy failed—skip silently.
        }
    }

    /** Allow injection via @RegisterExtension with a shared instance if needed. */
    @Override public boolean supportsParameter(ParameterContext pc, ExtensionContext ec) { return false; }
    @Override public Object resolveParameter(ParameterContext pc, ExtensionContext ec) { return null; }
}
