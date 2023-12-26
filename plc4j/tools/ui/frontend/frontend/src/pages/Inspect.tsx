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
import {TabPanel, TabView} from "primereact/tabview";
import {Splitter, SplitterPanel} from "primereact/splitter";
import {useState} from "react";
import {RestApplicationClient} from "../generated/plc4j-tools-ui-frontend.ts";
import axios from "axios";
import {TreeItemData} from "../model/TreeItemData.ts";
import {ScrollPanel} from "primereact/scrollpanel";
import NavigationTree from "../components/NavigationTree.tsx";
import PlcConnection from "../components/PlcConnection.tsx";

export default function Inspect() {
    const [driverTreeRoots, setDriverTreeRoots] = useState<TreeItemData[]>([])
    const restClient = new RestApplicationClient(axios);

    function updateDriverList() {
        const driverList = restClient.getDriverList();
        driverList.then(response => {
            let newDriverTreeRoots: TreeItemData[] = [];
            response.data.map(driverValue => {
                newDriverTreeRoots = [...newDriverTreeRoots, {
                    id: driverValue.code,
                    name: driverValue.name,
                    type: "DRIVER",
                    supportsDiscovery: driverValue.supportsDiscovery,
                    supportsBrowsing: false,
                    supportsReading: false,
                    supportsWriting: false,
                    supportsSubscribing: false,
                    supportsPublishing: false,
                }]
            })
            setDriverTreeRoots(newDriverTreeRoots)
        })
    }

    if(driverTreeRoots.length == 0) {
        updateDriverList()
    }

    return (
        <Splitter className="h-full">
            <SplitterPanel
                size={16} minSize={1}
                className="flex align-items-center justify-content-center">
                <TabView style={{width: '100%', height:'100%'}}>
                    <TabPanel header="By Driver" className="m-0">
                        <ScrollPanel style={{width: '100%', height:'100%'}} className="h-full">
                            <NavigationTree treeItems={driverTreeRoots}/>
                        </ScrollPanel>
                    </TabPanel>
                    <TabPanel header="By Device">
                        <ScrollPanel style={{width: '100%', height:'100%'}}>
                            <NavigationTree  treeItems={driverTreeRoots}/>
                        </ScrollPanel>
                    </TabPanel>
                </TabView>
            </SplitterPanel>
            <SplitterPanel size={84} className="flex align-items-center justify-content-center">
                <TabView className="h-full w-full" panelContainerClassName="h-full">
                    <PlcConnection connectionString={"ads://192.168.23.20"}/>
                </TabView>
            </SplitterPanel>
        </Splitter>
    )
}
