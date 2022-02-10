//go:build !change
// +build !change

package main

import "image/color"

var Cyan = color.RGBA{R: 100, G: 200, B: 200, A: 0xff}

const (
	Zero = `........
.111111.
.111111.
.11..11.
.11..11.
.11..11.
.11..11.
.11..11.
.11..11.
.111111.
.111111.
........`

	One = `........
.....11.
.....11.
.....11.
.....11.
.....11.
.....11.
.....11.
.....11.
.....11.
.....11.
........`

	Two = `........
.111111.
.111111.
.....11.
.....11.
.111111.
.111111.
.11.....
.11.....
.111111.
.111111.
........`

	Three = `........
.111111.
.111111.
.....11.
.....11.
.111111.
.111111.
.....11.
.....11.
.111111.
.111111.
........`

	Four = `........
.11..11.
.11..11.
.11..11.
.11..11.
.111111.
.111111.
.....11.
.....11.
.....11.
.....11.
........`

	Five = `........
.111111.
.111111.
.11.....
.11.....
.111111.
.111111.
.....11.
.....11.
.111111.
.111111.
........`

	Six = `........
.111111.
.111111.
.11.....
.11.....
.111111.
.111111.
.11..11.
.11..11.
.111111.
.111111.
........`

	Seven = `........
.111111.
.111111.
.....11.
.....11.
.....11.
.....11.
.....11.
.....11.
.....11.
.....11.
........`

	Eight = `........
.111111.
.111111.
.11..11.
.11..11.
.111111.
.111111.
.11..11.
.11..11.
.111111.
.111111.
........`

	Nine = `........
.111111.
.111111.
.11..11.
.11..11.
.111111.
.111111.
.....11.
.....11.
.....11.
.....11.
........`

	Colon = `....
....
....
....
.11.
.11.
....
.11.
.11.
....
....
....`
)
