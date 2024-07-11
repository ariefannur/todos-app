pipeline {
    agent any 
    tools { go 'go1.19' }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
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
        stage('Deploy') {
            steps {
                azCommands('servicePrincipalId', ['az login --service-principal -u $AZURE_CRED_CLIENT_ID -p $AZURE_CRED_CLIENT_SECRET -t $AZURE_CRED_TENANT_ID'])
            }
        }
    }
}