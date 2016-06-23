Lazyblog [![Build Status](https://semaphoreci.com/api/v1/bentranter/lazyblog/branches/master/shields_badge.svg)](https://semaphoreci.com/bentranter/lazyblog)
---

Lazyblog: The Viral-Proof Blogging Platform

*The what?*

Chances are, when you start a blog for this first time, you're going to opt for the cheapest hosting possible, and easiest setup. This will lead you to shared hosting and Wordpress -- nothing wrong with doing that!

However, let's say you write a blog post, and someone with a lot of Twitter followers shares it, or it gets posted on Reddit, or you get #1 on Hacker News, or something happens that makes your post go viral -- what happens?

Your site crashes!

But it doesn't have to be that way. With Lazyblog, your site is viral-proof! Using our secret formula, we've created a way for your blog posts to withstand incredible load. Just check out the bechmarks at the bottom of this doc!

Notice
---

There are two versions of this: the crappy old one (that worked but the code is hard to read) and the new, beautiful one.

To run the old, working code, navigate to the `cmd` directory and use either `go run` or `go build`.

If you want to try the new stuff, use the same commands in the top level of this repo. It'll use the `main.go` file there.

If that's not enough, check out of the work-in progress branches. Those branches have been benchmarked around 275,000 req/sec with the hope of getting up to 450,000 req/sec before merging back to master.

Usage
---

Setting up a new "Lazyblog" is (almost) very simple. Just clone and compile the repo, and run `lazyblog -username="<your-username>" -password="<your-password>` to setup your site, then login with those same credentials at `/admin/login`. If you want to use this in development mode, you'll need to set the environment variable `LAZYBLOG_ENV=dev`, otherwise it will default to running on port 80.

In the future, I _promise_ I'll make this much easier.

Kind-Of Lazy Deployment
---

With Supervisor installed, you can use the file,

```conf
[program:lazyblog]
command=/root/go/src/github.com/<your-fork>/lazyblog/cmd/cmd -username="<your-username>" -password="<your-password>"
autostart=true
autorestart=true
startretries=10
user=<your-user>
redirect_stderr=true
stdout_logfile=/var/log/supervisor/lazyblog.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=10
```

Secret Formula
---

1. One(1) [httprouter](https://github.com/julienschmidt/httprouter).
2. Two(2) [BoltDB](https://github.com/boltdb/bolt) buckets.
3. Cache all rendered blog posts and JSON.

When you submit a completed blog post, that post (along with its template) are rendered to create a complete HTML page. That HTML page is then saved in a BoltDB bucket for that post. Next time someone visits the page that contains that blog post, the HTML is loaded directly from BoltDB, meaning that it doesn't need to be rendered again. This gives you incredibly low latency, and allows the server to handle a huge amount of requests.

Benchmark Results
---

**New Version**

```
Running 30s test @ http://localhost:8000
  8 threads and 32 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     5.64ms    9.31ms 161.28ms   88.63%
    Req/Sec    48.38k     7.78k  169.37k    80.93%
  Latency Distribution
     50%    3.38ms
     75%    6.25ms
     90%   18.19ms
     99%    0.00us
  11572339 requests in 30.10s, 1.62GB read
Requests/sec: 384504.41
Transfer/sec:     55.00MB
```

**Old Version**

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

I've tested this, and I think I've made a discovery: BoltDB itself is not the bottleneck - the bottleneck is the HTTP package! That means that your content literally cannot be served any faster using the normal HTTP package. Pretty cool to think that your content _feels_ dynamic, but is served as fast as `net/http` can serve some "Hello, World" bytes.

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

Available on all platforms. Run `make release` to build for every platform.
