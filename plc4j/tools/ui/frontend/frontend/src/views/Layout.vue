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
import Menu from "@/components/Menu.vue";
import Nav from "@/components/Nav.vue";
import {defineComponent} from "vue";

export default defineComponent({
  name: "Layout",
  components: {Menu, Nav},
  data() {
    return {
      title: 'lalala'
    }
  },
  created() {
    this.title = 'Apache PLC4X: ' + this.$route.meta?.title as string ?? 'unknown'
  },
  updated() {
    this.title = 'Apache PLC4X: ' + this.$route.meta?.title as string ?? 'unknown'
  }
})
</script>

<template>
    <v-layout class="rounded rounded-md">
        <v-app-bar title="Application bar"></v-app-bar>

        <v-navigation-drawer>
            <v-list>
                <v-list-item title="Navigation drawer"></v-list-item>
            </v-list>
        </v-navigation-drawer>

        <v-main class="d-flex align-center justify-center" style="min-height: 300px;">
            <router-view/>
        </v-main>
    </v-layout>
</template>

<style scoped>

</style>