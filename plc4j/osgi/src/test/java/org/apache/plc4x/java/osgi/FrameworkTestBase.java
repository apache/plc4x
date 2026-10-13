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
package org.apache.plc4x.java.osgi;

import org.apache.felix.framework.FrameworkFactory;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.io.TempDir;
import org.osgi.framework.BundleContext;
import org.osgi.framework.Constants;
import org.osgi.framework.launch.Framework;

import java.nio.file.Path;
import java.util.Map;

/**
 * Runs each test against a real, embedded Felix framework, the one Karaf 4.4 is built on. The activators get the
 * system bundle's context, whose class loader is the test class path, so they find the test META-INF/services.
 */
abstract class FrameworkTestBase {

    @TempDir
    Path storage;

    private Framework framework;
    protected BundleContext context;

    @BeforeEach
    void startFramework() throws Exception {
        framework = new FrameworkFactory().newFramework(Map.of(
            Constants.FRAMEWORK_STORAGE, storage.toString(),
            Constants.FRAMEWORK_STORAGE_CLEAN, Constants.FRAMEWORK_STORAGE_CLEAN_ONFIRSTINIT));
        framework.start();
        context = framework.getBundleContext();
    }

    @AfterEach
    void stopFramework() throws Exception {
        framework.stop();
        framework.waitForStop(10_000);
    }

}
