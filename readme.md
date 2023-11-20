# Check Semantic Version 

### generation executable linux and extract zip file
#### 1 terminal 
- `docker build -t daniel101/check-semantic-version .`
- `docker run -it daniel101/check-semantic-version`

#### 2 Terminal
- `docker ps`
- `docker cp YOUR_WORKING_CONTAINER_ID:/go/src/checkSemanticVersion/main.zip .`

### Add github permission
- `git update-index --chmod=+x ./scripts/executable`