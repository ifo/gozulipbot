## GoZulipBot

`gozulipbot` is a library to interact with Zulip in Go.
It is primarily targeted toward making bots.

## Installation

`go get github.com/ifo/gozulipbot`

## Usage

Make sure to add `gozulipbot` to your imports:

```go
import (
  gzb "github.com/ifo/gozulipbot"
)
```

Check out the examples directory for more info.

### Credentials
NB! Unlike in [matterbridge](https://github.com/42wim/matterbridge/wiki/Section-Zulip-(basic)), `APIURL` is the full path like `https://yourZulipDomain.zulipchat.com/api/v1/`

[`examples/`](examples/) use [`flag.GetConfigFromFlags`](flag.go#L10) to read credential setup from command line flags `--apiurl`, `--apiurl` and `--email`, or from the environment like

```
ZULIP_APIURL=https://yourZulipDomain.zulipchat.com/api/v1/
ZULIP_APIKEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
ZULIP_EMAIL=you@domain.tld
```
