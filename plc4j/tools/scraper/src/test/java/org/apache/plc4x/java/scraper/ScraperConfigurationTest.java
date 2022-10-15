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
package org.apache.plc4x.java.scraper;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.exc.MismatchedInputException;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.scraper.config.JobConfiguration;
import org.apache.plc4x.java.scraper.config.JobConfigurationClassicImpl;
import org.apache.plc4x.java.scraper.config.ScraperConfiguration;
import org.apache.plc4x.java.scraper.config.ScraperConfigurationClassicImpl;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.junit.jupiter.MockitoExtension;

import java.io.IOException;
import java.util.List;

@ExtendWith(MockitoExtension.class)
class ScraperConfigurationTest implements WithAssertions {

    ObjectMapper mapper = new ObjectMapper(new YAMLFactory());

    @Test
    void parseJobs_fromString() throws IOException {
        String yaml =   "sources:\n" +
                        "    a1: b\n" +
                        "    a2: b\n" +
                        "    a3: b\n" +
                        "jobs:\n" +
                        "    - name: job1\n" +
                        "      scrapeRate: 10\n" +
                        "      sources:\n" +
                        "        - a1\n" +
                        "        - a2\n" +
                        "        - a3\n" +
                        "      fields:\n" +
                        "        a: DBasdf\n" +
                        "        b: DBbsdf\n";

        ScraperConfiguration configuration = mapper.readValue(yaml, ScraperConfigurationClassicImpl.class);

        assertThat(configuration.getJobConfigurations()).hasSize(1);
        JobConfiguration conf = configuration.getJobConfigurations().get(0);

        assertThat(configuration.getSources())
            .isNotEmpty()
            .hasSize(3)
            .containsEntry("a1", "b")
            .containsEntry("a2", "b")
            .containsEntry("a3", "b");

        assertThat(conf.getName()).isEqualTo("job1");
        assertThat(conf.getScrapeRate()).isEqualTo(10);
        assertThat(conf.getSources())
            .hasSize(3);

        assertThat(conf.getFields())
            .hasSize(2)
            .containsEntry("a", "DBasdf")
            .containsEntry("b", "DBbsdf");
    }

    @Test
    void parseJobs_missingEntries_fails() {
        String jobs =   "sources:\n" +
                        "    a: b\n" +
                        "jobs:\n" +
                        "    - name: job1\n" +
                        "      scrapeRate: 10\n" +
                        "      sources:\n" +
                        "        - a1\n";

        assertThatThrownBy(() -> mapper.readValue(jobs, ScraperConfigurationClassicImpl.class))
            .isInstanceOf(MismatchedInputException.class);
    }

    @Test
    void fromYaml_loads() {
        String yaml =   "sources:\n" +
                        "  a1: b\n" +
                        "  a2: b\n" +
                        "  a3: b\n" +
                        "jobs:\n" +
                        "  - name: job1\n" +
                        "    scrapeRate: 10\n" +
                        "    sources:\n" +
                        "      - a1\n" +
                        "      - a2\n" +
                        "      - a3\n" +
                        "    fields:\n" +
                        "      a: DBasdf\n" +
                        "      b: DBbsdf\n";

        assertThatCode(() -> ScraperConfiguration.fromYaml(yaml, ScraperConfigurationClassicImpl.class))
            .doesNotThrowAnyException();
    }

    @Test
    void fromString_loads() {
        String json =   "{\n" +
                        "    \"sources\": {\n" +
                        "        \"a1\": \"b\",\n" +
                        "        \"a2\": \"b\",\n" +
                        "        \"a3\": \"b\"\n" +
                        "    },\n" +
                        "    \"jobs\": [\n" +
                        "        {\n" +
                        "            \"name\": \"job1\",\n" +
                        "            \"scrapeRate\": 10,\n" +
                        "            \"sources\": [\n" +
                        "                \"a1\",\n" +
                        "                \"a2\",\n" +
                        "                \"a3\"\n" +
                        "            ],\n" +
                        "            \"fields\": {\n" +
                        "                \"a\": \"DBasdf\",\n" +
                        "                \"b\": \"DBbsdf\"\n" +
                        "            }\n" +
                        "        }\n" +
                        "    ]\n" +
                        "}";

        assertThatCode(() -> ScraperConfiguration.fromJson(json, ScraperConfigurationClassicImpl.class))
            .doesNotThrowAnyException();
    }

    @Test
    void new_notAllSourceAliasesAreResolvable_throws() {
        String yaml =   "sources:\n" +
                        "  b: c\n" +
                        "jobs:\n" +
                        "  - name: job1\n" +
                        "    scrapeRate: 10\n" +
                        "    sources:\n" +
                        "      - s1\n" +
                        "    fields:\n";

        assertThatThrownBy(() -> ScraperConfiguration.fromYaml(yaml, ScraperConfigurationClassicImpl.class))
            .isInstanceOf(PlcRuntimeException.class)
            .hasStackTraceContaining("unreferenced sources: [s1]");
    }

    @Test
    void generateScrapeJobs_fromConfig() throws ScraperException {
        String yaml =   "sources:\n" +
                        "  source1: 'connection string'\n" +
                        "jobs:\n" +
                        "  - name: job1\n" +
                        "    scrapeRate: 10\n" +
                        "    sources:\n" +
                        "      - source1\n" +
                        "    fields:\n" +
                        "      field1: 'DB1 Field 1'\n";

        List<ScrapeJob> jobs = ScraperConfiguration.fromYaml(yaml, ScraperConfigurationClassicImpl.class).getJobs();
        assertThat(jobs).hasSize(1);

        ScrapeJob job = jobs.get(0);

        assertThat(job.getJobName()).isEqualTo("job1");
        assertThat(job.getScrapeRate()).isEqualTo(10);
        assertThat(job.getSourceConnections())
            .hasSize(1)
            .containsEntry("source1", "connection string");
        assertThat(job.getFields())
            .hasSize(1)
            .containsEntry("field1", "DB1 Field 1");
    }

    @Nested
    class Files {

        @Test
        void json() throws IOException {
            ScraperConfiguration conf = ScraperConfiguration.fromFile("src/test/resources/config.json", ScraperConfigurationClassicImpl.class);
        }

        @Test
        void yaml() throws IOException {
            ScraperConfiguration conf = ScraperConfiguration.fromFile("src/test/resources/config.yml", ScraperConfigurationClassicImpl.class);
        }
    }
}