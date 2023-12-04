<!--
  - Licensed to the Apache Software Foundation (ASF) under one
  - or more contributor license agreements.  See the NOTICE file
  - distributed with this work for additional information
  - regarding copyright ownership.  The ASF licenses this file
  - to you under the Apache License, Version 2.0 (the
  - "License"); you may not use this file except in compliance
  - with the License.  You may obtain a copy of the License at
  -
  -   http://www.apache.org/licenses/LICENSE-2.0
  -
  - Unless required by applicable law or agreed to in writing,
  - software distributed under the License is distributed on an
  - "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  - KIND, either express or implied.  See the License for the
  - specific language governing permissions and limitations
  - under the License.
  -->
<script lang="ts">
import { defineComponent } from 'vue'
import { useMeta } from 'vue-meta'
import {TreeItem} from "@/models/treeItem";

export default defineComponent({
    name: 'App',
    metaInfo: {
        title: 'Default Title',
    },
    data: () => ({
        selectedNavigationTab: null,
        byProtocolItems: [
            {
                id: 1,
                name: 'Applications :',
                children: [
                    { id: 2, name: 'Calendar : app' },
                    { id: 3, name: 'Chrome : app' },
                    { id: 4, name: 'Webstorm : app' },
                ],
            },
            {
                id: 5,
                name: 'Documents :',
                children: [
                    {
                        id: 6,
                        name: 'vuetify :',
                        children: [
                            {
                                id: 7,
                                name: 'src :',
                                children: [
                                    { id: 8, name: 'index : ts' },
                                    { id: 9, name: 'bootstrap : ts' },
                                ],
                            },
                        ],
                    },
                    {
                        id: 10,
                        name: 'material2 :',
                        children: [
                            {
                                id: 11,
                                name: 'src :',
                                children: [
                                    { id: 12, name: 'v-btn : ts' },
                                    { id: 13, name: 'v-card : ts' },
                                    { id: 14, name: 'v-window : ts' },
                                ],
                            },
                        ],
                    },
                ],
            },
            {
                id: 15,
                name: 'Downloads :',
                children: [
                    { id: 16, name: 'October : pdf' },
                    { id: 17, name: 'November : pdf' },
                    { id: 18, name: 'Tutorial : html' },
                ],
            },
            {
                id: 19,
                name: 'Videos :',
                children: [
                    {
                        id: 20,
                        name: 'Tutorials :',
                        children: [
                            { id: 21, name: 'Basic layouts : mp4' },
                            { id: 22, name: 'Advanced techniques : mp4' },
                            { id: 23, name: 'All about app : dir' },
                        ],
                    },
                    { id: 24, name: 'Intro : mov' },
                    { id: 25, name: 'Conference introduction : avi' },
                ],
            },        ],
        byAddressItems: [
            {
                id: 1,
                name: 'Applications :',
                children: [
                    { id: 2, name: 'Calendar : app' },
                    { id: 3, name: 'Chrome : app' },
                    { id: 4, name: 'Webstorm : app' },
                ],
            },
            {
                id: 5,
                name: 'Documents :',
                children: [
                    {
                        id: 6,
                        name: 'vuetify :',
                        children: [
                            {
                                id: 7,
                                name: 'src :',
                                children: [
                                    { id: 8, name: 'index : ts' },
                                    { id: 9, name: 'bootstrap : ts' },
                                ],
                            },
                        ],
                    },
                    {
                        id: 10,
                        name: 'material2 :',
                        children: [
                            {
                                id: 11,
                                name: 'src :',
                                children: [
                                    { id: 12, name: 'v-btn : ts' },
                                    { id: 13, name: 'v-card : ts' },
                                    { id: 14, name: 'v-window : ts' },
                                ],
                            },
                        ],
                    },
                ],
            },
            {
                id: 15,
                name: 'Downloads :',
                children: [
                    { id: 16, name: 'October : pdf' },
                    { id: 17, name: 'November : pdf' },
                    { id: 18, name: 'Tutorial : html' },
                ],
            },
            {
                id: 19,
                name: 'Videos :',
                children: [
                    {
                        id: 20,
                        name: 'Tutorials :',
                        children: [
                            { id: 21, name: 'Basic layouts : mp4' },
                            { id: 22, name: 'Advanced techniques : mp4' },
                            { id: 23, name: 'All about app : dir' },
                        ],
                    },
                    { id: 24, name: 'Intro : mov' },
                    { id: 25, name: 'Conference introduction : avi' },
                ],
            },
        ],
        selectedTab: null,
    }),
    methods: {
        removeTab(index: string) {
            console.log(index)
        },
    }
})
</script>

<template>
  <v-app>
      <v-layout class="rounded rounded-md">
          <v-app-bar>
              <v-img :src="require('./assets/plc4x-logo.svg')" max-width="260" class="ma-8"/>
              <v-app-bar-title>Application bar</v-app-bar-title>
          </v-app-bar>

          <v-navigation-drawer :width="400">
              <v-container fluid style="height: 100vh;">
                  <v-tabs v-model="selectedNavigationTab" class="border-b-sm" show-arrows>
                      <v-tab value="protocol">By Protocol</v-tab>
                      <v-tab value="address">By Address</v-tab>
                  </v-tabs>
                  <v-card-text>
                      <v-window v-model="selectedNavigationTab">
                          <v-window-item value="protocol">
                              <v-treeview :items="byProtocolItems"></v-treeview>
                          </v-window-item>
                          <v-window-item value="address">
                              <v-treeview :items="byAddressItems"></v-treeview>
                          </v-window-item>
                      </v-window>
                  </v-card-text>
              </v-container>
          </v-navigation-drawer>

          <v-main class="d-flex align-center justify-center">
              <v-container fluid style="height: 100vh;">
                  <v-tabs v-model="selectedTab" class="border-b-sm" show-arrows>
                      <v-tab value="s7://192.168.23.30">s7://192.168.23.30
                          <v-btn icon @click="removeTab('lalal')" class="ml-2"
                          ><v-icon x-small>mdi-close</v-icon></v-btn
                          >
                      </v-tab>
                      <v-tab value="ads://192.168.23.20">ads://192.168.23.20
                          <v-btn icon @click="removeTab('huiii')" class="ml-2"
                          ><v-icon x-small>mdi-close</v-icon></v-btn
                          >
                      </v-tab>
                  </v-tabs>
                  <v-card-text>
                      <v-window v-model="selectedTab">
                          <v-window-item value="s7://192.168.23.30">
                              One
                          </v-window-item>
                          <v-window-item value="ads://192.168.23.20">
                              Two
                          </v-window-item>
                      </v-window>
                  </v-card-text>
              </v-container>
              <!--router-view/-->
          </v-main>
      </v-layout>
  </v-app>
</template>

<style>
.bi {
  display: inline-block;
  width: 1rem;
  height: 1rem;
}

/*
 * Sidebar
 */

@media (min-width: 768px) {
  .sidebar .offcanvas-lg {
    position: -webkit-sticky;
    position: sticky;
    top: 48px;
  }
  .navbar-search {
    display: block;
  }
}

.sidebar .nav-link {
  font-size: .875rem;
  font-weight: 500;
}

.sidebar .nav-link.active {
  color: #2470dc;
}

.sidebar-heading {
  font-size: .75rem;
}

/*
 * Navbar
 */

.navbar-brand {
  padding-top: .75rem;
  padding-bottom: .75rem;
  background-color: rgba(0, 0, 0, .25);
  box-shadow: inset -1px 0 0 rgba(0, 0, 0, .25);
}

.navbar .form-control {
  padding: .75rem 1rem;
}
</style>