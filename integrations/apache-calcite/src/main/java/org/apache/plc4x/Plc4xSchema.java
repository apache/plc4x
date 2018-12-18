package org.apache.plc4x;

import org.apache.calcite.schema.Table;
import org.apache.calcite.schema.impl.AbstractSchema;
import org.apache.plc4x.java.scraper.ResultHandler;
import org.apache.plc4x.java.scraper.Scraper;
import org.apache.plc4x.java.scraper.config.JobConfiguration;
import org.apache.plc4x.java.scraper.config.ScraperConfiguration;

import java.util.Map;
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.BlockingQueue;
import java.util.stream.Collectors;

/**
 * Scraper -> Handler -> Table
 */
public class Plc4xSchema extends AbstractSchema {

    private final ScraperConfiguration configuration;
    private final Scraper scraper;
    private final QueueHandler handler;
    private final Map<String, BlockingQueue<Object[]>> queues;
    private final Map<String, Table> tableMap;

    Plc4xSchema(ScraperConfiguration configuration) {
        this.configuration = configuration;
        this.handler = new QueueHandler();
        this.scraper = new Scraper(configuration, handler);
        this.queues = configuration.getJobConfigurations().stream()
            .collect(Collectors.toMap(
                JobConfiguration::getName,
                conf -> new ArrayBlockingQueue<Object[]>(100)
            ));
        // Create the tables
        this.tableMap = configuration.getJobConfigurations().stream()
            .collect(Collectors.toMap(
                JobConfiguration::getName,
                conf -> new Plc4xTable(queues.get(conf.getName()), conf)
            ));
        // Start the scraper
        this.scraper.start();
    }

    @Override
    protected Map<String, Table> getTableMap() {
        // Return a map of all jobs
        return this.tableMap;
    }

    class QueueHandler implements ResultHandler {

        @Override
        public void handle(String job, String alias, Map<String, Object> results) {
            Object[] objects = results.values().toArray();
            try {
                queues.get(job).put(objects);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }

    }
}
