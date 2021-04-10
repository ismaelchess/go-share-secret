pipeline{
    agent any

    tools {
        go 'go1.15.5'
    }
    environment {
        GO111MODULE = 'on'
        XDG_CACHE_HOME="/tmp/.cache"
    }

    stages {
        stage('Pre Test') {
            steps {
                    sh '''
                    make get
                    '''
            }
        }
        stage('Build') {
            steps {
                    sh '''
                        make build
                    '''
            }
        }
        stage('Test') {
            steps {
                    sh '''
                        make test
                    '''
            }
        }

    }
}