pipeline{
    agent any

    environment {
        XDG_CACHE_HOME="/tmp/.cache"
    }

    stages {
        stage('Build') {
            steps {
                sh '''
                    make jk-docker-image-build
                '''
            }
        }
        stage('Lint'){
            steps {
                sh '''
                    make jk-docker-lint
                '''
            }
        }         
        stage('Test') {
            steps {
                    sh '''
                        make jk-docker-test
                    '''
            }
        }

    }
    post {
        always {
            sh '''
                make jk-docker-clean-all
            '''
        }
    }    
}