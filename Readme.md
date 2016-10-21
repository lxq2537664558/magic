<p align="left">
    <a href="https://vgo.io">
     <img  width="250" src="https://github.com/aiyun/openapm/blob/master/assets/images/vgo.png"></a>
</p>
Official Site
------------
<a href="http://vgo.io">http://openapm.io</a>

Ready To Use?
---------------
This project is still in early stage,but agent mecury is ready to use and we also used in our production enviroment

Features
------------
 - Pure Go
 - Modern architecture
 - Highly customizable
 - Simplicity
 - Extreme high performance
 - Easy installing and using

Functionality
-------------
 - Metric collect  (mecury)
 - Metric report  (mecury
 - Metric parse  (mecury)
 - Dapper log collect (mecury)
 - Metric store、show (influxdb + grafana)
 - Data streaming (vgo)
 - Alarm  (vgo)
 - Dapper log store 、 anlaysize、show(vgo)

Components
------------
 - openapm : The monitor platform,which is server side. 
 - mecury : A plug-in client side agent,Collec and report metrics to output(Default is to vgo).
 - etcd
 - influxdb
 - grafana
 
Installation(only mecury now)
------------
These part all are Runable applications,so,we could get with following command:  </br>
   ```bash
    $ go get -u github.com/aiyun/openapm
    $ cd $GOPATH/src/github.com/aiyun/openapm/mecury
   ```
Then we get two independent binarys in our $GOPATH/bin directory: 
   *  mecury.


Usage(only mecury now)
------------
 - For now ,you can use mecury to collect system info、java metrics info(see input plugins),and output to console(for debug)、influxdb、nats.
 - You can configure these in mecury.toml, [global_filters] is used to filtered out the undesired inputs or outputs.
 - In the default mecury.toml, we collect the system info and start the java metris httplistener, 
 and output these metrics to console for debug.For more details,please see the corresponding readme.md in 
the sub directories of pluginsy
 - Start the mecury with default mecury.toml and see what displays in your console!
