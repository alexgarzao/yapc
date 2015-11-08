# YAPC
Yet another proxy caching (YAPC) is an open source Transparent HTTP Proxy Caching, similar to NGINX, ATS and SQUID.

My goal with this project is to learn Golang. Of course, if someone interesting in contribute or use, feel free to contact me :-)

Nowadays, I only started to think about what the project is (and isn't), your basic architecture, planned versions, ...

Bellow I will put what in my mind :-)


# Planned versions
* Only proxy, with basic access control. With the map configuration will be possible to specifier what can be accessed and what cannot.
* Proxy + cache
* Transparent caching proxy

# Doubts
* Is possible one cache system reaches the goals of the three versions above?
* Reverse proxy is basically a proxy cache that access only local content?
* Where (and how) the URL could be rewritten?
* How NGX define domain's rules and resources' rules?
* Log manager/statistics: is a good idea uses the same approach that varnish? Or it must be considered an improvement?
* For implement BDD, Must be used Python+Lettuce? Or is better to try a BDD framework specific for Golang? Or there is another approach?
* Is a good idea to have more than one option to do the storage, log, and so on?
  * Maybe, based in the config, I can rewrite some code to use a specific storage backend, log system, … With this approach, this project can be used like a framework to construct a personalized proxy cache system :-)

# Startup sequence (final version)
* If the config was changed, generate the source's configuration, compile and build the new proxy binary
* Run the proxy binary
* Bind in the network port
* Listen
* For each connection request
  * Accept
  * Read request
  * Match request in configuration
  * Send 404 if appropriate
  * Connect in the upstream
  * While there is object data
    * Read object data (4KB buffer)
    * Send data to modules' pipeline
    * Send data to cache (storage)
    * Send data to downstream
