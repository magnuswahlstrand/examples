# Durable Object in-memory storage: Location tracking

Example of a durable object that stores data in memory. Suitable for caching, but can be evicted.

[DEMO](https://durable-object-in-memory.wahlstrand.workers.dev/)

## Features

Will output geo information about the caller if called for the first time, and information about the PREVIOUS caller if called again.

First request:

```
This is the first request to this Durable object instance. Location was not set.

New state:
City: Stockholm
Country Code: SE
Postal Code: 100 29
Is EU Country: Yes
```

Second request

```
Durable object was already loaded with an in memory state. Updating state.  
    
Previous state:
City: Stockholm
Country Code: SE
Postal Code: 100 29
Is EU Country: Yes


New state:
City: Stockholm
Country Code: SE
Postal Code: 100 29
Is EU Country: Yes
```

## Run locally

```
pnpm install
pnpm dev
```

## Deploy

```
pnpm run deploy
```

## Inspiration

Inspired
by [Cloudflare Examples: Durables Objects in-memory state](https://developers.cloudflare.com/durable-objects/examples/durable-object-in-memory-state/).

Added a few features like:

* Typescript support
* Using RPC instead of HTTP fetch
* More information about the caller

