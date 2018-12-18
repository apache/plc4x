package org.apache.plc4x;

import org.apache.calcite.schema.Schema;
import org.apache.calcite.schema.SchemaFactory;
import org.apache.calcite.schema.SchemaPlus;
import org.apache.commons.lang3.Validate;
import org.apache.plc4x.java.scraper.config.ScraperConfiguration;

import java.io.IOException;
import java.util.Map;

public class Plc4xSchemaFactory implements SchemaFactory {

    @Override
    public Schema create(SchemaPlus parentSchema, String name, Map<String, Object> operand) {
        Object config = operand.get("config");
        Validate.notNull(config, "No configuration file given. Please specify one with 'config=...'");
        // Load configuration from file
        ScraperConfiguration configuration;
        try {
            configuration = ScraperConfiguration.fromFile(config.toString());
        } catch (IOException e) {
            throw new RuntimeException("Unable to load configuration file!", e);
        }
        // Pass the configuration to the Schema
        return new Plc4xSchema(configuration);
    }

}
