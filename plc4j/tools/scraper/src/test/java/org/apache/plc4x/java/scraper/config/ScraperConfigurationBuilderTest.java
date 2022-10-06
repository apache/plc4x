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
package org.apache.plc4x.java.scraper.config;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import org.assertj.core.api.WithAssertions;

import java.util.Arrays;
import java.util.List;

class ScraperConfigurationBuilderTest implements WithAssertions {

    //ToDo test is failing idon't know why (Tim)
    void builder_usage_example() throws JsonProcessingException {
        ScraperConfigurationClassicImplBuilder builder = new ScraperConfigurationClassicImplBuilder();
        List<String> sources = Arrays.asList("s1", "s2");
        List<String> jobs = Arrays.asList("j1", "j2");

        sources.forEach(source -> builder.addSource(source, source));
        for (String job : jobs) {
            JobConfigurationClassicImplBuilder jobConfigurationClassicImplBuilder = builder.job(job, 10);
            sources.forEach(jobConfigurationClassicImplBuilder::source);
            for (int i = 1; i <= 10; i++) {
                jobConfigurationClassicImplBuilder.field("f" + i, "qry" + i);
            }
            jobConfigurationClassicImplBuilder.build();
        }

        ScraperConfiguration configuration = builder.build();

        ObjectMapper mapper = new ObjectMapper(new YAMLFactory());
        String s = mapper.writeValueAsString(configuration);

        assertThat(configuration.getJobConfigurations())
            .hasSize(2);
        assertThat(configuration.getSources())
            .hasSize(2);
        assertThat(s).contains("sources:\n" +
                                "  s1: \"s1\"\n" +
                                "  s2: \"s2\"\n" +
                                "jobs:\n" +
                                "- name: \"j1\"\n" +
                                "  scrapeRate: 10\n" +
                                "  connections:\n" +
                                "    s1: \"s1\"\n" +
                                "    s2: \"s2\"\n" +
                                "  fields:");

    }
}