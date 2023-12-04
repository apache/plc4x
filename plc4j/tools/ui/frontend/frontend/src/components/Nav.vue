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
import {defineComponent} from "vue";

export default defineComponent({
  name: "Nav",
  props: [
      "title"
  ],
})

</script>

<template>
  <header class="navbar sticky-top bg-dark flex-md-nowrap p-0 shadow" data-bs-theme="dark">
    <v-app-bar title="{{ title }}">
      <v-img :src="require('./../assets/plc4x-logo.svg')" max-width="260" class="ma-8"/>
      <!--v-toolbar-title>{{ title }}</v-toolbar-title-->
    </v-app-bar>
  </header>
</template>
