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
