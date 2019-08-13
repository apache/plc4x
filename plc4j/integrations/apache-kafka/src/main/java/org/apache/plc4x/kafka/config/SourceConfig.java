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
package org.apache.plc4x.kafka.config;

import com.google.gson.*;
import org.apache.plc4x.kafka.exceptions.ConfigurationException;

import java.util.List;
import java.util.Map;

public class SourceConfig {

    private final List<Source> sources;
    private final List<Job> jobs;

    public static SourceConfig fromPropertyMap(Map<String, String> properties) throws ConfigurationException {
        try {
            JsonObject jsonConfig = propertiesToJson(properties);
            return new Gson().fromJson(jsonConfig, SourceConfig.class);
        } catch (Exception e) {
            throw new ConfigurationException("Error configuring.", e);
        }
    }

    private static JsonObject propertiesToJson(Map<String, String> properties) {
        JsonObject config = new JsonObject();

        String defaultTopic = properties.getOrDefault("defaults.topic", null);
        // Generally create the JSON tree structure.
        for (Map.Entry<String, String> entry : properties.entrySet()) {
            String key = entry.getKey();
            if("defaults.topic".equals(key)) {
                continue;
            }
            JsonObject context = config;
            String[] segments = key.split("\\.");
            for (int i = 0; i < segments.length; i++) {
                String segment = segments[i];
                // If this is the leaf segment, set the value.
                if(i == segments.length - 1) {
                    context.add(segment, new JsonPrimitive(properties.get(key)));
                }
                // if it's not a leaf segment, ensure the structure exists.
                else {
                    if (!context.has(segment)) {
                        context.add(segment, new JsonObject());
                    }
                    context = (JsonObject) context.get(segment);
                }
            }
        }

        // Transform the structure for the sources into a real JSON (Jackson serializable) form.
        JsonObject oldSources = config.getAsJsonObject("sources");
        JsonArray sources = new JsonArray();
        for (String sourceName : oldSources.keySet()) {
            JsonObject oldSource = oldSources.getAsJsonObject(sourceName);

            // Create a new JSON object with the correct structure.
            JsonObject source = new JsonObject();
            source.add("name", new JsonPrimitive(sourceName));
            String sourceTopic = oldSource.has("topic") ? oldSource.get("topic").getAsString() : defaultTopic;
            // Copy all the other properties over.
            for (String sourceProperty : oldSource.keySet()) {
                if("topic".equals(sourceProperty)) {
                    // Filter out the topic setting as this is set on jobReference level.
                } else if("jobReferences".equals(sourceProperty)) {
                    JsonObject oldJobReferences = oldSource.getAsJsonObject(sourceProperty);
                    JsonArray jobReferences = new JsonArray();
                    for (String jobReferenceName : oldJobReferences.keySet()) {
                        JsonObject oldJobReference = oldJobReferences.getAsJsonObject(jobReferenceName);
                        JsonObject jobReference = new JsonObject();
                        jobReference.add("name", new JsonPrimitive(jobReferenceName));
                        for (String jobReferenceProperty : oldJobReference.keySet()) {
                            jobReference.add(jobReferenceProperty, oldJobReference.get(jobReferenceProperty));
                        }
                        if(!jobReference.has("topic")) {
                            jobReference.add("topic", new JsonPrimitive(sourceTopic));
                        }
                        jobReferences.add(jobReference);
                    }
                    source.add(sourceProperty, jobReferences);
                } else {
                    source.add(sourceProperty, oldSource.get(sourceProperty));
                }
            }

            // Add the new source to the existing sources.
            sources.add(source);
        }
        config.remove("sources");
        config.add("sources", sources);

        // Transform the structure of the jobs into a real JSON form.
        JsonObject oldJobs = config.getAsJsonObject("jobs");
        JsonArray jobs = new JsonArray();
        for (String jobName : oldJobs.keySet()) {
            JsonObject oldJob = oldJobs.getAsJsonObject(jobName);

            // Create a new JSON object with the correct structure.
            JsonObject job = new JsonObject();
            job.add("name", new JsonPrimitive(jobName));
            // Copy all the other properties over.
            for (String jobProperty : oldJob.keySet()) {
                job.add(jobProperty, oldJob.get(jobProperty));
            }

            // Add the new source to the existing sources.
            jobs.add(job);
        }
        config.remove("jobs");
        config.add("jobs", jobs);

        return config;
    }

    public SourceConfig(List<Source> sources, List<Job> jobs) {
        this.sources = sources;
        this.jobs = jobs;
    }

    public List<Source> getSources() {
        return sources;
    }

    public List<Job> getJobs() {
        return jobs;
    }

}
