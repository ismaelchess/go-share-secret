pipeline{
    agent any

    environment {
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