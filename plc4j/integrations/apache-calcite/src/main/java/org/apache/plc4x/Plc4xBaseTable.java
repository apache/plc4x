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

import org.apache.calcite.DataContext;
import org.apache.calcite.linq4j.AbstractEnumerable;
import org.apache.calcite.linq4j.Enumerable;
import org.apache.calcite.linq4j.Enumerator;
import org.apache.calcite.rel.*;
import org.apache.calcite.rel.type.RelDataType;
import org.apache.calcite.rel.type.RelDataTypeFactory;
import org.apache.calcite.schema.Statistic;
import org.apache.calcite.schema.impl.AbstractTable;
import org.apache.calcite.util.ImmutableBitSet;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.scraper.config.JobConfiguration;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.Timestamp;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.Supplier;
import java.util.stream.Collectors;
import java.util.stream.Stream;

/**
 * Base for Stream and "Table" version of the Plc4xTable.
 * Needs to be subclassed due to usage of "instanceof" in Calcites internals.
 */
public abstract class Plc4xBaseTable extends AbstractTable {

    private static final Logger logger = LoggerFactory.getLogger(Plc4xBaseTable.class);

    private final BlockingQueue<Plc4xSchema.Record> queue;
    private final JobConfiguration conf;
    private final long tableCutoff;
    private Plc4xSchema.Record current;
    private final List<String> names;

    public Plc4xBaseTable(BlockingQueue<Plc4xSchema.Record> queue, JobConfiguration conf, long tableCutoff) {
        this.tableCutoff = tableCutoff;
        logger.info("Instantiating new PLC4X Table with configuration: {}", conf);
        this.queue = queue;
        this.conf = conf;
        // Extract names
        names = new ArrayList<>(conf.getTags().keySet());
    }

    @Override
    public Statistic getStatistic() {
        return new Statistic() {

            public Double getRowCount() {
                return tableCutoff > 0 ? (double)tableCutoff : null;
            }

            public boolean isKey(ImmutableBitSet columns) {
                return false;
            }

            public List<RelReferentialConstraint> getReferentialConstraints() {
                return Collections.emptyList();
            }

            public List<RelCollation> getCollations() {
                return Collections.singletonList(RelCollationImpl.of(new RelFieldCollation(0, RelFieldCollation.Direction.ASCENDING)));
            }

            public RelDistribution getDistribution() {
                return RelDistributionTraitDef.INSTANCE.getDefault();
            }
        };
    }

    @Override
    public RelDataType getRowType(RelDataTypeFactory typeFactory) {
        // Create the table spec
        // Block until the first result is in the queue
        CompletableFuture<Plc4xSchema.Record> future = CompletableFuture.supplyAsync(new FirstElementFetcher(queue));
        Plc4xSchema.Record first;
        try {
            first = future.get(5, TimeUnit.SECONDS);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcRuntimeException("Thread was interrupted!", e);
        } catch (ExecutionException | TimeoutException e) {
            throw new PlcRuntimeException("Unable to fetch first record and infer arguments!", e);
        }
        logger.info("Inferring types for Table '{}' based on values: {}", conf.getName(), first.values);
        // Extract types
        List<RelDataType> types = names.stream()
            .map(n -> {
                Object o = first.values.get(n);
                logger.debug("Infer tag '{}' as class '{}'", n, o.getClass());
                return typeFactory.createJavaType(o.getClass());
            })
            .collect(Collectors.toList());
        List<String> pre = new ArrayList<>(Arrays.asList("timestamp", "source"));
        pre.addAll(names);
        List<RelDataType> preTypes = Stream.of(Timestamp.class, String.class)
            .map(typeFactory::createJavaType)
            .collect(Collectors.toList());
        preTypes.addAll(types);
        return typeFactory.createStructType(preTypes, pre);
    }

    /**
     * if tableCutoff is positive, then the row gets limited to that.
     */
    public Enumerable<Object[]> scan(DataContext root) {
        return new AbstractEnumerable<>() {
            @Override
            public Enumerator<Object[]> enumerator() {
                return new Enumerator<>() {

                    private final AtomicLong counter = new AtomicLong(0);

                    @Override
                    public Object[] current() {
                        List<Object> objects = new ArrayList<>(Arrays.asList(new Timestamp(current.timestamp.toEpochMilli()), current.source));
                        List<Object> objects2 = names.stream().map(name -> current.values.get(name)).collect(Collectors.toList());
                        objects.addAll(objects2);
                        return objects.toArray();
                    }

                    @Override
                    public boolean moveNext() {
                        try {
                            current = queue.take();
                            // If stream, simply return
                            if (tableCutoff <= 0L) {
                                return true;
                            }
                            // If table, return if below cutoff
                            return counter.getAndIncrement() < tableCutoff;
                        } catch (InterruptedException e) {
                            Thread.currentThread().interrupt();
                        }
                        return false;
                    }

                    @Override
                    public void reset() {
                        counter.set(0);
                    }

                    @Override
                    public void close() {
                        // Unimplemented
                    }
                };
            }
        };
    }

    /**
     * Waits until a first (non null) element is in the queue
     */
    private static class FirstElementFetcher implements Supplier<Plc4xSchema.Record> {

        private final BlockingQueue<Plc4xSchema.Record> queue;

        private FirstElementFetcher(BlockingQueue<Plc4xSchema.Record> queue) {
            this.queue = queue;
        }

        @Override
        public Plc4xSchema.Record get() {
            Plc4xSchema.Record first;
            do {
                first = queue.peek();
            } while (first == null);
            return first;
        }
    }

}
