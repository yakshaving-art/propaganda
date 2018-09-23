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
  -webhook-url string
    	slack webhook url (default SLACK_WEBHOOK_URL env var)
```
