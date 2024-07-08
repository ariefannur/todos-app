pipeline {
    agent any 

    tools {
        go 'go1.20'
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
    }

    stages {
       stage('Build') {
            steps {
                go build
            }
        }
        stage('Test') {
            steps {
                go test -v
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
            }
        }
    }
}