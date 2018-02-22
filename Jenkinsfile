#!groovy

/*
 *
 *  Licensed to the Apache Software Foundation (ASF) under one or more
 *  contributor license agreements.  See the NOTICE file distributed with
 *  this work for additional information regarding copyright ownership.
 *  The ASF licenses this file to You under the Apache License, Version 2.0
 *  (the "License"); you may not use this file except in compliance with
 *  the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
pipeline {
    agent {
        node {
            label 'ubuntu && !H33'
        }
    }

    environment {
        CC = 'clang'
        MVN_HOME = "${tool 'Maven 3 (latest)'}"
        JAVA_HOME = "${tool 'JDK 1.8 (latest)'}"
        PATH = "${env.JAVA_HOME}/bin:${env.PATH}"
    }

    stages {
        stage('Initialization') {
            steps {
                echo 'Building Branch: ' + env.BRANCH_NAME
            }
        }

        stage('Cleanup') {
            steps {
                echo 'Cleaning up the workspace'
                deleteDir()
            }
        }

        stage('Checkout') {
            steps {
                echo 'Checking out branch ' + env.BRANCH_NAME
                checkout scm
            }
        }

        stage('Build') {
            when {
                expression {
                    env.BRANCH_NAME != 'master'
                }
            }
            steps {
                echo 'Building'
                sh "${MVN_HOME}/bin/mvn -Pjenkins-build -Dmaven.repo.local=.repository clean install"
            }
            post {
                always {
                    junit(testResults: '**/surefire-reports/*.xml', allowEmptyResults: true)
                    junit(testResults: '**/failsafe-reports/*.xml', allowEmptyResults: true)
                }
            }
        }

        stage('Build') {
            when {
                branch 'master'
            }
            steps {
                echo 'Building'
                sh "${MVN_HOME}/bin/mvn -Pjenkins-build clean deploy sonar:sonar site:site"
            }
            post {
                always {
                    junit(testResults: '**/surefire-reports/*.xml', allowEmptyResults: true)
                    junit(testResults: '**/failsafe-reports/*.xml', allowEmptyResults: true)
                }
            }
        }

        stage('Stage Site') {
            when {
                branch 'master'
            }
            steps {
                echo 'Staging Site'
                sh "${MVN_HOME}/bin/mvn -Pjenkins-build ${mavenLocalRepo} site:stage"
            }
        }

    }
}