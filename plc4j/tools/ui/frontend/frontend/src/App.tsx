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

import './App.css'
import NavigationTree from "./components/NavigationTree.tsx";
import axios from 'axios';
import {useState} from "react";
import {RestApplicationClient} from "./generated/plc4j-tools-ui-frontend.ts";
import {TreeItemData} from "./model/TreeItemData.ts";
import {TabPanel, TabView} from "primereact/tabview";
import {Menubar} from "primereact/menubar";
import {MenuItem} from "primereact/menuitem";
import plc4xLogo from "./assets/plc4x-logo.svg";
import { Splitter, SplitterPanel } from 'primereact/splitter';
import {ScrollPanel} from "primereact/scrollpanel";

axios.defaults.baseURL = 'http://localhost:8080';

function App() {
    const [driverTreeRoots, setDriverTreeRoots] = useState<TreeItemData[]>([])
    const restClient = new RestApplicationClient(axios);
    const start = <img alt="logo" src={plc4xLogo} height="40" className="mr-8"></img>;
    const menuModel: MenuItem[] = [
        /*{
            label: 'File',
            icon: 'pi pi-fw pi-file',
            items: [
                {
                    label: 'New',
                    icon: 'pi pi-fw pi-plus',
                    items: [
                        {
                            label: 'Bookmark',
                            icon: 'pi pi-fw pi-bookmark'
                        },
                        {
                            label: 'Video',
                            icon: 'pi pi-fw pi-video'
                        },

                    ]
                },
                {
                    label: 'Delete',
                    icon: 'pi pi-fw pi-trash'
                },
                {
                    separator: true
                },
                {
                    label: 'Export',
                    icon: 'pi pi-fw pi-external-link'
                }
            ]
        },
        {
            label: 'Edit',
            icon: 'pi pi-fw pi-pencil',
            items: [
                {
                    label: 'Left',
                    icon: 'pi pi-fw pi-align-left'
                },
                {
                    label: 'Right',
                    icon: 'pi pi-fw pi-align-right'
                },
                {
                    label: 'Center',
                    icon: 'pi pi-fw pi-align-center'
                },
                {
                    label: 'Justify',
                    icon: 'pi pi-fw pi-align-justify'
                },

            ]
        },
        {
            label: 'Users',
            icon: 'pi pi-fw pi-user',
            items: [
                {
                    label: 'New',
                    icon: 'pi pi-fw pi-user-plus',

                },
                {
                    label: 'Delete',
                    icon: 'pi pi-fw pi-user-minus',

                },
                {
                    label: 'Search',
                    icon: 'pi pi-fw pi-users',
                    items: [
                        {
                            label: 'Filter',
                            icon: 'pi pi-fw pi-filter',
                            items: [
                                {
                                    label: 'Print',
                                    icon: 'pi pi-fw pi-print'
                                }
                            ]
                        },
                        {
                            icon: 'pi pi-fw pi-bars',
                            label: 'List'
                        }
                    ]
                }
            ]
        },
        {
            label: 'Events',
            icon: 'pi pi-fw pi-calendar',
            items: [
                {
                    label: 'Edit',
                    icon: 'pi pi-fw pi-pencil',
                    items: [
                        {
                            label: 'Save',
                            icon: 'pi pi-fw pi-calendar-plus'
                        },
                        {
                            label: 'Delete',
                            icon: 'pi pi-fw pi-calendar-minus'
                        }
                    ]
                },
                {
                    label: 'Archive',
                    icon: 'pi pi-fw pi-calendar-times',
                    items: [
                        {
                            label: 'Remove',
                            icon: 'pi pi-fw pi-calendar-minus'
                        }
                    ]
                }
            ]
        },*/
        {
            label: 'Dummy',
            icon: 'pi pi-fw pi-file',
        }
    ];

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
        <div className="flex flex-column w-screen h-screen">
            <header>
                <Menubar start={start} model={menuModel}/>
            </header>
            <Splitter className="flex h-full">
                <SplitterPanel className="flex">
                    <TabView style={{width: '100%', height:'100%'}}>
                        <TabPanel header="By Driver" className="m-0">
                            <ScrollPanel style={{width: '100%', height:'100%'}} className="h-full">
                                <NavigationTree treeItems={driverTreeRoots}/>
                            </ScrollPanel>
                        </TabPanel>
                        <TabPanel header="By Device">
                            <ScrollPanel style={{width: '100%', height:'100%'}}>
                                <NavigationTree treeItems={driverTreeRoots}/>
                            </ScrollPanel>
                        </TabPanel>
                    </TabView>
                </SplitterPanel>
                <SplitterPanel className="w-full">
                    <TabView>

                        <TabPanel key={"ads://192.168.23.20"} header="ads://192.168.23.20" closable={true}>
                            <p>test</p>
                        </TabPanel>
                    </TabView>
                </SplitterPanel>
            </Splitter>
        </div>
    )
}

export default App
