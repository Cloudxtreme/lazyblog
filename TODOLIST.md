### Todo

1. Add option to use `github.com/go-redis/redis` over BoltDB (allegedly it's "faster")
1. Add nice Markdown editor
1. Check bottlenecks with Uber's `go-torch`. -- I did this, and... found nothing interesting. The only thing I think I discovered is that you can respond to the request a _tiny_ bit faster if you return `*bytes.Buffer`s up through the store interface method to the handler, then use `buf.WriteTo(w)`, but that'd need to be benchmarked.
1. Implement cert generation.
1. Implement cert signing/challenge like Caddy does.
