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
package org.apache.plc4x.hop.metadata.util;

import io.netty.util.AbstractReferenceCounted;
import io.netty.util.ReferenceCounted;
import org.apache.plc4x.java.api.PlcConnection;

/**
 * This is a wrapper around a PlcConnection connection from the Apache Plc4x library.
 * Its goal is for it to be shared between the different transform that 
 * may be running simultaneously on the local or remote Hop engine. 
 * Its use with other kind of engines must be certified.
 * The concept is very simple, "the last to leave closes the connection".
 */
public class Plc4xWrapperConnection extends AbstractReferenceCounted {

    private final PlcConnection connection;
    private final String id;
            
    public Plc4xWrapperConnection(PlcConnection connection, String id){
        super();
        this.connection = connection;
        this.id = id;
    }
    
    public PlcConnection getConnection(){
        return this.connection;
    }    
        
    @Override
    protected void deallocate() {
        try{
            connection.close();
        } catch (Exception ex){
            ex.printStackTrace();
        }
    }
    
    @Override
    public ReferenceCounted touch(Object hint) {
        return this;
    }
    
    public String getId() {
        return id;
    }
    
}
