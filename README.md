wtw
===

`wtw` tell you what to wear on your run based on the current weather
and the type of run.  It uses data collected from Runner's
World's [What to Wear](http://www.runnersworld.com/what-to-wear) page.

## Usage

```
Usage of wtw:
  -conditions string
    	c (clear), pc (partly cloudy), o (overcast), r (heavy rain), lr (light rain) or s (snowing) (default "c")
  -feel string
    	c (cool), ib (in between) or w (warm) (default "ib")
  -gender string
    	m (male) or f (female) (default "m")
  -intensity string
    	n (easy run), lr (long run), h (hard workout) or r (race) (default "n")
  -temp int
    	temp (°F) (default 60)
  -time string
    	dawn, day, dusk or night (default "day")
  -wind string
    	nw (now win), lw (light wind), hw (heavy wind) (default "nw")
```

## Example

```
$ wtw -gender m -temp 60 -conditions c -wind nw -time day -intensity n -feel ib
Sunglasses
Singlet
Shorts
Running Shoes
Sunblock
```

