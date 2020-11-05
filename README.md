# github-actions-upstream-watch ![CI](https://github.com/bots-house/github-actions-upstream-watch/workflows/CI/badge.svg)

Trigger Github Actions Workflow when someone commit to another repo. 

## Usage

### Docker Compose

```yaml
version: '3.8'

volumes:
  watcher-state:
    driver: local

services:
  watcher:
    image: ghcr.io/bots-house/github-actions-upstream-watch:latest 
    environment: 
      # watch for new commits in github.com/tdlib/telegram-bot-api
      SRC_ORG: tdlib
      SRC_REPO: telegram-bot-api
      SRC_BRANCH: master
      
      # dispatch event upstream_commit to github.com/bots-house/docker-telegram-bot-api
      DST_ORG: bots-house
      DST_REPO: docker-telegram-bot-api
      DST_EVENT: upstream_commit

      # get this token from https://github.com/settings/tokens/new with repo:* perms
      TOKEN: {TOKEN}

      # check for new commits every 1 minute
      PERIOD: 1m
    volumes:
      - watcher-state:/data
```
