# Cloudflare-DynDNS

## Description

cloudflare-dyndns is a CLI application written in Go and delivered as a single executable that lets you run your own dynamic DNS service with Cloudflare.

An api token is needed to run the tool -> [how to create a cloudflare api token](https://developers.cloudflare.com/api/tokens/create/)

## Todos

- Docker Image


## Usage

### Creates DNS Records

```bash
❯ cloudflare-dyndns update --record "dyn-dns-test" --zone "test.com" --api-token "mytoken"
Creating A Record with content "1.1.1.1" & ttl 1
Creating AAAA Record with content "2a02::123" & ttl 1
```

### Updates DNS Records

```bash
❯ cloudflare-dyndns update --record "dyn-dns-test" --zone "test.com" --api-token "mytoken"
Updating A Record with content "1.1.1.2" & ttl 1
No update needed for AAAA Record
```

### Help

```
❯ cloudflare-dyndns update --help
Update DynDNS Entry

Usage:
  cloudflare-dyndns update [flags]

Flags:
  -t, --api-token string   cloudflare api token
  -h, --help               help for update
      --log-level string   log level (default "warning")
  -r, --record string      dns record to change
      --ttl int            ttl for record (default 1)
  -z, --zone string        zone name
```
