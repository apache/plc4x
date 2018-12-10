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
            label 'plc4x'
        }
    }

    environment {
        // It seems the login the jenkins slave uses, doesn't pick up the environment changes,
        // so we have to try to manually add theme here.
        MAVEN_HOME = '/opt/maven'
        PATH = "${MAVEN_HOME}/bin:${env.PATH}"

        PLC4X_BUILD_ON_JENKINS = true
        JENKINS_PROFILE = 'jenkins-build'
        // On non develop build we don't want to pollute the global m2 repo
        MVN_LOCAL_REPO_OPT = '-Dmaven.repo.local=.repository'
        // Testfails will be handled by the jenkins junit steps and mark the build as unstable.
        MVN_TEST_FAIL_IGNORE = '-Dmaven.test.failure.ignore=true'
        // Make JQAssistant run with Neo4j 3.1.3 instead of 2.x
        JQASSISTANT_NEO4J_VERSION = '-Djqassistant.neo4jVersion=3'
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
                echo 'Using PATH = ' + env.PATH
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
                    env.BRANCH_NAME != 'develop'
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

        stage('Build develop') {
            when {
                branch 'develop'
            }
            steps {
                echo 'Building'
                // Make sure the directory is wiped.
                sh 'rm -rf ./local-snapshots-dir'

                // We'll deploy to a relative directory so we can save
                // that and deploy in a later step on a different node
                sh 'mvn -P${JENKINS_PROFILE} ${MVN_TEST_FAIL_IGNORE} ${JQASSISTANT_NEO4J_VERSION} -DaltDeploymentRepository=snapshot-repo::default::file:./local-snapshots-dir clean deploy'

                // Stash the build results so we can deploy them on another node
                stash name: 'plc4x-build-snapshots', includes: 'local-snapshots-dir/**'
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
                branch 'develop'
            }
            steps {
                echo 'Checking Code Quality'
                withSonarQubeEnv('ASF Sonar Analysis') {
                    sh 'mvn -P${JENKINS_PROFILE} sonar:sonar'
                }
            }
        }

        stage('Deploy') {
            when {
                branch 'develop'
            }
            // Only the official build nodes have the credentials to deploy setup.
            agent {
                node {
                    label 'ubuntu'
                }
            }
            steps {
                echo 'Deploying'
                // Clean up the snapshots directory.
                dir("local-snapshots-dir/") {
                    deleteDir()
                }

                // Unstash the previously stashed build results.
                unstash name: 'plc4x-build-snapshots'

                // Deploy the artifacts using the wagon-maven-plugin.
                sh 'mvn -f jenkins.pom -X -P deploy-snapshots wagon:upload'
            }
        }

        stage('Build site') {
            when {
                branch 'develop'
            }
            steps {
                echo 'Building Site'
                sh 'mvn -P${JENKINS_PROFILE} site'
            }
        }

        stage('Stage site') {
            when {
                branch 'develop'
            }
            steps {
                echo 'Staging Site'
                sh 'mvn -P${JENKINS_PROFILE} site:stage'
                // Stash the generated site so we can publish it on the 'git-website' node.
                stash includes: 'target/staging/**/*', name: 'plc4x-site'
            }
        }

        stage('Deploy site') {
            when {
                branch 'develop'
            }
            // Only the nodes labeled 'git-websites' have the credentials to commit to the.
            agent {
                node {
                    label 'git-websites'
                }
            }
            steps {
                echo 'Deploying Site'
                // Unstash the previously stashed site.
                unstash 'plc4x-site'
                // Publish the site with the scm-publish plugin.
                sh 'mvn -P${JENKINS_PROFILE} scm-publish:publish-scm'
            }
        }
    }

    // Send out notifications on unsuccessfull builds.
    post {
         failure {
            emailext (
                subject: "[BUILD-FAILURE]: Job '${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]'",
                body: """<p>BUILD-FAILURE: Job '${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]':</p><p>Check console output at "<a href="${env.BUILD_URL}">${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]</a>"</p>""",
                to: "dev@plc4x.apache.org",
                recipientProviders: [[$class: 'DevelopersRecipientProvider']]
            )
        }
        unstable {
            emailext (
                subject: "[BUILD-UNSTABLE]: Job '${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]'",
                body: """<p>BUILD-UNSTABLE: Job '${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]':</p><p>Check console output at "<a href="${env.BUILD_URL}">${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]</a>"</p>""",
                to: "dev@plc4x.apache.org",
                recipientProviders: [[$class: 'DevelopersRecipientProvider']]
            )
        }
    }

}