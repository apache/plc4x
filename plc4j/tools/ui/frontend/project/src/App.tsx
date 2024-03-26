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

import './App.css'
import axios from 'axios';
import {createBrowserRouter, RouterProvider} from "react-router-dom";
import MainLayout from "./layouts/MainLayout.tsx";
import Inspect from "./pages/Inspect.tsx";
import OpcUa from "./pages/OpcUa.tsx";
import Mqtt from "./pages/Mqtt.tsx";
import Settings from "./pages/Settings.tsx";
import About from "./pages/About.tsx";
import useWebSocket from 'react-use-websocket';
import {useState} from "react";
import {DeviceEvent, RestApplicationClient} from "./generated/plc4j-tools-ui-frontend.ts";
import store, {
    addDevice,
    deleteDevice,
    DeviceAction,
    InitializeConnectionsAction,
    initializeLists,
    updateDevice
} from "./store";

axios.defaults.baseURL = 'http://localhost:8080';
const restClient = new RestApplicationClient(axios);

// We're actually just using this concept in order to separate the layout for the general page from the content.
const router = createBrowserRouter([
    {
        path: '/',
        element: <MainLayout/>,
        children: [
            {path: '/inspect', element: <Inspect/>},
            {path: '/opcua', element: <OpcUa/>},
            {path: '/mqtt', element: <Mqtt/>},
            {path: '/settings', element: <Settings/>},
            {path: '/about', element: <About/>},
        ]
    },
])

function App() {
    const [initialized, setInitialized] = useState(false)

    useWebSocket( 'ws://localhost:8080/ws', {
        onOpen: () => {
            console.log('WebSocket connection established.');
        },
        onMessage: event => {
            const deviceEvent = JSON.parse(event.data) as DeviceEvent;
            if(deviceEvent != null) {
                switch (deviceEvent.eventType) {
                    case "CREATED": {
                        const action:DeviceAction = {device: deviceEvent.source}
                        store.dispatch(addDevice(action));
                        break;
                    }
                    case "UPDATED": {
                        const action:DeviceAction = {device: deviceEvent.source}
                        store.dispatch(updateDevice(action));
                        break;
                    }
                    case "DELETED": {
                        const action: DeviceAction = {device: deviceEvent.source}
                        store.dispatch(deleteDevice(action));
                        break;
                    }
                }
            }
        }
    });

    // Load the initial list of drivers and connections and initialize the store with that.
    if(!initialized) {
        setInitialized(true);
        restClient.getAllDrivers().then(driverList => {
            restClient.getAllDevices().then(deviceList => {
                const action:InitializeConnectionsAction = {driverList:driverList.data, deviceList: deviceList.data}
                store.dispatch(initializeLists(action))
            })
        })
    }

    return (
        <RouterProvider router={router}/>
    )
}

export default App
