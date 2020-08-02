# go-git-app

Note: please, use it as a part of bitbucket.org/vishjosh/bipp-go-git project:
- create a folder experimental-app in it's root
- run go test 
 It's the easiest way to integrate

1. Install Docker
2. Install docker-compose
3. Run: `docker-compose up`

4. Go to http://yourhost:9000/install and follow gitea installation wizard instructions
Don't forget to set `root url` and administration account

Change:

SSH Server Domain: yourhost
SSH Server Port: 222
Gitea Base URL: http://yourhost:9000/


In Administrator Account Settings set:

Administrator Username: gitea
Password: secret123
Email: gitea@gitea.com

5. Create in gitea repository with name "testrepo" under the user "gitea" to verify 
that everything is alright and for later usage in tests



Deployment:

Frontend building

1. cd go-git-app
2. docker build -t go-git-app app

Backend building

NOTE: you must have access to private repository via ssh. Nowadays, go modules have no integration with bitbucket. So, please, use github.com. For example, github.com/ujent/go-git. Or copy from bitbucket to your github repo.

1. cd go-git-app
2. ssh-add (only once)
3. DOCKER_BUILDKIT=1 docker build --ssh default -t go-git-api .
4. docker save -o go-git-api.tar go-git-api
5. Copy go-git-api.tar to server
6. docker load -i go-git-api.tar

To start application:
1. cd docker/app
2. docker-compose up -d

Local development:

UI:

To start a development server:
1. cd app
2. npm install
2. npm run start

Server:
1. To deploy test MySQL: 
    - cd docker/test
    - docker-compose up

2. APP_SERVER_PORT=4000 GIT_DB_CONN_STRING=root:secret@/gogittest FS_TYPE=1 go run main.go server.go 
or for using local filesystem
APP_SERVER_PORT=4000 FS_TYPE=2 GIT_ROOT=/home/ujent/code/go-git-app/testdata go run main.go server.go 


