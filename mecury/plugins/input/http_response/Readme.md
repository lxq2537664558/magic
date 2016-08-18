# Example Input Plugin

This input plugin will test HTTP/HTTPS connections.

### Configuration:

```
# HTTP/HTTPS request given an address a method and a timeout
[[inputs.http_response]]
  ## Server address(Required)
  address = "http://github.com"
  ## Set response_timeout (default 5 seconds)
  response_timeout = "5s"
  ## HTTP Request Method
  method = "GET"
  ## Whether to follow redirects from the server (defaults to false)
  follow_redirects = true
  ## HTTP Request Headers (all values must be strings)
  # [inputs.http_response.headers]
  #   Host = "github.com"
  ## Optional HTTP Request Body
  # body = '''
  # {'fake':'data'}
  # '''

  ## Optional SSL Config
  # ssl_ca = "/etc/mecury/ca.pem"
  # ssl_cert = "/etc/mecury/cert.pem"
  # ssl_key = "/etc/mecury/key.pem"
  ## Use SSL but skip chain & host verification
  # insecure_skip_verify = false
```

### Measurements & Fields:

- http_response
    - response_time (float, seconds)
    - http_response_code (int) #The code received

### Tags:

- All measurements have the following tags:
    - server
    - method

### Console Outputs
```
Console Output-------------------------- 2016-08-18 10:03:30.009504268 +0800 CST
http_response,host=mac.local,method=GET,server=http://github.com http_response_code=200i,response_time=2.340581561 1471485802000000000
```