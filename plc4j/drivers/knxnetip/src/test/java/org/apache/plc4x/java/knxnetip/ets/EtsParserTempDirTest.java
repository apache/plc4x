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
package org.apache.plc4x.java.knxnetip.ets;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.HashSet;
import java.util.Set;
import java.util.stream.Stream;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

/**
 * Reading a knxproj unpacks part of it beside the process, so the directory it unpacks into has to
 * go again afterwards - including when the read fails, which is the case an operator pointing the
 * driver at the wrong file hits repeatedly.
 */
class EtsParserTempDirTest {

    private static Set<Path> tempDirectories() throws IOException {
        Path tmp = Path.of(System.getProperty("java.io.tmpdir"));
        try (Stream<Path> entries = Files.list(tmp)) {
            Set<Path> dirs = new HashSet<>();
            entries.filter(Files::isDirectory).forEach(dirs::add);
            return dirs;
        }
    }

    @Test
    void aFailedReadLeavesNothingBehind(@TempDir Path workDir) throws IOException {
        File notAKnxproj = workDir.resolve("broken.knxproj").toFile();
        Files.write(notAKnxproj.toPath(), "this is not a zip archive".getBytes());

        Set<Path> before = tempDirectories();
        // Reading it has to fail - the point is what it leaves behind when it does.
        assertThrows(RuntimeException.class, () -> new EtsParser().parse(notAKnxproj, null));
        Set<Path> after = tempDirectories();

        after.removeAll(before);
        assertEquals(Set.of(), after,
            "a failed read must not leave a directory behind in the temp directory");
    }

    @Test
    void repeatedFailedReadsDoNotAccumulate(@TempDir Path workDir) throws IOException {
        File notAKnxproj = workDir.resolve("broken.knxproj").toFile();
        Files.write(notAKnxproj.toPath(), "this is not a zip archive".getBytes());

        Set<Path> before = tempDirectories();
        for (int i = 0; i < 20; i++) {
            assertThrows(RuntimeException.class, () -> new EtsParser().parse(notAKnxproj, null));
        }
        Set<Path> after = tempDirectories();

        after.removeAll(before);
        assertEquals(Set.of(), after,
            "twenty failed reads must not leave twenty directories behind");
    }
}
