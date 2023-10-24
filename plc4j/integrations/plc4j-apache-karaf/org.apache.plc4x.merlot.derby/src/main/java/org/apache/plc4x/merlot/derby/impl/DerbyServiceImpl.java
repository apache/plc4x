/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.apache.plc4x.merlot.derby.impl;

import java.util.Properties;
import java.util.logging.Level;
import java.util.logging.Logger;
import org.apache.derby.drda.NetworkServerControl;
import org.apache.plc4x.merlot.derby.api.DerbyService;

public class DerbyServiceImpl implements DerbyService {

    NetworkServerControl server;
    
    @Override
    public void init() {
        try {            
            System.setProperty("derby.system.home", "data/derby");
            System.setProperty("derby.storage.tempDirectory", "data/derby/tmp");            
            server = new NetworkServerControl();           
            server.start(null); 
            Properties props = server.getCurrentProperties();
        } catch (Exception ex) {
            System.out.println(ex.getMessage());
            Logger.getLogger(DerbyServiceImpl.class.getName()).log(Level.SEVERE, null, ex);
        }
    }

    @Override
    public void destroy() {
        try {
            server.shutdown();
        } catch (Exception ex) {
            System.out.println(ex.getMessage());            
            Logger.getLogger(DerbyServiceImpl.class.getName()).log(Level.SEVERE, null, ex);
        }
    }
    
}