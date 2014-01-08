gomahawk
========

WIP implementation of the [Network Protocol of tomahawk](http://wiki.tomahawk-player.org/index.php/Network_Protocol) in golang

Install
=======

	go get github.com/MStoykov/gomahawk


Usage 
=====

Import
------
	
	import (
		"github.com/MStoykov/gomahawk"
	)

Implement Tomahawk Interface
----------------------------

look at [examples](./examples/)


Make new Gomahawk instance
--------------------------

	g := NewGomahawkImpl()
	gs := gomahawk.NewGomahawkServer(g)
	err := gs.ListenTo(net.IPv4(192, 168, 1, 13), "50210") // Listen on 192.168.1.13 and the default tomahawk port
	// error checking
	gs.Start()

Requesting DBConnection
-----------------------
This is used to get the changes from the remote

	var gs GomahawkServer
	t := gs.Tomahawks()[0] // get the first of the connected Tomahawks
	err := t.RequestDBConnection() 
	// error handleing
	// the Gomahawk registered with the GomahawkServer gs will get call to 
	// NewDBConnection with the Tomahawk instance and the given DBConnection can be used to initializa fetching
	dbConnection.FetchOps(fetchOps, "") // get all changes 
	// the methods on fetchOps will be called sequentually and at the end the Close() method will be called
 	// that signals that all current changes have been transmitted
 


LICENSE
=======
The MIT License (MIT)

Copyright (c) 2013 Mihail Stoykov

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
