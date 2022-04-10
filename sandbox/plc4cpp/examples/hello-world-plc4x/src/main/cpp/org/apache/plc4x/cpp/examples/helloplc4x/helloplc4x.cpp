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
//package org.apache.plc4x.java.examples.helloplc4x;
#include <iostream>

#include <org/apache/plc4x/cpp/PlcDriverManager.h>
#include <org/apache/plc4x/cpp/api/PlcConnection.h>
#include <org/apache/plc4x/cpp/api/messages/PlcReadRequest.h>
#include <org/apache/plc4x/cpp/api/messages/PlcReadResponse.h>
#include <org/apache/plc4x/cpp/api/types/PlcResponseCode.h>
#include <org/apache/plc4x/cpp/utils/logger/BLogger.h>

using namespace std;
using namespace org::apache::plc4x::cpp;

/**
    * Example code do demonstrate using PLC4X.
    *
    * @param args ignored.
    */
int main(int argc, char *argv[]) 
{
    /* CliOptions options = CliOptions.fromArgs(args);
    if (options == null) {
        CliOptions.printHelp();
        // Could not parse.
        System.exit(1);
    }*/

    // Establish a connection to the plc using the url provided as first argument
    try 
    {

        PlcDriverManager plcDriverManager;
        //options.getConnectionString()

        PlcConnection* plcConnection = plcDriverManager.getConnection("s7://10.10.64.20/0/1");

        // Check if this connection support reading of data.
        /*if (!plcConnection.getMetadata().canRead()) {
            logger.error("This connection doesn't support reading.");
            return;
        }*/

        // Create a new read request:
        // - Give the single item requested the alias name "value"
        /*PlcReadRequest::Builder builder = plcConnection.readRequestBuilder();
        for (int i = 0; i < options.getFieldAddress().length; i++) {
            builder.addItem("value-" + i, options.getFieldAddress()[i]);
        }
        PlcReadRequest readRequest = builder.build();*/

        //////////////////////////////////////////////////////////
        // Read synchronously ...
        // NOTICE: the ".get()" immediately lets this thread pause until
        // the response is processed and available.*/
        LOG_INFO("Synchronous request ...");
        /*PlcReadResponse syncResponse = readRequest.execute().get();
        // Simply iterating over the field names returned in the response.
        printResponse(syncResponse);

        //////////////////////////////////////////////////////////
        // Read asynchronously ...
        // Register a callback executed as soon as a response arrives.
        logger.info("Asynchronous request ...");
        CompletableFuture<? extends PlcReadResponse> asyncResponse = readRequest.execute();
        asyncResponse.whenComplete((readResponse, throwable) -> {
            if (readResponse != null) {
                printResponse(readResponse);
            } else {
                logger.error("An error occurred: " + throwable.getMessage(), throwable);
            }
        }*/
    }
    catch (...)
    {

    }

    return 0;
}

/*private static void printResponse(PlcReadResponse response) {
    for (String fieldName : response.getFieldNames()) {
        if(response.getResponseCode(fieldName) == PlcResponseCode.OK) {
            int numValues = response.getNumberOfValues(fieldName);
            // If it's just one element, output just one single line.
            if(numValues == 1) {
                logger.info("Value[" + fieldName + "]: " + response.getObject(fieldName));
            }
            // If it's more than one element, output each in a single row.
            else {
                logger.info("Value[" + fieldName + "]:");
                for(int i = 0; i < numValues; i++) {
                    logger.info(" - " + response.getObject(fieldName, i));
                }
            }
        }
        // Something went wrong, to output an error message instead.
        else {
            logger.error("Error[" + fieldName + "]: " + response.getResponseCode(fieldName).name());
        }
    }
}*/


