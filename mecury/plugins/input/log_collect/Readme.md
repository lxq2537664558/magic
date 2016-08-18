# log_collect Input Plugin

The log_collect plugin streams and parses the given logfiles.Beside the default ones you can also define your 
own log parsers!

### Configuration:

```toml
[[inputs.log_collect]]
  ## Log files to parse.
  ## These accept standard unix glob matching rules, but with the addition of
  ## ** as a "super asterisk". ie:
  ##   /var/log/**.log     -> recursively find all .log files in /var/log
  ##   /var/log/*/*.log    -> find all .log files with a parent dir in /var/log
  ##   /var/log/apache.log -> only tail the apache log file
  files = ["/var/log/apache/access.log"]
  ## Read file from beginning(If you want start read logs from last offset, please set to true).
  from_beginning = true

    ## the content of line log is : aabb
    ##Console Output-------------------------- 2016-08-18 15:17:10.003733768 +0800 CST
    ##test,host=mac.local log="aabb" 1471504613374704255
  [inputs.log_collect.raw]
    name = "test"
    field_name = "log"
    '''

```
