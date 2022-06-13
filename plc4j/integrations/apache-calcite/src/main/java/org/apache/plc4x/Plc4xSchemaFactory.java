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
package org.apache.plc4x;

import org.apache.calcite.schema.Schema;
import org.apache.calcite.schema.SchemaFactory;
import org.apache.calcite.schema.SchemaPlus;
import org.apache.commons.lang3.Validate;
import org.apache.plc4x.java.scraper.config.ScraperConfiguration;
import org.apache.plc4x.java.scraper.config.triggeredscraper.ScraperConfigurationTriggeredImpl;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.util.Map;

public class Plc4xSchemaFactory implements SchemaFactory {
    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4xSchemaFactory.class);

    @Override
    public Schema create(SchemaPlus parentSchema, String name, Map<String, Object> operand) {
        // Fetch config
        Object config = operand.get("config");
        Validate.notNull(config, "No configuration file given. Please specify operand 'config'...'");
        // Load configuration from file
        ScraperConfiguration configuration;
        try {
            configuration = ScraperConfiguration.fromFile(config.toString(), ScraperConfigurationTriggeredImpl.class);
        } catch (IOException e) {
            throw new IllegalArgumentException("Unable to load configuration file!", e);
        }

        // Fetch limit
        Object limit = operand.get("limit");
        Validate.notNull(limit, "No limit for the number of rows for a table. Please specify operand 'config'...'");
        long parsedLimit;
        try {
            parsedLimit = Long.parseLong(limit.toString());
        } catch (NumberFormatException e) {
            throw new IllegalArgumentException("Given limit '" + limit + "' cannot be parsed to valid long!", e);
        }
        // Pass the configuration to the Schema
        try {
            return new Plc4xSchema(configuration, parsedLimit);
        } catch (ScraperException e) {
            LOGGER.warn("Could not evaluate Plc4xSchema",e);
            //ToDo Exception, but interface does not accept ... null is fishy
            return null;
        }
    }

}
