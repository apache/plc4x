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

import 'primeicons/primeicons.css';
import {TabPanel} from "primereact/tabview";
import {TreeTable} from "primereact/treetable";
import {Column} from "primereact/column";

type PlcConnectionProps = {
    connectionString: string;
}

export default function PlcConnection({connectionString}: PlcConnectionProps) {
    return (
        <TabPanel key={connectionString} header={connectionString} closable={true} contentClassName="h-full">
            <TreeTable>
                <Column field="name" header="Name" expander/>
                <Column field="dataType" header="Datatype"/>
                <Column field="size" header="Size"/>
                <Column field="readable" header="Readable"/>
                <Column field="writable" header="Writable"/>
                <Column field="subscribable" header="Subscribable"/>
                <Column field="publishable" header="Publishable"/>
                <Column field="value" header="Value"/>
                <Column field="setValue" header="Set Value"/>
            </TreeTable>
        </TabPanel>
    )
}