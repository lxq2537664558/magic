<p align="center">
    <a href="https://vgo.io">
     <img  width="600" src="https://github.com/corego/vgo/blob/master/assets/images/vgo.png"></a>
</p>
Official Site
------------
<a href="http://vgo.io">vgo.io</a>

Features
------------
 - Modern and plug-in,Simplicity And Performance.
 - Write entirely in Go. Much easier to use than others,while it's highly customizable.


Components
------------
 - vgo : The monitor platform,which is server side. 
 - mecury : A plug-in client side agent,Collec and report metrics to output(Default is to vgo).

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

