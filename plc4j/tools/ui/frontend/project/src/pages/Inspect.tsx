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
import {TabPanel, TabView} from "primereact/tabview";
import {Splitter, SplitterPanel} from "primereact/splitter";
import {ScrollPanel} from "primereact/scrollpanel";
import NavigationTree from "../components/NavigationTree.tsx";
import PlcConnection from "../components/PlcConnection.tsx";
import { useSelector } from "react-redux";
import {Device, Driver} from "../generated/plc4j-tools-ui-frontend.ts";
import {TreeItemData} from "../model/TreeItemData.ts";
import {RootState} from "../store";

function getByDriverTree(driverList: Driver[], deviceList: Device[]):TreeItemData[] {
    console.log("getByDriverTree " + JSON.stringify(driverList) + " " + JSON.stringify(deviceList))
    if(driverList && deviceList) {
        const driverMap = new Map<string, TreeItemData>()
        let drivers:TreeItemData[] = []
        driverList.forEach(value => {
            const driverEntry:TreeItemData = {
                id: value.code,
                name: value.name,
                type: "DRIVER",
                driver: value,
                supportsDiscovery: value.metadata.discoverySupported,
                supportsBrowsing: false,
                supportsReading: false,
                supportsWriting: false,
                supportsSubscribing: false,
                supportsPublishing: false,
                children: []
            }
            driverMap.set(value.code, driverEntry)
            drivers = [...drivers, driverEntry]
        })
        drivers.sort((a,b) => (a.name > b.name) ? 1 : ((b.name > a.name) ? -1 : 0))

        deviceList.forEach(value => {
            const curDriverTreeItem = driverMap.get(value.protocolCode);
            if(curDriverTreeItem && curDriverTreeItem.children) {
                curDriverTreeItem.children = [...curDriverTreeItem.children, {
                    id: value.id?.toString(),
                    name: value.name,
                    type: "DEVICE",
                    device: value,
                    // Even if most connections will support reading/writing, in order to really know, we need a
                    // connection first. As we're not offering any connection-actions here anyway we just set all
                    // to false.
                    supportsDiscovery: false,
                    supportsBrowsing: false,
                    supportsReading: false,
                    supportsWriting: false,
                    supportsSubscribing: false,
                    supportsPublishing: false} as TreeItemData
                ]
            }
        })

        return [{
            id: "all",
            name: "Root",
            type: "ROOT",
            supportsDiscovery: true,
            supportsBrowsing: false,
            supportsReading: false,
            supportsWriting: false,
            supportsSubscribing: false,
            supportsPublishing: false,
            children: drivers
        }]
    }
    return []
}

function getByDeviceTree(deviceList: Device[]):TreeItemData[] {
    console.log("getByDeviceTree " + JSON.stringify(deviceList))
    if(deviceList) {
        // Group the connections by transport-url.
        // TODO: Possibly create filters for the different types of urls (IP, Hostname, Port, Mac-Address, ...)
        const deviceMap = new Map<string, Device[]>
        deviceList.forEach(device => {
            const devices = deviceMap.get(device.transportUrl)
            if(devices) {
                deviceMap.set(device.transportUrl, [...devices, device])
            } else {
                deviceMap.set(device.transportUrl, [device])
            }
        })

        // Build a tree based on the grouped locations.
    }
    return []
}

export default function Inspect() {
    const lists = useSelector((state: RootState) => {
        console.log("Updated connections " + JSON.stringify(state.connections))
        return state.connections
    })
    return (
        <Splitter className="h-full">
            <SplitterPanel
                size={16} minSize={1}
                className="flex align-items-center justify-content-center">
                <TabView style={{width: '100%', height:'100%'}}>
                    <TabPanel header="By Driver" className="m-0">
                        <ScrollPanel style={{width: '100%', height:'100%'}} className="h-full">
                            <NavigationTree treeItems={getByDriverTree(lists.driverList, lists.deviceList)}/>
                        </ScrollPanel>
                    </TabPanel>
                    <TabPanel header="By Device">
                        <ScrollPanel style={{width: '100%', height:'100%'}}>
                            <NavigationTree  treeItems={getByDeviceTree(lists.deviceList)}/>
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
