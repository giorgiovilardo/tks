# tks

Trekin's Key Statistics is a powerful tool for analyzing football data.

`justfile` contains all the needed commands.

## Roadmap

* integrate the excel calculations;
* refactor to hexagonal architecture or in general more modular and testable code;
* add tests 😇;
* add CI;
* maybe add a sqlite database;
* learn `htmx` better?

## Retrospective

After some time to poke at the project with Rust and Python, given the constraints of

* having to distribute a binary to people using windows;
* wanting to embed resources in the binary to avoid extra files, keeping everything in one binary;
* not wanting to use a database, preferring to just use the csv files and do the calculations in the binary;

the decision to use Go was absolutely on point.
