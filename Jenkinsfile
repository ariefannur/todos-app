pipeline {
    agent any 

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