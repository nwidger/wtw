wtw
===

`wtw` tells you what to wear on your run based on the current weather
and the type of run.  It uses data collected from Runner's
World's [What to Wear](http://www.runnersworld.com/what-to-wear) page.
Weather data is [Powered by Yahoo!](https://www.yahoo.com/?ilc=401).

## Try it out online

Try out the Javascript version of wtw in your
browser [here](https://nwidger.github.io/wtw).

## Build

### Mac OS X/Linux/Windows

```
$ go get -u github.com/nwidger/wtw/cmd/wtw
```

### Javascript

1. Install [GopherJS](https://github.com/gopherjs/gopherjs)

   ```
   go get -u github.com/gopherjs/gopherjs
   ```

2. Download wtw

   ```
   go get -d -u github.com/nwidger/wtw
   ```

3. Build `docs/main_js.js`

   ```
   cd $GOPATH/src/github.com/nwidger/wtw
   GOARCH=js gopherjs build github.com/nwidger/wtw/docs -m -o docs/main_js.js
   ```

4. Open `docs/index.html` in your browser.

## Usage

```
Usage of wtw:
  -browser
    	open answer in the browser
  -conditions string
    	c (clear), pc (partly cloudy), o (overcast), r (heavy rain), lr (light rain) or s (snowing) (default "c")
  -feel string
    	c (cool), ib (in between) or w (warm) (default "ib")
  -gender string
    	m (male) or f (female) (default "m")
  -intensity string
    	n (easy run), lr (long run), h (hard workout) or r (race) (default "n")
  -location string
    	get current conditions for location, overrides -temp, -conditions and -wind
  -temp int
    	temp (°F) (default 60)
  -time string
    	dawn, day, dusk, night or current (default "current")
  -v	print conditions before answer
  -wind string
    	nw (no wind), lw (light wind), hw (heavy wind) (default "nw")
```

## Example

With `-location`, `wtw` will retrieve the current weather from your
current location.  You only need to specify `-gender`, `-intensity`,
and `-feel`:

```
$ wtw -location 03801 -gender m -intensity n -feel ib
Sunglasses
Singlet
Shorts
Running Shoes
Sunblock
```

Without `-location`, you will need to specify `-temp`, `-conditions`
and `-wind` as well:

```
$ wtw -temp 60 -conditions c -wind nw -gender m -intensity n -feel ib
Sunglasses
Singlet
Shorts
Running Shoes
Sunblock
```

Most users will not need to specify `-time` as its value is
automatically selected based on the current time.

Specifying `-v` will cause `wtw` to print the conditions before the
answer:

```
$ wtw -location 03820 -v
wtw -gender m -temp 20 -conditions c -wind hw -time day -intensity lr -feel ib -v
Winter Cap
Sunglasses
Heavy Jacket
Long-Sleeve Shirt
Gloves
Tights
Running Shoes
Sunblock
```

This can be useful if for some reason the conditions determined by
`-location` or `-time current` don't quite match up with what you were
expecting.
