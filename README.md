# Meiji chocolate puzzle solver

## What is this?

I got nerdsniped one day with a Meiji chocolate puzzle. It's a puzzle that looks like a Meiji chocolate tablet. Each piece is a bunch of chocolate squares, and the overall pieces form a rectangle. The goal of the puzzle is to find a way to recreate that rectangle out of the pieces. This type of puzzle is also known as polyomino puzzle.

This solver is written in Go and can solve all 3 Meiji puzzles (white, milk, chocolate) in about 1~5 seconds. There's some luck involved in the search "algorithm" (brute-force), as it is a parallelized brute-force solver. It explores the search space in pretty random order, though in practice it will generally find solutions with pieces defined "earlier" in the puzzle definition as being generally closer to the top-left corner of the puzzle, since the solver stops looking as soon as it finds one valid solution.

The solver works by parsing the visual representation of pieces and creates bitmasks out of each possible orientation of each piece. Then it tries to fit every bitmask onto a board-sized bitfield until it finds a way to fill it all completely.

Enjoy!

## Solutions

All puzzles have multiple solutions. The following contains 2 solutions for each puzzle. The solver finds a random one on each run.

### Meiji White Chocolate Puzzle

<div align="center">
	<img src="https://github.com/EtiennePerot/meiji-chocolate-puzzle-solver/blob/master/img/white.jpg?raw=true" alt="Meiji White Chocolate Puzzle"/>
</div>

The white chocolate puzzle is a 8x5 grid (40 cells) with 8 puzzle pieces.

```shell
$ go run white.go
┌─────────┬─┬───┐
│F F F F F│A│G G│
│ ┌─┬─────┘ ├─┐ │
│F│C│A A A A│H│G│
├─┤ └─┬─────┤ │ │
│E│C C│D D D│H│G│
│ │   ├───┐ │ └─┤
│E│C C│B B│D│H H│
│ └───┴─┐ │ │   │
│E E E E│B│D│H H│
└───────┴─┴─┴───┘

$ go run white.go
┌───────┬─┬─────┐
│A A A A│B│D D D│
│ ┌───┬─┤ └─┬─┐ │
│A│H H│E│B B│G│D│
├─┤   │ ├───┘ │ │
│C│H H│E│G G G│D│
│ └─┐ │ └─────┼─┤
│C C│H│E E E E│F│
│   │ ├───────┘ │
│C C│H│F F F F F│
└───┴─┴─────────┘
```

### Meiji Milk Chocolate Puzzle

<div align="center">
	<img src="https://github.com/EtiennePerot/meiji-chocolate-puzzle-solver/blob/master/img/milk.jpg?raw=true" alt="Meiji Milk Chocolate Puzzle"/>
</div>

The Milk chocolate puzzle is a 10x6 grid (60 cells) with 12 puzzle pieces.

```shell
$ go run milk.go
┌───┬───────┬─┬─────┐
│G G│A A A A│C│J J J│
├─┐ │ ┌─┬───┘ └─┐ ┌─┤
│F│G│A│I│C C C C│J│L│
│ │ └─┤ └───┬───┤ │ │
│F│G G│I I I│H H│J│L│
│ └───┼─┐ ┌─┘ ┌─┴─┤ │
│F F F│D│I│H H│E E│L│
├───┬─┘ └─┤ ┌─┴─┐ │ │
│B B│D D D│H│K K│E│L│
│   └─┐ ┌─┴─┘ ┌─┘ │ │
│B B B│D│K K K│E E│L│
└─────┴─┴─────┴───┴─┘

$ go run milk.go
┌───────┬─────┬─────┐
│A A A A│B B B│E E E│
│ ┌─────┴─┐   │ ┌─┐ │
│A│C C C C│B B│E│D│E│
├─┴─┐ ┌───┴─┬─┼─┘ └─┤
│G G│C│K K K│I│D D D│
├─┐ ├─┘ ┌───┤ └─┐ ┌─┤
│F│G│K K│H H│I I│D│J│
│ │ └─┬─┘ ┌─┘ ┌─┴─┘ │
│F│G G│H H│I I│J J J│
│ └───┤ ┌─┴───┴───┐ │
│F F F│H│L L L L L│J│
└─────┴─┴─────────┴─┘
```

### Meiji Black Chocolate Puzzle

<div align="center">
	<img src="https://github.com/EtiennePerot/meiji-chocolate-puzzle-solver/blob/master/img/black.jpg?raw=true" alt="Meiji Black Chocolate Puzzle"/>
</div>

The Black chocolate puzzle is a 11x6 grid (66 cells) with 11 puzzle pieces. What makes it interesting is that its 66 cells mean it does not fit in a 64-bit bitfield.

There are only 2 solutions for the Black puzzle. Each of them is simply the 180-degree-rotated version of the other.

```shell
$ go run black.go
┌─────┬───────┬───┬───┐
│K K K│B B B B│I I│F F│
│   ┌─┼─────┐ └─┐ └─┐ │
│K K│G│E E E│B B│I I│F│
│ ┌─┘ └───┐ ├─┬─┘ ┌─┘ │
│K│G G G G│E│D│I I│F F│
├─┴─┐ ┌─┬─┘ │ ├───┴─┐ │
│H H│G│H│E E│D│C C C│F│
├─┐ └─┘ ├───┘ ├─┐   └─┤
│J│H H H│D D D│A│C C C│
│ └─────┴─┐ ┌─┘ └─────┤
│J J J J J│D│A A A A A│
└─────────┴─┴─────────┘

$ go run black.go
┌─────────┬─┬─────────┐
│A A A A A│D│J J J J J│
├─────┐ ┌─┘ └─┬─────┐ │
│C C C│A│D D D│H H H│J│
├─┐   └─┤ ┌───┤ ┌─┐ └─┤
│F│C C C│D│E E│H│G│H H│
│ └─┬───┤ │ ┌─┴─┘ └─┬─┤
│F F│I I│D│E│G G G G│K│
│ ┌─┘ ┌─┴─┤ └───┐ ┌─┘ │
│F│I I│B B│E E E│G│K K│
│ └─┐ └─┐ └─────┼─┘   │
│F F│I I│B B B B│K K K│
└───┴───┴───────┴─────┘
```

## License

This solver is licensed under the [WTFPL](http://www.wtfpl.net/).
