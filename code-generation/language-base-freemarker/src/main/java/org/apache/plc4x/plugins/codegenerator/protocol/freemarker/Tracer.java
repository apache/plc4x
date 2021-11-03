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
package org.apache.plc4x.plugins.codegenerator.protocol.freemarker;

import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Can be used to annotate generated code with traces out of the generator
 * <p>
 * uses c-style comments to inject traces. To customize sub this class and override separator(), prefix(), suffix()
 */
public class Tracer {

    protected static boolean ENABLED = "true".equalsIgnoreCase(System.getenv().get("PLC4X_TRACE_CODE_GEN"));

    protected final String currentTrace;

    protected Tracer(String currentTrace) {
        this.currentTrace = currentTrace;
    }

    /**
     * use this method to start a trace
     *
     * @param base usually the method name
     * @return a new trace containing the method name as base
     */
    public static Tracer start(String base) {
        return new Tracer(base);
    }

    /**
     * Returns a new Tracer with the appended sub
     *
     * @param sub usually a logical if name
     * @return a new trace with current trace + sub trace
     */
    public Tracer dive(String sub) {
        Tracer that = this;
        return new Tracer(currentTrace + separator() + sub) {
            @Override
            protected String prefix() {
                return that.prefix();
            }

            @Override
            protected String suffix() {
                return that.suffix();
            }
        };
    }

    /**
     * Can be used to remove traces from a traced string.
     *
     * @param somethingContainingTraces something containing traces
     * @return de-traced string
     */
    public String removeTraces(String somethingContainingTraces) {
        if (somethingContainingTraces == null) {
            return null;
        }
        return somethingContainingTraces.replaceAll(Pattern.quote(prefix()) + ".*" + Pattern.quote(suffix()), "");
    }

    /**
     * Can be used to extract traces from a traced string.
     *
     * @param somethingContainingTraces something containing traces
     * @return trace of something containing traces or "" if not traces available
     */
    public String extractTraces(String somethingContainingTraces) {
        Pattern pattern = Pattern.compile("(" + Pattern.quote(prefix()) + ".*" + Pattern.quote(suffix()) + ").*");
        Matcher matcher = pattern.matcher(somethingContainingTraces);
        if (!matcher.find()) {
            return "";
        }
        return matcher.group(1);
    }

    protected String separator() {
        return "/";
    }

    protected String prefix() {
        return "/*";
    }

    protected String suffix() {
        return "*/";
    }

    protected boolean isEnabled() {
        return ENABLED;
    }

    @Override
    public String toString() {
        if (!isEnabled()) {
            return "";
        }
        return prefix() + currentTrace + suffix();
    }
}
