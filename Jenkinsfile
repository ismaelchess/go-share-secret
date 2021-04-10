pipeline{
    agent any

    environment {
        XDG_CACHE_HOME="/tmp/.cache"
    }

    stages {
        stage('Pre Test') {
            steps {
                withCredentials([string(credentialsId: 'github-token', variable: 'NETRC')]) {
                    sh '''
                        make get
                    '''
                }
            }
        }
        stage('Build') {
            steps {
                withCredentials([string(credentialsId: 'github-token', variable: 'NETRC')]) {
                    sh '''
                        make build
                    '''   
                }
            }
        }
        stage('Test') {
            steps {
                withCredentials([string(credentialsId: 'github-token', variable: 'NETRC')]) {
                    sh '''
                        make test
                    '''
                }
            }
        }

    }
}