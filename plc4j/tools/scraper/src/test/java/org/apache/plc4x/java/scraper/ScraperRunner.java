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
package org.apache.plc4x.java.scraper;

import org.apache.plc4x.java.scraper.config.ScraperConfiguration;
import org.apache.plc4x.java.scraper.config.ScraperConfigurationClassicImpl;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;

public class ScraperRunner {

    private static final Logger LOGGER = LoggerFactory.getLogger(ScraperRunner.class);

    public static void main(String[] args) throws IOException, ScraperException {
        ScraperConfiguration configuration = ScraperConfiguration.fromFile("plc4j/utils/scraper/src/test/resources/example.yml", ScraperConfigurationClassicImpl.class);
        Scraper scraper = new ScraperImpl(configuration, (j, a, m) -> LOGGER.info("Results from {}/{}: {}", j, a, m));

        scraper.start();
    }
}
