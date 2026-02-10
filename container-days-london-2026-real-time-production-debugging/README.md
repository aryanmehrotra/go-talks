# Real-Time Production Debugging (Container Days London 2026)
## Go Profiling Resources (pprof + GoFr)

A small collection of resources for profiling Go services using `pprof`,
especially when building microservices with GoFr.

Profiling is usually the fastest way to understand CPU usage, memory growth,
goroutine buildup, and latency issues in production systems.

---

## Maintainers / Speakers

**Aryan Mehrotra**

LinkedIn: [https://www.linkedin.com/in/aryanmehrotra](https://www.linkedin.com/in/aryanmehrotra)

X: [https://x.com/_aryanmehrotra](https://x.com/_aryanmehrotra)

GitHub: [https://github.com/aryanmehrotra](https://github.com/aryanmehrotra)


**Umang Mundhra**

LinkedIn: [https://www.linkedin.com/in/umang01-hash](https://www.linkedin.com/in/umang01-hash)

X: [https://x.com/umang01hash](https://x.com/umang01hash)

GitHub: [https://github.com/Umang01-hash](https://github.com/Umang01-hash)


---

## GoFr

Repository:
[https://github.com/gofr-dev/gofr](https://github.com/gofr-dev/gofr)

Documentation:
[https://gofr.dev/docs](https://gofr.dev/docs)

GoFr exposes standard Go runtime behavior, so all native Go profiling tools
(`pprof`, tracing, etc.) work without any framework-specific changes.

---

## pprof Documentation

Official Go profiling docs:
[https://pkg.go.dev/net/http/pprof](https://pkg.go.dev/net/http/pprof)

Go blog (profiling):
[https://go.dev/blog/pprof](https://go.dev/blog/pprof)

Interactive pprof guide:
[https://github.com/google/pprof/blob/main/doc/README.md](https://github.com/google/pprof/blob/main/doc/README.md)

---

## Talks

Using GoFr for building Microservices - Aryan Mehrotra & Umang Mundhra (GopherCon)

[https://youtu.be/EAMwtmaZoPY?si=OXX0y5c7_Ygg1Wor](https://youtu.be/EAMwtmaZoPY?si=OXX0y5c7_Ygg1Wor)

Two Go Programs, Three Profiling Techniques — Dave Cheney (GopherCon)

[https://www.youtube.com/watch?v=nok0aYiGiYA](https://www.youtube.com/watch?v=nok0aYiGiYA)

Profiling Request Latency with Critical Path Analysis (GopherCon)

[https://www.youtube.com/watch?v=BayZ3k-QkFw](https://www.youtube.com/watch?v=BayZ3k-QkFw)

Building a Go Profiler with eBPF (GopherCon Europe)

[https://www.youtube.com/watch?v=OlHQ6gkwqyA](https://www.youtube.com/watch?v=OlHQ6gkwqyA)

---

## Notes

Profiling should be part of normal debugging — not something saved for
“performance emergencies”.

Most production issues in Go services come down to:

* unexpected allocations
* blocking operations
* goroutine leaks
* lock contention
