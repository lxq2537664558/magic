<p align="left">
    <a href="https://vgo.io">
     <img  width="250" src="https://github.com/corego/vgo/blob/master/assets/images/vgo.png"></a>
</p>
Official Site
------------
<a href="http://vgo.io">http://vgo.io</a>

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
 - vgo : The monitor platform,which is server side. 
 - mecury : A plug-in client side agent,Collec and report metrics to output(Default is to vgo).
 - etcd
 - influxdb
 - grafana
 
Installation
------------
These part all are Runable applications,so,we could get with following command:  </br>
   ```bash
    $ go get -u github.com/corego/vgo
    $ cd $GOPATH/src/github.com/corego/vgo/mecury && go install
    $ cd $GOPATH/src/github.com/corego/vgo/vgo && go install
   ```
Then we get two independent binarys in our $GOPATH/bin directory: 
   *  mecury.
   *  vgo.

