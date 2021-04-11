pipeline{
    agent any

    environment {
        GO111MODULE = 'on'
        XDG_CACHE_HOME="/tmp/.cache"
    }

    stages {
        stage('Build') {
            steps {
                withCredentials([string(credentialsId: 'github-token', variable: 'NETRC')]) {
                        sh '''
                            make docker-image-build
                        '''
                }
            }
        }
        stage('Lint'){
            steps {
                script {
                    withCredentials([string(credentialsId: 'github-token', variable: 'NETRC')]) {
                        sh '''
                            make docker-lint
                        '''
                    }
                }
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