Lazyblog
---

Lazyblog is viral-proof blogging platform.

*A what?*

Chances are, when you start a blog for this first time, you're going to opt for the cheapest hosting possible, and easiest setup. This will lead you to shared hosting and Wordpress -- nothing wrong with doing that!

However, let's say you write a blog post, and someone with a lot of Twitter followers shares it, or it gets posted on Reddit, or you get #1 on Hacker News, or something happens that makes your post go viral -- what happens?

Your site crashes!

But it doesn't have to be that way. With Lazyblog, your site is viral-proof! Using our secret formula, we've created a way for your blog posts to withstand incredible load. Just check out the bechmarks at the bottom of this doc!

Secret Formula
---

1. One(1) [httprouter](https://github.com/julienschmidt/httprouter).
2. Two(2) [BoltDB](https://github.com/boltdb/bolt) buckets.
3. Cache all rendered blog posts and JSON.

When you submit a completed blog post, that post (along with its template) are rendered to create a complete HTML page. That HTML page is then saved in a BoltDB bucket for that post. Next time someone visits the page that contains that blog post, the HTML is loaded directly from BoltDB, meaning that it doesn't need to be rendered again. This gives you incredibly low latency, and allows the server to handle a huge amount of requests.

Benchmark Results
---

**Cached post**:

```
λ wrk -d10 -c20 -t10 "http://localhost:3000/posts/here-is-me-post-56ce5f87"
Running 10s test @ http://localhost:3000/posts/here-is-me-post-56ce5f87
  10 threads and 20 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   450.93us  733.62us  20.26ms   97.15%
    Req/Sec     5.29k   451.35     8.56k    85.98%
  529239 requests in 10.10s, 262.46MB read
Requests/sec:  52393.49
Transfer/sec:     25.98MB
```

**Post JSON**:

```
λ wrk -d10 -c20 -t10 "http://localhost:3000/api/posts/here-is-me-post-56ce5f87"
Running 10s test @ http://localhost:3000/api/posts/here-is-me-post-56ce5f87
  10 threads and 20 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   448.74us  647.37us  14.53ms   96.56%
    Req/Sec     5.26k   334.52     6.43k    84.65%
  528850 requests in 10.10s, 135.17MB read
Requests/sec:  52353.43
Transfer/sec:     13.38MB
```

I think BoltDB might be the bottleneck in two tests above, will investigate one day.

**Homepage (not cached)**:

```
λ wrk -d10 -c20 -t10 "http://localhost:3000/"
Running 10s test @ http://localhost:3000/
  10 threads and 20 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.07ms    4.32ms 157.50ms   94.47%
    Req/Sec     1.45k   196.00     2.04k    72.40%
  144711 requests in 10.01s, 169.75MB read
Requests/sec:  14458.39
Transfer/sec:     16.96MB
```
