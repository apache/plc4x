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

import groovy.text.SimpleTemplateEngine
import org.jsoup.Jsoup

import java.nio.file.Files
import java.nio.file.StandardCopyOption

// Make sure the cache directory exists in the used maven local repo
def localRepoBaseDir = session.getLocalRepository().getBasedir()
def cacheDir = new File(localRepoBaseDir, ".cache/bacnet-vendorids")
if (!cacheDir.exists()) {
    cacheDir.mkdirs()
}

// Check if a previous version exists and check if we need to re-download
// If the file is less than 24h old, we won't re-download in order to avoid
// being banned on the bacnet server.
def bacnetVendorHtm = new File(cacheDir, "BACnet Vendor IDs.htm")
def update = true
if (bacnetVendorHtm.exists()) {
    // If the last update was less than 24h before, don't update it again.
    if (bacnetVendorHtm.lastModified() > (new Date().getTime() - 86400000)) {
        update = false
    }
}

// If we need to update the vendor ids
if (update) {
    try {
        InputStream inputStream = new URL("http://www.bacnet.org/VendorID/BACnet%20Vendor%20IDs.htm").openStream()
        Files.copy(inputStream, bacnetVendorHtm.toPath(), StandardCopyOption.REPLACE_EXISTING)
        println "Successfully updated BACnet Vendor IDs.htm"
    } catch (Exception e) {
        println "Got an error updating BACnet Vendor IDs.htm. Intentionally not failing the build as we might just be offline: " + e.getMessage()
    }
} else {
    println "Skipped updating BACnet Vendor IDs.htm as it's fresh enough"
}

// Copy the knx-master-data to the current target directory
def targetDir = new File(project.getBasedir(), "target/downloads")
if (!targetDir.exists()) {
    targetDir.mkdirs()
}
def targetFile = new File(targetDir, "BACnet Vendor IDs.htm")
Files.copy(bacnetVendorHtm.toPath(), targetFile.toPath(), StandardCopyOption.REPLACE_EXISTING)

def reservedIds = ["555", "666", "777", "888", "911", "999", "1111"]

// Parse that now
def doc = Jsoup.parse(targetFile.text)

def table = doc.select("table").first();
def iterator = table.select("td").iterator();
def vendors = []
def foundOrganization = [:]
while (iterator.hasNext()) {
    def vendorId = iterator.next().text()
    def organization = iterator.next().text()
    // Reserved ones have a org so we abort here
    if (vendorId in reservedIds) {
        continue
    }
    def organizationSanitized = organization
        .replaceAll(/[^A-Za-z-0-1]/, "")
        .replaceAll(/([A-Z][a-z])/, /_$1/)
        .toUpperCase()
    // Remove starting underlines
        .replaceAll(/^_/, '')
    // Remove - in names
        .replaceAll(/-/, '')
    // Prefix digits with n
        .replaceAll(/^(\d)/, /N_$1/)
    def exitingOrganizationCount = foundOrganization[organizationSanitized]
    if (exitingOrganizationCount) {
        println "$organization found ${exitingOrganizationCount+1} times"
        organizationSanitized += "$exitingOrganizationCount"
        foundOrganization[organizationSanitized] = exitingOrganizationCount + 1
    } else {
        foundOrganization[organizationSanitized] = 1
    }
    def contactPerson = iterator.next().text()
    def address = iterator.next().text()
    println "Found vendorId:$vendorId organization:$organization contactPerson:$contactPerson address:$address"
    vendors << [vendorId: vendorId, organization: organization, organizationSanitized: organizationSanitized, contactPerson: contactPerson, address: address]
}

mspecTemplate = """
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
[enum uint 16 BACnetVendorId(uint 16 vendorId, string 8 vendorName)
<% for (item in values) { %> ['<%= item.vendorId %>' <%= item.organizationSanitized %> ['<%= item.vendorId %>', '"<%= item.organization %>"']]\\n <% } %>
  ['0xFFFF' UNKNOWN_VENDOR ['0xFFFF', '"Unknown"']]
]

[type BACnetVendorIdTagged(uint 8 tagNumber, TagClass tagClass)
    [simple   BACnetTagHeader
                        header                                                                               ]
    [validation    'header.tagClass == tagClass'    "tag class doesn't match"                                ]
    [validation    '(header.tagClass == TagClass.APPLICATION_TAGS && header.actualTagNumber == 2) || (header.actualTagNumber == tagNumber)'
                                                    "tagnumber doesn't match" shouldFail=false               ]
    [manual   BACnetVendorId
                    value
                        'STATIC_CALL("readEnumGeneric", readBuffer, header.actualLength, BACnetVendorId.UNKNOWN_VENDOR)'
                        'STATIC_CALL("writeEnumGeneric", writeBuffer, value)'
                        'header.actualLength * 8'                                                            ]
    [virtual  bit   isUnknownId
                        'value == BACnetVendorId.UNKNOWN_VENDOR'                                             ]
    //TODO: change to uint32 once cast is inserted
    [manual   uint 32
                    unknownId
                        'STATIC_CALL("readProprietaryEnumGeneric", readBuffer, header.actualLength, isUnknownId)'
                        'STATIC_CALL("writeProprietaryEnumGeneric", writeBuffer, unknownId, isUnknownId)'
                        '_value.isUnknownId?(header.actualLength * 8):0'                                       ]
]
"""
SimpleTemplateEngine templateEngine = new SimpleTemplateEngine()
def template = templateEngine.createTemplate(mspecTemplate).make([values: vendors])

def mspecTargetDir = new File(project.build.getOutputDirectory(), "/protocols/bacnetip")
if (!mspecTargetDir.exists()) {
    mspecTargetDir.mkdirs()
}
def vendorIdMspec = new File(mspecTargetDir, "bacnet-vendorids.mspec")

template.writeTo(vendorIdMspec.newWriter())
