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
package org.apache.plc4x.examples.plc4j.s7event;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class PlcReadAndWriteVIPA {

    private static final Logger logger = LoggerFactory.getLogger(PlcReadAndWriteVIPA.class);       
    
    /**
     * @param args the command line arguments
     */
    public static void main(String[] args) {
        logger.info("*****************************************************"); 
        logger.info("* 1. I need a real VIPA CPU.."); 
        logger.info("*    Press [ENTER]");        
        logger.info("*****************************************************");  
    }
    
}
