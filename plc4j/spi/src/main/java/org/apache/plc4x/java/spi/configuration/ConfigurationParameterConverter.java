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
package org.apache.plc4x.java.spi.configuration;

/**
 * Interface which allows to convert parameter from URI into its complex form.
 */
public interface ConfigurationParameterConverter<T> {

    /**
     * Type of supported configuration parameter.
     *
     * Returned value determines Java type to which this converter is able to turn string representation. Only if field
     * type is assignable to returned type conversion attempt will be made.
     *
     * @return Java type constructed by converter.
     */
    Class<T> getType();

    /**
     * Executes conversion of parameter textual representation into java object.
     *
     * @param value Parameter value.
     * @return Object representing passed string value.
     */
    T convert(String value);

}
