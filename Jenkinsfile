pipeline{
    agent any

    environment {
        XDG_CACHE_HOME="/tmp/.cache"
    }

    stages {
        stage('Pre Test') {
            steps {
                    sh '''
                    whoami
                    echo "=============="
                    echo $PATH
                    go version; \\
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