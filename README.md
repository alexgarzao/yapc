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


# Tasks
* Define the possibles module's types and where them can be executed
  * Some modules can generate content
  * Some modules can change the content that was read from the upstream
  * Some modules can change the content that was read from the cache
  * In specific situations, one module can generate the headers' content or alter them
  * We cannot forget that headers are the first content that must be sent
* BDD, TDD, performance tests, stress tests, FDD, …
* Focus on test since the first line of code!!!!
* Varnish tools are valuable? It's funny to implement, but they are useful?
  * varnishlog, varnish top, varnishhist, varnishstat
* purge?


# Sample config
    .*.youtube.com:	pass
    aabbcc.*		pass
    abcxy.br		-> 192.168.1.2 // proxy
    // Is it a good approach? Maybe, using a module would be more interesting, but how can I integrate this on the config? Well, a NGX's module can register your tokens in configuration…
    .*			fail // proxy reverso
    .*			accept // proxy ou proxy cache

Future, in this config, I can put cache parameters (global and per rule), TTL, ...


# Proxy architecture

Components:
* Config manager: Responsible for read the configuration, generate the sources, compile and build the new binary;
* Startup ou Setup Manager: Responsible for Load the plugins (or they will be together with the binary?), network setup (bind, listen);
* Requests manager: Responsible for accept the requests, read them, match and validate in the config, hit/read or miss/fetch, send to modules' pipeline, save in storage, send to downstream;
* Permanent log manager;
* Object manager: Responsible for keep the object's' index in memory to provide a fast search (hit or miss?);
* Storage manager: Responsible for save and load the objects. Is possible to have a lot of storage managers like, for example, memory storage, fdd storage, ssd storage, MMAP version, one file per object, one file with a lot of objects, … And, in some cases, the objects are unique in the storages, but in other cases the objects can be duplicated (for fast retrieve);
* Evict manager: Responsible for select the objects that will be removed;
* Memory manager: Responsible for the memory management, avoiding syscall to OS;
* Statistics manager: Responsible for keep object's statistics like hits, misses, fetch, object count, KB in/out per second, new objects, removed objects, and so on…;
* Core: Responsible for the workflow of the other components, signals treatment and error handler.
