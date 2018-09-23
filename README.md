[![pipeline status](https://gitlab.com/yakshaving.art/propaganda/badges/master/pipeline.svg)](https://gitlab.com/yakshaving.art/propaganda/commits/master)
[![coverage report](https://gitlab.com/yakshaving.art/propaganda/badges/master/coverage.svg)](https://gitlab.com/yakshaving.art/propaganda/commits/master)

# Propaganda

Announces merges in GitHub and GitLab to a slack incoming webhook.

Propaganda, from the spanish meaning: advertisement

> a notice or announcement in a public medium promoting a product, service, or event.

## Usage

```
Usage of ./propaganda:
  -address string
    	listening address (default ":9092")
  -config string
    	configuration file to use (default "propaganda.yml")
  -debug
    	enable debug logging
  -match-pattern string
    	match string regex (default "\\[announce\\]")
  -metrics string
    	metrics path (default "/metrics")
  -version
    	show version and exit
  -webhook-url string
    	slack webhook url (default SLACK_WEBHOOK_URL env var)
```

## Registering Webhooks

### Github

Be sure to pick `application/json` as content type, else all the webhooks will simply fail to be parsed.

### Gitlab

No particularities, simply pick only Merge Request events to reduce noise.
