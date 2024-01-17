# Check Semantic Version 

### generation executable linux and extract zip file
#### 1 terminal 
- `docker build -t daniel101/check-semantic-version .`
- `docker run -it daniel101/check-semantic-version`

#### 2 Terminal
- `docker ps`
- `docker cp YOUR_WORKING_CONTAINER_ID:/go/src/checkSemanticVersion/executable.zip .`

### Add github permission
- `git update-index --chmod=+x ./scripts/executable`

### How to use

- `./scripts/executable verify -versionOld v0.0.1 -versionNew v0.0.2`
- `./scripts/executable update-github-vars -owner OWNER -repository REPOSITORY -varName VAR_NAME -varValue VAR_VALUE -gbtoken PA_TOKEN`
- `./scripts/executable commits-verify -owner OWNER -repository REPOSITORY -merge_request_id 65 -bypass NO -gbtoken PA_TOKEN`