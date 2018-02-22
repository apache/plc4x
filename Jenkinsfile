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
        PLC4X_BUILD = true
        JENKINS_PROFILE = 'jenkins-build'
        MVN_LOCAL_REPO_OPT = '-Dmaven.repo.local=.repository'
        MVN_TEST_FAIL_IGNORE = '-Dmaven.test.failure.ignore=true'
    }

    tools {
        maven 'Maven 3 (latest)'
        jdk 'JDK 1.8 (latest)'
    }

    options {
        timeout(time: 1, unit: 'HOURS')
        skipStagesAfterUnstable()
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
                sh 'mvn -P${JENKINS_PROFILE} ${MVN_TEST_FAIL_IGNORE} ${MVN_LOCAL_REPO_OPT} clean install'
            }
            post {
                always {
                    junit(testResults: '**/surefire-reports/*.xml', allowEmptyResults: true)
                    junit(testResults: '**/failsafe-reports/*.xml', allowEmptyResults: true)
                }
            }
        }

        stage('Build master') {
            when {
                branch 'master'
            }
            steps {
                echo 'Building'
                sh 'mvn -P${JENKINS_PROFILE} ${MVN_TEST_FAIL_IGNORE} clean install'
            }
            post {
                always {
                    junit(testResults: '**/surefire-reports/*.xml', allowEmptyResults: true)
                    junit(testResults: '**/failsafe-reports/*.xml', allowEmptyResults: true)
                }
            }
        }

        stage('Code Quality') {
            when {
                branch 'master'
            }
            steps {
                echo 'Building'
                sh 'mvn -P${JENKINS_PROFILE} sonar:sonar'
            }
        }

        stage('Deploy') {
            when {
                branch 'master'
            }
            steps {
                echo 'Building'
                sh 'mvn -P${JENKINS_PROFILE} deploy'
            }
        }

        stage('Build site') {
            when {
                branch 'master'
            }
            steps {
                echo 'Building'
                sh 'mvn -P${JENKINS_PROFILE} site:site'
            }
        }

        stage('Stage Site') {
            when {
                branch 'master'
            }
            steps {
                echo 'Staging Site'
                sh 'mvn -P${JENKINS_PROFILE} ${MVN_LOCAL_REPO_OPT} site:stage'
            }
        }

    }
}