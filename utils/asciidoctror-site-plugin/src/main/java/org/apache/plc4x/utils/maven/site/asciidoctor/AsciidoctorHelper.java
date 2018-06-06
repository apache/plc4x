/*
 * Licensed under the Apache License, Version 2.0 (the &quot;License&quot;);
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an &quot;AS IS&quot; BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.plc4x.utils.maven.site.asciidoctor;

import org.asciidoctor.Attributes;
import org.asciidoctor.AttributesBuilder;

import java.util.Map;

/**
 * Utility class for re-usable logic.
 */
public class AsciidoctorHelper {

    /**
     * Adds attributes from a {@link Map} into a {@link AttributesBuilder} taking care of Maven's XML parsing special
     * cases like toggles, nulls, etc.
     */
    public static void addAttributes(final Map<String, Object> attributes, AttributesBuilder attributesBuilder) {
        // TODO Figure out how to reliably set other values (like boolean values, dates, times, etc)
        for (Map.Entry<String, Object> attributeEntry : attributes.entrySet()) {
            addAttribute(attributeEntry.getKey(), attributeEntry.getValue(), attributesBuilder);
        }
    }

    /**
     * Adds an attribute into a {@link AttributesBuilder} taking care of Maven's XML parsing special cases like
     * toggles toggles, nulls, etc.
     */
    public static void addAttribute(String attribute, Object value, AttributesBuilder attributesBuilder) {
        // NOTE Maven interprets an empty value as null, so we need to explicitly convert it to empty string (see #36)
        // NOTE In Asciidoctor, an empty string represents a true value
        if (value == null || "true".equals(value)) {
            attributesBuilder.attribute(attribute, "");
        }
        // NOTE a value of false is effectively the same as a null value, so recommend the use of the string "false"
        else if ("false".equals(value)) {
            attributesBuilder.attribute(attribute, null);
        }
        // NOTE Maven can't assign a Boolean value from the XML-based configuration, but a client may
        else if (value instanceof Boolean) {
            attributesBuilder.attribute(attribute, Attributes.toAsciidoctorFlag((Boolean) value));
        } else {
            // Can't do anything about dates and times because all that logic is private in Attributes
            attributesBuilder.attribute(attribute, value);
        }
    }

}
