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
        PLC4X_BUILD_ON_JENKINS = true
        JENKINS_PROFILE = 'jenkins-build'
        // On non master build we don't want to pollute the global m2 repo
        MVN_LOCAL_REPO_OPT = '-Dmaven.repo.local=.repository'
        // Testfails will be handled by the jenkins junit steps and mark the build as unstable.
        MVN_TEST_FAIL_IGNORE = '-Dmaven.test.failure.ignore=true'
    }

    tools {
        maven 'Maven 3 (latest)'
        jdk 'JDK 1.8 (latest)'
    }

    options {
        timeout(time: 1, unit: 'HOURS')
        // When we have test-fails e.g. we don't need to run the remaining steps
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

        stage('Stage site') {
            when {
                branch 'master'
            }
            steps {
                echo 'Staging Site'
                sh 'mvn -P${JENKINS_PROFILE} site:stage'
            }
        }

        stage('Deploy site') {
            when {
                branch 'master'
            }
            agent {
                node {
                    label 'git-websites'
                }
            }
            steps {
                echo 'Deploy Site'
                // We need to regenerate the site for deploy as we switch the node. We could save time by stash/unstash.
                sh 'mvn -P${JENKINS_PROFILE} site:site'
                sh 'mvn -P${JENKINS_PROFILE} site:deploy'
            }
        }

    }
}