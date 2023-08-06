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
package org.apache.plc4x.merlot.db.command;

import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.karaf.shell.support.table.ShellTable;
import org.epics.pvdata.pv.ScalarType;


@Command(scope = "db", name = "scalars", description = "List default scalar types for PV database.")
@Service
public class DBScalarCommand implements Action {

    @Override
    public Object execute() throws Exception {
        ShowScalars();
        return null;
    }
        
    private void ShowScalars(){
        System.out.println("Basic scalar types supported by PV Database: ");
        ShellTable table = new ShellTable();
        table.column("PV Type");        
        table.column("PV Type (String)");        
        ScalarType[] scts = ScalarType.values();
        for (ScalarType sct:scts){
            table.addRow().addContent(sct.name(),sct.toString());
        }
        System.out.println();
        table.print(System.out);
        System.out.println();            
    }    
    
}
