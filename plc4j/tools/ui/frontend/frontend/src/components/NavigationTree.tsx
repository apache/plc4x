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
import {Toast} from "primereact/toast";
import {ContextMenu} from "primereact/contextmenu";
import {RestApplicationClient} from "../generated/plc4j-tools-ui-frontend.ts";
import axios from "axios";

type NavigationTreeProps = {
    treeItems: TreeItemData[];
}

const restClient = new RestApplicationClient(axios);

export default function NavigationTree({treeItems}: NavigationTreeProps) {
    const [selectedNodeKey, setSelectedNodeKey] = useState<string>("");
    const toast = useRef<Toast>(null);
    const cm = useRef<ContextMenu>(null);
    const menu = [
        {
            label: 'Discover',
            icon: 'pi pi-search',
            command: () => {
                restClient.discover(selectedNodeKey);
            }
        },
        {
            label: 'Add',
            icon: 'pi pi-plus-circle',
            command: () => {
                console.log("Toggle")
            }
        }
    ];
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
            children: curItem.children?.map(value => createTreeNode(value))
        }
    }

    const treeNodes: TreeNode[] = treeItems.map(value => createTreeNode(value))
    return(
        <div>
            <Toast ref={toast} />

            <ContextMenu model={menu} ref={cm} onHide={() => setSelectedNodeKey("")} />

            <Tree value={treeNodes} selectionMode="single"
                  contextMenuSelectionKey={selectedNodeKey}
                  onContextMenuSelectionChange={event => setSelectedNodeKey(event.value as string)}
                  onContextMenu={event => {
                      if(cm.current) {
                          cm.current.show(event.originalEvent);
                      }
                  }}/>
        </div>)
}