# Simple benchmark program. The test consists in sending sequential requests using 1 or more process.
# This test uses the ab tool (apache benchmark utility).
import subprocess

BENCHMARK_APP = "ab"
PROXY_ADDRESS = "localhost:8098"

# Types of tests.
test_types = [
    #   ( "test_description"             , executions_number (-n), concurrent_workers (-c) )
        ( "1 request, 1 process"         , 1                     , 1                       ),
        ( "10 requests, 1 process"       , 10                    , 1                       ),
        ( "10 requests, 10 processes"    , 10                    , 10                      ),
        ( "100 requests, 10 processes"   , 100                   , 10                      ),
        ( "1000 requests, 50 processes"  , 1000                  , 50                      ),
        ( "1000 requests, 100 processes" , 1000                  , 100                     ),
#        ( "10000 requests, 10 processes" , 10000                 , 10                      ),
]

object_list = [
    "http://i.imgur.com/t0x900S.jpg",                                             # 118401 bytes (~115 KB).
    "http://i.imgur.com/BMM9XWU.png",                                             # 52659 bytes (~51 KB).
    "http://graphics.stanford.edu/~seander/bithacks.html",                        # Less than 100 KB.
    "http://www.learn2crack.com/wp-content/uploads/2015/05/logo31.png",           # Less than 23 KB.
    "http://pbs.twimg.com/profile_images/603610759671611392/JRQtMqMR_normal.png", # With less than 2 KB (1655 bytes).
    "http://arquivos.oi.com.br/M1.zip",                                           # With 1MB size.
    "http://arquivos.oi.com.br/M5.zip",                                           # With 5MB size.
    "http://arquivos.oi.com.br/M10.zip",                                          # With 10MB size.
#    "http://arquivos.oi.com.br/M50.zip",                                          # With 50MB size.
]

# For each object, execute the test list.
for object_name in object_list:

    # For each type, run the test.
    for (test_description, executions_number, concurrent_workers) in test_types:
        print "**** File: {0}. Test: {1}. ****".format(object_name, test_description)
        command_line = "{0} -X {1} -n {2} -c {3} {4}".format(BENCHMARK_APP, PROXY_ADDRESS, executions_number, concurrent_workers, object_name)
        return_code = subprocess.call(command_line, shell=True)
        if return_code != 0:
            break
        print "\n\n\n"
