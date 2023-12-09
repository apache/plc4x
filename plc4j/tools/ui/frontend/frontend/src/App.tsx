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
import axios from 'axios';
import {createBrowserRouter, RouterProvider} from "react-router-dom";
import MainLayout from "./layouts/MainLayout.tsx";
import Editor from "./pages/Editor.tsx";

axios.defaults.baseURL = 'http://localhost:8080';

// We're actually just using this concept in order to separate the layout for the general page from the content.
const router = createBrowserRouter([
    {
        path: '/',
        element: <MainLayout/>,
        children: [
            {path: '/', element: <Editor/>},
        ]
    },
])

function App() {
    return (
        <RouterProvider router={router}/>
    )
}

export default App
