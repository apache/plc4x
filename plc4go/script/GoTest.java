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

package org.apache.plc4x.go;

import org.junit.jupiter.api.DynamicContainer;
import org.junit.jupiter.api.TestFactory;

import java.io.BufferedReader;
import java.io.File;
import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.file.Path;
import java.util.Map;
import java.util.stream.Stream;

import static java.util.Locale.ROOT;
import static java.util.Objects.requireNonNull;
import static java.util.stream.Collectors.groupingBy;
import static java.util.stream.Collectors.toList;
import static org.junit.jupiter.api.Assertions.fail;
import static org.junit.jupiter.api.DynamicContainer.dynamicContainer;
import static org.junit.jupiter.api.DynamicTest.dynamicTest;

class GoTest {
    @TestFactory
    Stream<DynamicContainer> goTests() throws IOException, InterruptedException {
        final var goVersion = System.getProperty("go.version", "1.20.7"); // propagated from pom ${go.version}/surefire
        final var osName = System.getProperty("os.name", "").toLowerCase(ROOT);
        final var os = osName.contains("linux") ?
            "linux" :
            // todo: check these two
            (osName.contains("win") ? "win" : "darwin");
        final var arch = System.getProperty("os.arch");
        final var goBase = Path.of(System.getProperty("user.home"))
            .resolve(".mvnGoLang")
            .resolve("go" + goVersion + "." + os + "-" + arch);
        final var go = goBase.resolve("bin/go");
        final var base = requireNonNull(go.getParent().getParent());
        final var target = Path.of("target").toAbsolutePath();
        final var env = Map.of(
            "GOROOT", base.toString(),
            "GOPATH",
            base.getParent().toAbsolutePath().resolve(".go_path") + File.pathSeparator +
                Path.of(".").toAbsolutePath(),
            "GOBIN", target.toString(),
            "GOCACHE", target.resolve("goBuildCache").toString(),
            "PATH", go.getParent().toAbsolutePath() + File.pathSeparator + System.getenv("PATH"));

        return readTestsFromOutput(env, go)
            .collect(groupingBy(it -> {
                final int sep = it.indexOf('_');
                if (sep > 0) {
                    return it.substring(0, sep);
                }
                return it;
            }))
            .entrySet().stream()
            .map(group -> dynamicContainer(group.getKey(), group.getValue().stream()
                .map(test -> dynamicTest(test.substring(test.indexOf('_') + 1), () -> runGoTest(env, go, test)))));
    }

    private void runGoTest(final Map<String, String> env, final Path go, final String test) throws IOException, InterruptedException {
        final var result = go(env, go, "test", "-v", "-run", test, "./...");
        if (result.exitValue() != 0) {
            fail("status=" + result.exitValue() + "\n" + readOutput(result));
        } else { // if (logOutput)
            try {
                try (final var stdout = result.getInputStream()) {
                    System.out.println(new String(stdout.readAllBytes()));
                }
            } catch (final RuntimeException | IOException re) {
                // ignore
            }
            try {
                try (final var stderr = result.getErrorStream()) {
                    System.err.println(new String(stderr.readAllBytes()));
                }
            } catch (final RuntimeException | IOException re) {
                // ignore
            }
        }
    }

    private Stream<String> readTestsFromOutput(final Map<String, String> env, final Path go) throws IOException, InterruptedException {
        final var list = go(env, go, "test", "-list", ".", "./...");
        if (list.exitValue() != 0) {
            fail("Can't list go tests: status=" + list.exitValue() + "\n" + readOutput(list));
        }
        try (final var reader = new BufferedReader(new InputStreamReader(list.getInputStream()))) {
            final var allTests = reader.lines()
                .filter(it -> !it.isBlank() && !it.contains(" ") && !it.contains("."))
                .collect(toList()); // meterialize before closing the reader
            return allTests.stream();
        }
    }

    private String readOutput(final Process list) {
        final var out = new StringBuilder();
        try {
            try (final var stdout = list.getInputStream()) {
                out.append(new String(stdout.readAllBytes()));
            }
            out.append("\n");
        } catch (final RuntimeException | IOException re) {
            // ignore
        }
        try {
            try (final var stderr = list.getErrorStream()) {
                out.append(new String(stderr.readAllBytes()));
            }
        } catch (final RuntimeException | IOException re) {
            // ignore
        }
        return out.toString();
    }

    private Process go(final Map<String, String> env, final Path go, final String... args) throws IOException, InterruptedException {
        final var processBuilder = new ProcessBuilder(Stream.concat(Stream.of(go.toString()), Stream.of(args)).toArray(String[]::new));
        processBuilder.environment().putAll(env);
        final var process = processBuilder.start();
        process.waitFor();
        return process;
    }
}