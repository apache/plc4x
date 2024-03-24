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
import {Menubar} from "primereact/menubar";
import {MenuItem} from "primereact/menuitem";
import plc4xLogo from "../assets/plc4x-logo.svg";
import {Image} from "primereact/image";
import {Outlet, useNavigate} from "react-router-dom";

export default function MainLayout() {
    const navigate = useNavigate();

    const menuItems = [
        {
            label: 'Inspect',
            command() {
                navigate('/inspect');
            }
        },
        {
            label: 'OPC-UA Server',
            command() {
                navigate('/opcua');
            }
        },
        {
            label: 'MQTT Emitter',
            command() {
                navigate('/mqtt');
            }
        },
        {
            label: 'Settings',
            command() {
                navigate('/settings');
            }
        },
        {
            label: 'About',
            command() {
                navigate('/about');
            }
        }
    ] as MenuItem[];

    const startLogo = <Image src={plc4xLogo} width="200px" className="m-3"/>;

    return (
        <>
            <nav>
                <Menubar className="flex w-12" model={menuItems} start={startLogo}/>
            </nav>
            <main className="h-full" style={{paddingTop: "10px"}}>
                <Outlet/>
            </main>
        </>
    )
}