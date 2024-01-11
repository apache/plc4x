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
import {useRef} from "react";
import {ContextMenu} from "primereact/contextmenu";
import {MenuItem} from "primereact/menuitem";
//import {RestApplicationClient} from "../generated/plc4j-tools-ui-frontend.ts";
//import axios from "axios";

type NavigationTreeProps = {
    treeItems: TreeItemData[];
}

//const restClient = new RestApplicationClient(axios);

export default function NavigationTree({treeItems}: NavigationTreeProps) {
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
        } as MenuItem
    ] as MenuItem[]
    function updateMenu(selectedItem:TreeItemData) {
        menu[0].disabled = selectedItem.supportsDiscovery
        menu[1].disabled = selectedItem.type != "DRIVER"
        console.log("Second element enabled: " + (menu[1].disabled))
    }
    function getIcon(curItem: TreeItemData):IconType<TreeNode> {
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
    function createTreeNode(curItem:TreeItemData):TreeNode {
        return {
            id: curItem.id,
            label: curItem.name,
            icon: getIcon(curItem),
            data: curItem,
            children: curItem.children?.map(value => createTreeNode(value))
        }
    }

    const treeNodes: TreeNode[] = treeItems.map(value => createTreeNode(value))
    return(
        <div>
            <ContextMenu model={menu} ref={cm}/>

            <Tree value={treeNodes}
                  selectionMode="single"
                  onContextMenu={event => {
                      if(cm.current) {
                          // Update the state of the menu (enabling/disabling some menu items)
                          updateMenu(event.node.data as TreeItemData)
                          //cm.current.props.model = menu
                          cm.current.show(event.originalEvent);
                      }
                  }}/>
        </div>)
}