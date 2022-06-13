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
import org.apache.calcite.linq4j.Enumerable;
import org.apache.calcite.schema.ScannableTable;
import org.apache.plc4x.java.scraper.config.JobConfiguration;

import java.util.concurrent.BlockingQueue;

public class Plc4xTable extends Plc4xBaseTable implements ScannableTable {

    public Plc4xTable(BlockingQueue<Plc4xSchema.Record> queue, JobConfiguration conf, long tableCutoff) {
        super(queue, conf, tableCutoff);
    }

    @Override
    public Enumerable<Object[]> scan(DataContext root) {
        return super.scan(root);
    }

}
