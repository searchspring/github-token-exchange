# Github Token Exchange

Takes an Oauth redirect code and exchanges it for a token and the user info and sends that info to the following URL

```
http://localhost:3827?user={"node_id":"value","avatar_url":"value","name":"value","token":"value"}
```

## Run with docker

Create a `.env` file and add the values: 
```
GITHUB_CLIENT_SECRET=*** find the secret in 1password ***
GITHUB_CLIENT_ID=e02c8965ff92aa84b6ee
GITHUB_REDIRECT_URL=http://localhost:3000
ALLOWLIST_REDIRECT_URLS=http://localhost,https://localhost,https://searchspring.github.io/snapp-explorer
```

Run the docker build

```bash
make docker-build
```

Run the docker image

```bash
make docker-run
```

## Run locally


Create a `.env` file as above.

```bash
make run
```


## Prod credentials

```
GITHUB_CLIENT_SECRET=*** find the secret in 1password ***
GITHUB_CLIENT_ID=5df635731e7fa3513c1d
GITHUB_REDIRECT_URL=https://token.kube.searchspring.io
ALLOWLIST_REDIRECT_URLS=http://localhost,https://localhost,https://searchspring.github.io/snapp-explorer
```
