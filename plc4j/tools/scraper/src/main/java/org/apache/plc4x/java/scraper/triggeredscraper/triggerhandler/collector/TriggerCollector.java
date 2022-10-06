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
package org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector;

import org.apache.plc4x.java.scraper.exception.ScraperException;

/**
 * defines the interface for implementing a TriggerCollector
 * that handles and acquires all triggerRequests at once that needs a PlcConnection
 */
public interface TriggerCollector {

    /**
     * submits a trigger request to TriggerCollector
     * @param plcField a (plc) field that is used for triggering procedure
     * @param plcConnectionString the connection string to the regarding source
     * @param maxAwaitingTime max awaiting time until request shall be submitted
     * @return a uuid under that the request is handled internally
     * @throws ScraperException something went wrong
     */
    String submitTrigger(String plcField, String plcConnectionString, long maxAwaitingTime) throws ScraperException;

    /**
     * requests the result of submitted plc request with default timeout
     * @param uuid uuid that represents the request
     * @return the object acquired by requesting plc instance
     * @throws ScraperException something went wrong
     */
    Object requestResult(String uuid) throws ScraperException;

    /**
     * requests the result of submitted plc request
     * @param uuid uuid that represents the request
     * @param timeout timeout until response shall be acquired
     * @return the object acquired by requesting plc instance
     * @throws ScraperException something went wrong
     */
    Object requestResult(String uuid, long timeout) throws ScraperException;

    /**
     * starts the acquirement of triggers
     */
    void start();

    /**
     * stops acquirement of triggers
     */
    void stop();
}
