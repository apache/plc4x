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
<script setup lang="ts">
</script>

<template>
  <v-navigation-drawer :width="320">
    <v-list-item>
      <v-img :src="require('./../assets/plc4x-logo.svg')"/>
    </v-list-item>
    <v-list-item title="Apache PLC4X"></v-list-item>
    <v-divider></v-divider>
    <v-list-item title="Dashboard" prepend-icon="mdi-monitor-dashboard" :to="'/dashboard'"></v-list-item>
    <v-list-item title="Jobs" prepend-icon="mdi-format-list-checks" :to="'/jobs'"></v-list-item>
    <v-divider class="h-25"></v-divider>
    <v-list-item title="Settings" prepend-icon="mdi-cog" :to="'/settings'"></v-list-item>
    <v-list-item title="About" prepend-icon="mdi-information-outline" :to="'/about'"></v-list-item>
  </v-navigation-drawer>
</template>
