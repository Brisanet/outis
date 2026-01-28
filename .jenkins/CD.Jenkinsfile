@Library('pipeline@v0.1.7') _

def tagInfo = [tag: null, shouldBuild: false, isRollback: false]
pipeline {
    agent any

    stages {
        stage('Checkout') {
            failFast true
            parallel {
                stage('Git') { steps { gitCheckout() } }
                stage('Environment') { steps { configChecker() } }
            }
        }

        stage('Code Versioning') {
            agent {
                docker {
                    image 'produtos.brisanet.net.br/v2/ci-tools:v2.1.13'
                    registryUrl 'https://produtos.brisanet.net.br/v2'
                    registryCredentialsId 'REGISTRY_CREDENTIALS'
                }
            }
            steps {
                script {
                    tagInfo = tagBuilder() 
                }
            }
        }

        stage('Build Image') {
            steps {
                script {
                    env.IMAGE_NAME = imageCraft(tagInfo)
                }
            }
        }

        stage('Vulnerability') {
            agent {
                docker {
                    image 'produtos.brisanet.net.br/v2/brisa-alpine:v1.1.1'
                    registryUrl 'https://produtos.brisanet.net.br/v2'
                    registryCredentialsId 'REGISTRY_CREDENTIALS'
                }
            }
            steps {
                vulnScan(env.IMAGE_NAME, tagInfo)
            }
        }

        stage('Sign Image') {
            agent {
                docker {
                    image 'produtos.brisanet.net.br/v2/consign:v1.0.3'
                    registryUrl 'https://produtos.brisanet.net.br/v2'
                    registryCredentialsId 'REGISTRY_CREDENTIALS'
                }
            }
            steps {
                imageSigner(env.IMAGE_NAME, tagInfo)
            }
        }

        stage('Upload Version') {
            steps {
                versionSync(tagInfo)
            }
        }
    }

    post {
        always {
            script {
                statusReporter(env.IMAGE_NAME, tagInfo)
            }
        }
    }
}