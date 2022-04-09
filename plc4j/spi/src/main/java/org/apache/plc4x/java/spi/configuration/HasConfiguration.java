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
package org.apache.plc4x.java.spi.configuration;

/**
 * (Marker) Interface which can be used to tell PLC4X that a class (that is instantiated by PLC4X)
 * has a Configuration.
 * PLC4X will then try to Instantiate the Class and populate it based on Connection Parameters.
 *
 * @param <CONFIGURATION> Class of the Configuration
 */
public interface HasConfiguration<CONFIGURATION extends Configuration> {

    /**
     * Is called directly after instantiation before the class is used somewhere.
     */
    void setConfiguration(CONFIGURATION configuration);

}
