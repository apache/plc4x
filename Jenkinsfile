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
        // Test failures will be handled by the jenkins junit steps and mark the build as unstable.
        MVN_TEST_FAIL_IGNORE = '-Dmaven.test.failure.ignore=true'

        SONARCLOUD_PARAMS = "-Dsonar.host.url=https://sonarcloud.io -Dsonar.organization=apache -Dsonar.projectKey=apache_plc4x -Dsonar.branch.name=develop"
    }

    tools {
        maven 'maven_3_latest'
        jdk 'jdk_11_latest'
    }

    options {
        // Kill this job after one hour.
        timeout(time: 24, unit: 'HOURS')
        // When we have test-fails e.g. we don't need to run the remaining steps
        skipStagesAfterUnstable()
        buildDiscarder(logRotator(numToKeepStr: '5', artifactNumToKeepStr: '3'))
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
                //sh 'mvn -P${JENKINS_PROFILE},skip-prerequisite-check,development,with-sandbox,with-c,with-cpp,with-boost,with-dotnet,with-python,with-proxies,with-logstash ${MVN_TEST_FAIL_IGNORE} ${MVN_LOCAL_REPO_OPT} clean install'
                //sh 'mvn -P${JENKINS_PROFILE},skip-prerequisite-check,development,with-sandbox,with-logstash,with-go ${MVN_TEST_FAIL_IGNORE} ${MVN_LOCAL_REPO_OPT} clean install'
                sh 'mvn -P${JENKINS_PROFILE},skip-prerequisite-check,with-sandbox,with-logstash,with-go ${MVN_TEST_FAIL_IGNORE} ${MVN_LOCAL_REPO_OPT} clean install'
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
                // Clean up the snapshots directory.
                dir("local-snapshots-dir/") {
                    deleteDir()
                }

                // We'll deploy to a relative directory so we can save
                // that and deploy in a later step on a different node
                //sh 'mvn -U -P${JENKINS_PROFILE},skip-prerequisite-check,development,with-sandbox,with-c,with-cpp,with-boost,with-dotnet,with-python,with-proxies,with-logstash ${MVN_TEST_FAIL_IGNORE} ${JQASSISTANT_NEO4J_VERSION} -DaltDeploymentRepository=snapshot-repo::default::file:./local-snapshots-dir clean deploy'
                //sh 'mvn -U -P${JENKINS_PROFILE},skip-prerequisite-check,development,with-sandbox,with-logstash,with-go ${MVN_TEST_FAIL_IGNORE} ${JQASSISTANT_NEO4J_VERSION} -DaltDeploymentRepository=snapshot-repo::default::file:./local-snapshots-dir clean deploy'
                sh 'mvn -U -P${JENKINS_PROFILE},skip-prerequisite-check,with-sandbox,with-logstash,with-go ${MVN_TEST_FAIL_IGNORE} ${JQASSISTANT_NEO4J_VERSION} -DaltDeploymentRepository=snapshot-repo::default::file:./local-snapshots-dir clean deploy'

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
                echo 'Checking Code Quality on SonarCloud'
                withCredentials([string(credentialsId: 'chris-sonarcloud-token', variable: 'SONAR_TOKEN')]) {
                    //sh 'mvn -P${JENKINS_PROFILE},skip-prerequisite-check,with-python,with-proxies,with-sandbox,with-logstash sonar:sonar ${SONARCLOUD_PARAMS} -Dsonar.login=${SONAR_TOKEN}'
                    sh 'mvn -P${JENKINS_PROFILE},skip-prerequisite-check,with-sandbox,with-logstash sonar:sonar ${SONARCLOUD_PARAMS} -Dsonar.login=${SONAR_TOKEN}'
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
                    // TODO: Disabled as H50 seems to be the only node with this label and it's currently offline for quite some time.
                    // label 'nexus-deploy'
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

                // Clean up the snapshots directory (freeing up more space after deploying).
                dir("local-snapshots-dir/") {
                    deleteDir()
                }
            }
        }

        stage('Build site') {
            when {
                branch 'develop'
            }
            steps {
                echo 'Building Site'
                //sh 'mvn -P${JENKINS_PROFILE},skip-prerequisite-check,with-proxies,with-logstash site'
                sh 'mvn -Djava.version=1.8 -P${JENKINS_PROFILE},skip-prerequisite-check,with-logstash site -X -pl .'
            }
        }

        stage('Stage site') {
            when {
                branch 'develop'
            }
            steps {
                echo 'Staging Site'
                // Build a directory containing the aggregated website.
                //sh 'mvn -P${JENKINS_PROFILE},skip-prerequisite-check,with-proxies,with-logstash site:stage'
                sh 'mvn -P${JENKINS_PROFILE},skip-prerequisite-check,with-logstash site:stage -pl .'
                // Make sure the script is executable.
                sh 'chmod +x tools/clean-site.sh'
                // Remove some redundant resources, which shouldn't be required.
                //sh 'tools/clean-site.sh'
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
                // Clean up the site directory.
                dir("target/staging") {
                    deleteDir()
                }

                // Unstash the previously stashed site.
                unstash 'plc4x-site'
                // Publish the site with the scm-publish plugin.
                sh 'mvn -f jenkins.pom -X -P deploy-site scm-publish:publish-scm'

                // Clean up the snapshots directory (freeing up more space after deploying).
                dir("target/staging") {
                    deleteDir()
                }
            }
        }
    }

    // Send out notifications on unsuccessful builds.
    post {
        // If this build failed, send an email to the list.
        failure {
            script {
                if(env.BRANCH_NAME == "develop") {
                    emailext(
                        subject: "[BUILD-FAILURE]: Job '${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]'",
                        body: """
BUILD-FAILURE: Job '${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]':

Check console output at "<a href="${env.BUILD_URL}">${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]</a>"
""",
                        to: "dev@plc4x.apache.org",
                        recipientProviders: [[$class: 'DevelopersRecipientProvider']]
                    )
                }
            }
        }

        // If this build didn't fail, but there were failing tests, send an email to the list.
        unstable {
            script {
                if(env.BRANCH_NAME == "develop") {
                    emailext(
                        subject: "[BUILD-UNSTABLE]: Job '${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]'",
                        body: """
BUILD-UNSTABLE: Job '${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]':

Check console output at "<a href="${env.BUILD_URL}">${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]</a>"
""",
                        to: "dev@plc4x.apache.org",
                        recipientProviders: [[$class: 'DevelopersRecipientProvider']]
                    )
                }
            }
        }

        // Send an email, if the last build was not successful and this one is.
        success {
            // Cleanup the build directory if the build was successful
            // (in this cae we probably don't have to do any post-build analysis)
            deleteDir()
            script {
                if ((env.BRANCH_NAME == "develop") && (currentBuild.previousBuild != null) && (currentBuild.previousBuild.result != 'SUCCESS')) {
                    emailext (
                        subject: "[BUILD-STABLE]: Job '${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]'",
                        body: """
BUILD-STABLE: Job '${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]':

Is back to normal.
""",
                        to: "dev@plc4x.apache.org",
                        recipientProviders: [[$class: 'DevelopersRecipientProvider']]
                    )
                }
            }
        }

        always {
            script {
                if(env.BRANCH_NAME == "release") {
                    // Double check if something was really changed as sometimes the
                    // build just runs without any changes.
                    if(currentBuild.changeSets.size() > 0) {
                        emailext(
                            subject: "[COMMIT-TO-RELEASE]: A commit to the release branch was made'",
                            body: """
COMMIT-TO-RELEASE: A commit to the release branch was made:

Check console output at "<a href="${env.BUILD_URL}">${env.JOB_NAME} [${env.BRANCH_NAME}] [${env.BUILD_NUMBER}]</a>"
""",
                            to: "dev@plc4x.apache.org",
                            recipientProviders: [[$class: 'DevelopersRecipientProvider']]
                        )
                    }
                }
            }
        }
    }

}
