@Library('pipeline@v0.1.7') _

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

        stage('Unit Test') {
            agent {
                docker {
                    image 'produtos.brisanet.net.br/v2/ci-tools:v2.1.13'
                    registryUrl 'https://produtos.brisanet.net.br/v2'
                    registryCredentialsId 'REGISTRY_CREDENTIALS'
                }
            }
            steps {
                testScope()
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
                vulnScan()
            }
        }

        stage('Build Image') {
            agent {
                docker {
                    image 'produtos.brisanet.net.br/v2/ci-tools:v2.1.13'
                    registryUrl 'https://produtos.brisanet.net.br/v2'
                    registryCredentialsId 'REGISTRY_CREDENTIALS'
                }
            }
            steps {
                imageCraft()
            }
        }
    }

    post {
        always {
            statusReporter()
        }
    }
}