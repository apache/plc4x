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

import {Action, createStore} from "redux";
import {Device, Driver} from "../generated/plc4j-tools-ui-frontend.ts";

export type ApplicationState = {
    driverList?: Driver[]
    deviceList?: Device[]
}

type InitializeAction = {
    type: string
    driverList: Driver[]
    deviceList: Device[]
}

type DeviceAction = {
    type: string
    device: Device
}

const storeReducer = (state: ApplicationState = {}, action: Action) => {
    console.log(action);
    if (action.type === 'initialize-lists') {
        const initializeAction: InitializeAction = action as InitializeAction
        return {
            ...state,
            deviceList: initializeAction.deviceList,
            driverList: initializeAction.driverList
        }
    } else if (action.type === 'add-device') {
        const deviceAction: DeviceAction = action as DeviceAction
        if(state.deviceList) {
            return {
                ...state,
                deviceList: [...state.deviceList, deviceAction.device]
            }
        }
        return {
            ...state,
            deviceList: [deviceAction.device]
        }
    } else if (action.type === 'update-device') {
        const deviceAction: DeviceAction = action as DeviceAction
        if(state.deviceList) {
            const device = state.deviceList.find(value => value.id == deviceAction.device.id);
            if(device) {
                const index = state.deviceList.indexOf(device);
                if(index) {
                    const newList = state.deviceList;
                    newList.splice(index, 1, deviceAction.device)
                    return {
                        ...state,
                        newList
                    }
                }
            }
        }
    } else if (action.type === 'delete-device') {
        const deviceAction: DeviceAction = action as DeviceAction
        if(state.deviceList) {
            const device = state.deviceList.find(value => value.id == deviceAction.device.id);
            if(device) {
                const index = state.deviceList.indexOf(device);
                if(index) {
                    const newList = state.deviceList;
                    newList.splice(index, 1)
                    return {
                        ...state,
                        newList
                    }
                }
            }
        }
    }
    return state
}

const store = createStore<ApplicationState, Action>(storeReducer)

export default store;
