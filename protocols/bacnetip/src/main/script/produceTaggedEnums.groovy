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
import org.apache.maven.project.MavenProject

project = (MavenProject) project
def bacnetEnumsFile = new File(project.basedir, "src/main/resources/protocols/bacnetip/bacnet-enums.mspec")

skippedEnums = ['BACnetRejectReason', 'BACnetAbortReason', 'BACnetPropertyStates']

foundEnums = []
enumPattern = ~/\[enum \w+ \d+ (\w+)/
bacnetEnumsFile.eachLine {
    def matcher = it =~ enumPattern
    if (matcher.find()) {
        def enumName = matcher[0][1]
        if (enumName in skippedEnums) {
            log.info("Skipping generation of tagged type for {}", enumName)
            return
        }
        log.info("Generating type for {}", enumName)
        foundEnums << [name: enumName]
    }
}

taggedEnumsTemplate = """
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
 
<% for (item in values) { %>
[type <%= item.name %>Tagged(uint 8 tagNumber, TagClass tagClass)
    [simple   BACnetTagHeader
                        header                                                                               ]
    [validation    'header.tagClass == tagClass'    "tag class doesn't match"                                ]
    [validation    '(header.tagClass == TagClass.APPLICATION_TAGS) || (header.actualTagNumber == tagNumber)'
                                                    "tagnumber doesn't match" shouldFail=false               ]
    [manual   <%= item.name %>
                    value
                        'STATIC_CALL("readEnumGeneric", readBuffer, header.actualLength, <%= item.name %>.VENDOR_PROPRIETARY_VALUE)'
                        'STATIC_CALL("writeEnumGeneric", writeBuffer, value)'
                        '_value.isProprietary?0:(header.actualLength * 8)'                                   ]
    [virtual  bit   isProprietary
                        'value == <%= item.name %>.VENDOR_PROPRIETARY_VALUE'                                 ]
    [manual   uint 32
                    proprietaryValue
                        'STATIC_CALL("readProprietaryEnumGeneric", readBuffer, header.actualLength, isProprietary)'
                        'STATIC_CALL("writeProprietaryEnumGeneric", writeBuffer, proprietaryValue, isProprietary)'
                        '_value.isProprietary?(header.actualLength * 8):0'                                   ]
]
<% } %>
"""
SimpleTemplateEngine templateEngine = new SimpleTemplateEngine()
def template = templateEngine.createTemplate(taggedEnumsTemplate).make([values: foundEnums])

def mspecTargetDir = new File(project.build.getOutputDirectory(), "/protocols/bacnetip")
if (!mspecTargetDir.exists()) {
    mspecTargetDir.mkdirs()
}
def vendorIdMspec = new File(mspecTargetDir, "bacnet-enums-tagged.mspec")

template.writeTo(vendorIdMspec.newWriter())
