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

package org.apache.plc4x.java.scraper.config;

import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class ScraperConfigurationBuilderTest {

    @Test
    void checkSyntax() {
        ScraperConfigurationBuilder builder = new ScraperConfigurationBuilder();
        List<String> sources = Arrays.asList("s1", "s2");
        List<String> jobs = Arrays.asList("j1", "j2");

        sources.forEach(source -> builder.addSource(source, source));
        for (String job : jobs) {
            JobConfigurationBuilder jobConfigurationBuilder = builder.job(job, 10);
            for (int i = 1; i <= 100; i++) {
                jobConfigurationBuilder.field("f" + i, "qry" + i);
            }
            jobConfigurationBuilder.build();
        }

        ScraperConfiguration configuration = builder.build();

        // TODO add assert.
        System.out.println(configuration);
    }
}