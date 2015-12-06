# YAPC
Yet another proxy caching (YAPC) is an open source Transparent HTTP Proxy Caching, similar to NGINX, ATS and SQUID.

My goal with this project is to learn Golang. Of course, if someone interesting in contribute or use, feel free to fork, contact me, ...

Nowadays, I only started to think about what the project is (and isn't), your basic architecture and planned versions. I did some work and I made some progress.

Bellow and at the wiki page there are more information.


# Planned versions
* ~~Only proxy. All requests are sent to the upstream. All requests are logged. Signals are intercepted. Initial performance tests. Initial BDD tests. TDD tests.~~
* Proxy with basic access control. With the map configuration will be possible to specifier what can be accessed and what cannot.
* ~~Proxy + cache (memory only)~~
* Proxy + cache (memory or disk)
* Transparent caching proxy
