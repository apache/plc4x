# ----------------------------------------------------------------------------
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
# ----------------------------------------------------------------------------

Gem::Specification.new do |s|
  s.name            = 'logstash-input-plc4x'
  s.version         = '${current-full-version}'
  s.licenses        = ['Apache-2.0']
  s.summary         = '${project.name}'
  s.description     = '${project.description}'
  s.authors         = ['Apache Software Foundation']
  s.email           = ['dev@plc4x.apache.org']
  s.homepage        = 'http://plc4x.apache.org/index.html'
  s.require_paths = ['lib', 'vendor/jar-dependencies']

  # Files
  s.files = Dir["lib/**/*","*.gemspec","*.md","CONTRIBUTORS","Gemfile","LICENSE","NOTICE.TXT", "vendor/jar-dependencies/**/*.jar", "vendor/jar-dependencies/**/*.rb", "VERSION", "docs/**/*"]

  # Special flag to let us know this is actually a logstash plugin
  s.metadata = { 'logstash_plugin' => 'true', 'logstash_group' => 'input', 'java_plugin' => 'true'}

  # Gem dependencies

  s.add_runtime_dependency 'logstash-core-plugin-api', '>= 1.60', '<= 2.99'
  s.add_runtime_dependency 'jar-dependencies', '~> 0'
  s.add_development_dependency 'logstash-devutils', '~> 0'
end