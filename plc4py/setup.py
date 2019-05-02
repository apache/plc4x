"""

  Licensed to the Apache Software Foundation (ASF) under one or more
  contributor license agreements.  See the NOTICE file distributed with
  this work for additional information regarding copyright ownership.
  The ASF licenses this file to You under the Apache License, Version 2.0
  (the "License"); you may not use this file except in compliance with
  the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

"""

import xml.etree.ElementTree as ET
from io import open
from os import path

from setuptools import setup, find_packages

# Remember the location of this file.
here = path.abspath(path.dirname(__file__))

# Read the entire README.md and save it's content in the "long_description" variable.
with open(path.join(here, 'README.md'), encoding='utf-8') as f:
    longDescription = f.read()

# Register the maven namespace.
ns = {'mvn': 'http://maven.apache.org/POM/4.0.0'}
# Load the pom.xml and extract some of it's information from it.
mavenPomRoot = ET.parse(path.join(here, 'pom.xml'))
mvnArtifactId = mavenPomRoot.find("mvn:artifactId", ns).text
mvnVersion = mavenPomRoot.find("mvn:parent/mvn:version", ns).text
mvnDescription = mavenPomRoot.find("mvn:description", ns).text

# Cut off the "-SNAPSHOT"
if mvnVersion.endswith('-SNAPSHOT'):
    mvnVersion = mvnVersion[:-9]

print("ArtifactId", mvnArtifactId)
print("Version", mvnVersion)

setup(
    name=mvnArtifactId,
    version=mvnVersion,
    description=mvnDescription,
    long_description=longDescription,
    long_description_content_type='text/markdown',
    url='https://plc4x.apache.org',
    author_email='dev@plc4x.apache.org',

    packages=find_packages(exclude=['target', 'src/test']),

    python_requires='>=2.7, !=3.0.*, !=3.1.*, !=3.2.*, !=3.3.*, <4',

    data_files=[('lib', ['src/main/resources/lib/log4j2.xml'])],

    project_urls={
        'Bug Reports': 'https://issues.apache.org/jira/projects/PLC4X',
        'Source': 'https://gitbox.apache.org/repos/asf/plc4x.git',
    },
)