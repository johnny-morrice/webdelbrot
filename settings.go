package main

// These settings are good for development
const __ADDR = "godelbrot.functorama.mor"
const __PORT = 9898
const __PREFIX = ""
const __TICKTIME = 10
const __DEBUG = true

// These settings are better for running an internet demo
/*
const __ADDR = "godelbrot.functorama.com"
const __PORT = 80
const __PREFIX = ""
const __TICKTIME = 250
const __DEBUG = false
*/

// ADVANCED SETTINGS
const __RESIZE_MS = 300
const __ZOOM_MS = 1000 / 60
// Fraction remaining per second
const __SHRINK_RATE = 0.9