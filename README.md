![Webdelbrot Mandelbrot Viewer](/webdelbrot.png)

# Webdelbrot

Web-based mandelbrot client of Godelbrot, the Unix-style Mandelbrot
renderer.

http://github.com/johnny-morrice/godelbrot


## Demo

http://godelbrot.functorama.com

## Usage

Webdelbrot is a pure gopherjs client of Godelbrot.

	# First launch Godelbrot webservice 
	$ configbrot -fix grow -numerics native | restfulbrot -debug -origins http://localhost:8080

	# Now serve with gopherjs
    $ gopherjs serve  -vw

## App controls

Left click to begin highlighting zoom region.  Hold longer for smaller zoom.  Right click to cancel zooming.

## Credits

**John Morrice**

http://functorama.com

https://github.com/johnny-morrice

## License

We use an MIT style license.  See LICENSE.txt for terms of use and distribution.
