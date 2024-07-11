pipeline {
    agent any 
    tools { go 'go1.19' }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
        REGISTRY_NAME = "simpletodos"
        ACR_LOGIN_SERVER = "${REGISTRY_NAME}.azurecr.io"
        REPOSITORY_NAME = "api"
    }

    stages {
       stage('Build') {
            steps {
                sh 'go build'
            }
        }
        stage('Test') {
            steps {
                sh 'go test -v'
            }
        }
        // Building Docker Image 
        stage ('Build Docker image') {
            steps {
                script {
                    sh "docker build -t ${REGISTRY_NAME}.azurecr.io/${REPOSITORY_NAME}:$BUILD_NUMBER ." 
                }
            }
        }
        stage('Deploy') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'azureKey', usernameVariable: 'AZURE_CRED_CLIENT_ID', passwordVariable: 'AZURE_CRED_CLIENT_SECRET')]) {
                    sh "docker login ${ACR_LOGIN_SERVER} -u $AZURE_CRED_CLIENT_ID -p $AZURE_CRED_CLIENT_SECRET"
                    sh " docker push ${REGISTRY_NAME}.azurecr.io/${REPOSITORY_NAME}:$BUILD_NUMBER"
                }
            }
        }
    }
}