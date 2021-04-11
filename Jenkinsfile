pipeline{
    agent any

    environment {
        XDG_CACHE_HOME="/tmp/.cache"
    }

    stages {
        stage('Build') {
            steps {
                sh '''
                    make docker-image-build
                '''
            }
        }
        stage('Lint'){
            steps {
                    sh '''
                        make docker-lint
                    '''
            }
        }         
        stage('Test') {
            steps {
                    sh '''
                        make docker-test
                    '''
            }
        }

    }
    post {
        always {
            sh '''
                make docker-network-clean docker-clean docker-clean-volumes
            '''
        }
    }    
}