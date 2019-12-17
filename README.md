# Orderbook

[![Go Report Card](https://goreportcard.com/badge/github.com/piquette/orderbook)](https://goreportcard.com/badge/github.com/piquette/orderbook)
[![Build Status](https://travis-ci.org/piquette/orderbook.svg?branch=master)](https://travis-ci.org/piquette/orderbook)
[![GoDoc](https://godoc.org/github.com/piquette/orderbook?status.svg)](https://godoc.org/github.com/piquette/orderbook)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Purpose
This project is intended as a way to build and explore the dynamics of performant orderbooks through implementing various data structures.

## Caveats
There are some architectural caveats (simplifications) that are made here to keep things nice. Some of them are:

* Concurrency is ignored here, operations on the order book are purely single-threaded and transactional.
* Orders are added and matched in continuous-time and are valid until cancelled.
* Prices are represented as integers for computational and educational simplicity.
* Market orders and advanced order type logic is ignored, every order must be submitted at a specific price.


## Operations
* Submit Order
* Cancel Order
* Get Top of Book
