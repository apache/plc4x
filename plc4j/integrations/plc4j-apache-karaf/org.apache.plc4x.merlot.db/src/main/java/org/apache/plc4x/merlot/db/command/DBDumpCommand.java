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

import java.util.List;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Argument;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.karaf.shell.support.table.ShellTable;
import org.epics.pvdatabase.PVDatabase;
import org.epics.pvdatabase.PVRecord;
import org.osgi.framework.BundleContext;


@Command(scope = "db", name = "dump", description = "Dump information from a list of records.")
@Service
public class DBDumpCommand  implements Action {
    @Reference
    BundleContext bundleContext;
    
    @Reference
    PVDatabase master;

    @Argument(index = 0, name = "records", description = "List of records to dump information", required = true, multiValued = true)
    List<String> records;    
    
    @Override
    public Object execute() throws Exception {
        for (String record:records){
            PVRecord pvRecord = master.findRecord(record);
            PrintPVRecord(record, pvRecord);
        }
        return null;
    }

    private void PrintPVRecord(String record, PVRecord pvRecord) {
        ShellTable table = new ShellTable();
        table.column("Field");
        table.column("Value");;
        table.emptyTableText("Record not found.");
        if (!(pvRecord == null)){
                System.out.println(pvRecord.toString()+"\r\n");
            } else {
                System.out.println("Record \"" + record + "\" not found.");
            }
        //table.print(System.out);
    }
           
}
