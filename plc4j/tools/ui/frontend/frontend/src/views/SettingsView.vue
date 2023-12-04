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

import {defineComponent, ref} from "vue";
import axios from "axios";

function getLabelName(item:any):string {
  return item.value.label + "  (" + item.value.mountPoint + ")"
}

export default defineComponent({
  name: "Settings",
  data() {
    return {
      storagePath: ref(),
      enableRestServer: false,
      storagePathItems: [],
    }
  },
  async mounted() {
    // Load the storage path options available
    const items = await axios.get("configuration/options/storage-path")
    this.storagePathItems = items.data

    // Load the current configuration (or use a default, if no settings exist yet)
    try {
      const {data} = await axios.get('configuration')
      for(let storageOption of this.storagePathItems) {
        if(storageOption['uuid'] == data.storagePath) {
          this.storagePath = storageOption
          break
        }
      }
      //this.storagePath = data.storagePath
      this.enableRestServer = data.enableRestServer
    } catch (err) {
      this.storagePath = ""
      this.enableRestServer = true
    }

    console.log("done")
  },
  methods: {
    async submit() {
      await axios.post('configuration', {
        "id": 1,
        "storagePath": this.storagePath.uuid,
        "enableRestServer": this.enableRestServer
      })
    }
  }
})

</script>

<template>
    <form @submit.prevent="submit">
      <v-container>
        <v-row>
          <v-col>
            <v-select
                v-model="storagePath"
                :required="true"
                :items="storagePathItems"
                label="Storage Path"
                hint="This is where all time-series data is stored. Please note that when using flash-based storage, this storage will `age` under heavy load. Especially for the internal flash storage, this can reduce the lifetime of your CtrlX device. Using an external hard-drive is strongly encouraged."
                persistent-hint
            >
              <template v-slot:selection="{ item }">
                <v-list-item :title="`${item.value.label} (${item.value.mountPoint})`"/>
              </template>
              <template v-slot:item="{ item, props: { onClick } }" >
                <v-list-item :title="`${item.value.label} (${item.value.mountPoint})`" :disabled="!item.value.mounted" @click="onClick"/>
              </template>
            </v-select>
          </v-col>
        </v-row>

        <v-row>
          <v-col>
            <v-checkbox
                v-model="enableRestServer"
                hint="Enabled the TimechoDB rest-server. This is especially needed, when using the IoT Dashboard"
                persistent-hint
                label="Enable Rest Server"
                type="checkbox"
            ></v-checkbox>
          </v-col>
        </v-row>

        <v-row>
          <v-col>
            <v-btn
                class="me-4"
                type="submit"
            >
              Submit
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </form>
</template>

<style scoped>

</style>