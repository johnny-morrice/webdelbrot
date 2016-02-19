# Webdelbrot

Web-based mandelbrot viewer based on Godelbrot

http://github.com/johnny-morrice/godelbrot

You will probably also want the static data files

http://github.com/johnny-morrice/webdelbrot-static

## Usage

Webdelbrot behaves like other godelbrot tools.  It is configured by sending the output of configbrot to stdin.

    $ configbrot | webdelbrot

Then point your browser to localhost:8080.

Note the -addr argument allows you to specify network interface.

Webdelbrot has a few other flags.  Try -help.

## App controls

Left click to begin highlighting zoom region.  Left click again to zoom.

Middle quick or "q" to cancel zoom selection.

## Credits

**John Morrice**

http://functorama.com

https://github.com/johnny-morrice

## License

We use an MIT style license.  See LICENSE.txt for terms of use and distribution.
