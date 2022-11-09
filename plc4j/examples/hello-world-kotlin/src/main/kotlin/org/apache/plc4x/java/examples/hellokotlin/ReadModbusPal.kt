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
package org.apache.plc4x.java.examples.hellokotlin;

import org.apache.plc4x.java.PlcDriverManager
import org.apache.plc4x.java.api.types.PlcResponseCode
import java.util.concurrent.TimeUnit
import kotlin.system.exitProcess

fun main() {
    PlcDriverManager()
        .getConnection("modbus-tcp://localhost:502")
        .use { conn ->
            if (!conn.metadata.canRead()) {
                println("Cannot read!!")
                return
            }

            val readRequest = conn.readRequestBuilder()
                .addFieldAddress("value-1", "coil:1")
                .addFieldAddress("value-2", "coil:3[4]")
                .addFieldAddress("value-3", "holding-register:1")
                .addFieldAddress("value-4", "holding-register:3[4]")
                .build()

            val response = readRequest.execute().get(1, TimeUnit.MINUTES)
            response.fieldNames.forEach { fieldName ->
                val responseCode = response.getResponseCode(fieldName)
                if (responseCode !== PlcResponseCode.OK) {
                    println("Error[$fieldName]: ${responseCode.name}")
                    return
                }
                val numValues = response.getNumberOfValues(fieldName)
                // If it's just one element, output just one single line.
                if (numValues == 1) {
                    println("Value[$fieldName]: ${response.getObject(fieldName)}")
                } else {
                    println("Value[$fieldName]:")
                    for (i in 0 until numValues) {
                        println(" - " + response.getObject(fieldName, i))
                    }
                }
            }
        }

    exitProcess(0)
}