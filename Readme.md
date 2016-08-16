Vgo
===
Modern and plug-in,Simplicity And Performance.
--------

#### Write entirely in Go. Much easier to use than others,while it's highly customizable.

This platform is composed of two parts:
    - vgo 
        The monitor platform,which is server side. 
    - mecury 
        A plug-in client side agent,Collec and report metrics to output(Default is to vgo). </br>

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

