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

import org.apache.plc4x.java.api.PlcConnection;
import org.openide.util.lookup.AbstractLookup;
import org.openide.util.lookup.InstanceContent;
import org.openide.util.lookup.InstanceContent.Convertor;

/*
* Plc4xLookup is the main repository for sharing connections within a 
* Hop pipeline/workflow instance.
* It's static nature limits its use within a JVM instance, 
* however this is not a limitation since in general the connections to 
* the PLCs are limited resources.
* A "Convertor" is included for the "Plc4xWrapperConnection" types 
* to facilitate their use, without limiting their use as a container
* for other types of objects.
*/
public class Plc4xLookup extends AbstractLookup {
    private InstanceContent content = null;
    private static Plc4xLookup INSTANCE = new Plc4xLookup (); 
    
    public Plc4xLookup(InstanceContent content) {
        super(content);
        this.content = content;
    }

    public Plc4xLookup() {
        this(new InstanceContent());
    }

    public void add(Object obj) {
        if (obj instanceof Plc4xWrapperConnection)  
            content.add((Plc4xWrapperConnection) obj, Plc4xConvertor);
        else 
            content.add(obj);
    }

    public void remove(Object obj) {
        if (obj instanceof Plc4xWrapperConnection)  
            content.remove((Plc4xWrapperConnection) obj, Plc4xConvertor);
        else 
            content.remove(obj);
    }

    public static Plc4xLookup getDefault(){
        return INSTANCE;
    }  
 
    private static final Convertor<Plc4xWrapperConnection , Plc4xWrapperConnection>  Plc4xConvertor = 
            new InstanceContent.Convertor<Plc4xWrapperConnection , Plc4xWrapperConnection>() {
                
        @Override
        public Plc4xWrapperConnection convert(Plc4xWrapperConnection obj) {
            return obj;
        }

        @Override
        public Class<? extends Plc4xWrapperConnection> type(Plc4xWrapperConnection obj) {
            return Plc4xWrapperConnection.class;
        }

        @Override
        public String id(Plc4xWrapperConnection obj) {
            return obj.getId();
        }

        @Override
        public String displayName(Plc4xWrapperConnection obj) {
            return obj.getConnection().toString();
        }
    
    };
     
}
