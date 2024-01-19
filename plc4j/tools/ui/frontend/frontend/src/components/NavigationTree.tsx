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

import {TreeItemData} from "../model/TreeItemData.ts";
import {Tree} from "primereact/tree";
import {TreeNode} from "primereact/treenode";
import {IconType} from "primereact/utils";
import 'primeicons/primeicons.css';
import {useRef, useState} from "react";
import {ContextMenu} from "primereact/contextmenu";
import {MenuItem} from "primereact/menuitem";
import {Device, RestApplicationClient} from "../generated/plc4j-tools-ui-frontend.ts";
import axios from "axios";
import {Counter} from "../utils/Counter.ts";
import DeviceDialog from "./DeviceDialog.tsx";

type NavigationTreeProps = {
    treeItems: TreeItemData[];
}

const restClient = new RestApplicationClient(axios);

export default function NavigationTree({treeItems}: NavigationTreeProps) {
    const dialogDevice = useRef<Device>({
        id: 0,
        name: "",
        protocolCode: "",
        transportCode: "",
        transportUrl: "",
        options: {},
        attributes: {},
    });
    const [showDeviceDialog, setShowDeviceDialog] = useState(false)

    const cm = useRef<ContextMenu>(null);
    const menu = [
        {
            key: "1",
            label: 'Discover',
            data: "discover-data",
            icon: 'pi pi-search',
            disabled: false,
        } as MenuItem,
        {
            key: "2",
            label: 'Add',
            data: "add-data",
            icon: 'pi pi-plus-circle',
            disabled: false,
        } as MenuItem,
        {
            key: "3",
            label: 'Edit',
            data: "edit-data",
            icon: 'pi pi-pencil',
            disabled: false,
        } as MenuItem,
        {
            key: "4",
            label: 'Delete',
            data: "delete-data",
            icon: 'pi pi-minus-circle',
            disabled: false,
        } as MenuItem,
        {
            key: "5",
            label: 'Connect',
            data: "connect-data",
            icon: 'pi pi-phone',
            disabled: false,
        } as MenuItem
    ] as MenuItem[]

    function updateMenu(selectedItem: TreeItemData) {
        // Discover
        menu[0].disabled = !selectedItem.supportsDiscovery && selectedItem.type != "ROOT"
        menu[0].command = () => {
            restClient.discover(selectedItem.id)
        }

        // Add
        menu[1].disabled = selectedItem.type != "DRIVER"
        menu[1].command = () => {
            dialogDevice.current.id = 0;
            dialogDevice.current.name = "";
            dialogDevice.current.transportCode = "";
            dialogDevice.current.transportUrl = "";
            dialogDevice.current.protocolCode = selectedItem.id;
            dialogDevice.current.options = {};
            dialogDevice.current.attributes = {};
            setShowDeviceDialog(true);
        }

        // Edit
        menu[2].disabled = selectedItem.type != "DEVICE"
        menu[2].command = () => {
            console.log(JSON.stringify(selectedItem))
            if(selectedItem.device) {
                dialogDevice.current.id = selectedItem.device.id;
                dialogDevice.current.name = selectedItem.device.name;
                dialogDevice.current.transportCode = selectedItem.device.transportCode;
                dialogDevice.current.transportUrl = selectedItem.device.transportUrl;
                dialogDevice.current.protocolCode = selectedItem.device.protocolCode;
                dialogDevice.current.options = selectedItem.device.options;
                dialogDevice.current.attributes = selectedItem.device.attributes;
                setShowDeviceDialog(true);
            }
        }
        // Delete
        menu[3].disabled = selectedItem.type != "DEVICE"
        menu[3].command = () => {
            // TODO: Ask for confirmation ...
        }
        // Connect
        menu[4].disabled = selectedItem.type != "DEVICE"
        menu[4].command = () => {
            // TODO: Open a new tab with the connection ...
        }
    }

    function getIcon(curItem: TreeItemData): IconType<TreeNode> {
        switch (curItem.type) {
            case "DRIVER":
                return "pi pi-fw pi-folder-open"//"material-icons md-18 folder_open"
            case "CONNECTION":
                return "pi pi-fw pi-phone"//"material-icons md-18 tty"
            case "DEVICE":
                return "pi pi-fw"//"material-icons md-18 computer"

            // discover:  "Radar"
            // browse:    "Manage Search"
            // read:      "Keyboard Arrow Down"
            // write:     "Keyboard Arrow Up"
            // subscribe: "Keyboard Double Arrow Down"
            // publish:   "Keyboard Double Arrow Up"
        }
    }

    function createTreeNode(curItem: TreeItemData, keyGenerator: Counter): TreeNode {
        return {
            key: keyGenerator.getAndIncrement(),
            id: curItem.id,
            label: curItem.name,
            icon: getIcon(curItem),
            data: curItem,
            children: curItem.children?.map(value => createTreeNode(value, keyGenerator))
        }
    }

    const treeNodes: TreeNode[] = treeItems.map(value => createTreeNode(value, new Counter()))
    return (
        <div>
            <DeviceDialog device={dialogDevice.current} visible={showDeviceDialog} />
            <ContextMenu model={menu} ref={cm}/>
            <Tree value={treeNodes}
                  selectionMode="single"
                  onContextMenu={event => {
                      if (cm.current) {
                          // Update the state of the menu (enabling/disabling some menu items)
                          updateMenu(event.node.data as TreeItemData)
                          //cm.current.props.model = menu
                          cm.current.show(event.originalEvent);
                      }
                  }}/>
        </div>)
}