/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.language;

import freemarker.template.Configuration;
import freemarker.template.Template;

import java.io.IOException;
import java.util.List;

public interface LanguageTemplate {

    /**
     * The name of the template is what the plugin will use to select the correct language module.
     *
     * @return the name of the template.
     */
    String getName();

    /**
     * As some times it might be necessary to generate more than one file from one type, this method
     * returns a list of templates which will all be applied against a given type.
     *
     * @return list of templates to use for generating output.
     * @throws IOException if something goes wrong while loading the templates.
     */
    List<Template> getTemplates(Configuration freemarkerConfiguration) throws IOException;

    /**
     * Helper that will be used by the templates (to know which type to generate).
     *
     * @return instance of a helper for the given language.
     */
    LanguageTemplateHelper getHelper();

}
