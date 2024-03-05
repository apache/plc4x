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

import {configureStore, createSlice, PayloadAction} from '@reduxjs/toolkit'
import {Device, Driver} from "../generated/plc4j-tools-ui-frontend.ts";
import {useDispatch} from "react-redux";

export type InitializeConnectionsAction = {
    driverList: Driver[]
    deviceList: Device[]
}

export type DeviceAction = {
    device: Device
}

export interface ConnectionsState {
    driverList: Driver[]
    deviceList: Device[]
}

const connectionsInitialState: ConnectionsState = {
    driverList: [] as Driver[],
    deviceList: [] as Device[],
}

const connectionsSlice = createSlice({
    name: 'connections',
    initialState: connectionsInitialState,
    reducers: {
        initializeLists: (state, action: PayloadAction<InitializeConnectionsAction>) => {
            state.driverList = action.payload.driverList
            state.deviceList = action.payload.deviceList
        },
        addDevice: (state, action: PayloadAction<DeviceAction>) => {
            console.log("ADD " + JSON.stringify(action))
            state.deviceList = [...state.deviceList, action.payload.device]
        },
        updateDevice: (state, action: PayloadAction<DeviceAction>) => {
            console.log("UPDATE " + JSON.stringify(action))
            const device = state.deviceList.find(value => value.id == action.payload.device.id);
            if (device) {
                const index = state.deviceList.indexOf(device);
                if (index) {
                    state.deviceList.splice(index, 1, action.payload.device);
                }
            }
        },
        deleteDevice: (state, action: PayloadAction<DeviceAction>) => {
            console.log("DELETE " + JSON.stringify(action))
            const device = state.deviceList.find(value => value.id == action.payload.device.id);
            if (device) {
                const index = state.deviceList.indexOf(device);
                if (index) {
                    state.deviceList.splice(index, 1);
                }
            }
        }
    }
})

export const { initializeLists, addDevice, updateDevice, deleteDevice } = connectionsSlice.actions

const store = configureStore({
    reducer: {
        connections: connectionsSlice.reducer
    }
})
export default store

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>
// Inferred type: {connections: ConnectionsState}
export type AppDispatch = typeof store.dispatch
export const useAppDispatch: () => AppDispatch = useDispatch
