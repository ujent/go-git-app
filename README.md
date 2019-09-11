# go-git-app

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

NOTE: you must have access to private repository (github.com/ujent/go-git) via ssh
1. cd go-git-app
2. ssh-add (only once)
3. DOCKER_BUILDKIT=1 docker build --ssh default -t go-git-api .