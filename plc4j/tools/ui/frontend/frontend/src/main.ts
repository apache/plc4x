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
import "bootstrap/dist/css/bootstrap.min.css"
import 'element-plus/dist/index.css'

import "bootstrap"
import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'
import {loadFonts} from './plugins/webfontloader'
import axios from 'axios';
import ElementPlus from 'element-plus'
import {createMetaManager} from 'vue-meta'

// Set the default prefix for all API requests.
// TODO: Somehow change this to automatically use the url the page is served from.
//axios.defaults.baseURL = 'http://192.168.42.100:8080/api/'

// Keep the cookies, that we get from the server and keep on sending them with every request.
//axios.defaults.withCredentials = true

loadFonts()

createApp(App)
    .use(router)
    .use(store)
    .use(vuetify)
    .use(ElementPlus)
    .use(createMetaManager())
    .mount('#app')
