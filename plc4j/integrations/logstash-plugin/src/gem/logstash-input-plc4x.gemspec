Gem::Specification.new do |s|
  s.name            = '${project.artifactId}'
  s.version         = ${project.version}
  s.licenses        = ['Apache-2.0']
  s.summary         = "${project.name}"
  s.description     = "${project.description}"
  s.authors         = ['Apache Software Foundation']
  s.email           = 'dev@plc4x.apache.org'
  s.homepage        = "http://plc4x.apache.org/index.html"
  s.require_paths = ['lib', 'vendor/jar-dependencies']

  # Files
  s.files = Dir["lib/**/*","spec/**/*","*.gemspec","*.md","CONTRIBUTORS","Gemfile","LICENSE","NOTICE.TXT", "vendor/jar-dependencies/**/*.jar", "vendor/jar-dependencies/**/*.rb", "VERSION", "docs/**/*"]

  # Special flag to let us know this is actually a logstash plugin
  s.metadata = { 'logstash_plugin' => 'true', 'logstash_group' => 'input'}

  # Gem dependencies
  s.add_runtime_dependency "logstash-core-plugin-api", ">= 1.60", "<= 2.99"
  s.add_runtime_dependency 'jar-dependencies'

  s.add_development_dependency 'logstash-devutils'
end