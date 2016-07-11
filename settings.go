package main

// These settings are good for development
const __ADDR = "localhost"
const __PORT = 9898
const __PREFIX = ""
const __TICKTIME = 100
const __DEBUG = true
const __TRACE = false

// These settings are better for running an internet demo
/*
const __ADDR = "godelbrot.functorama.com"
const __PORT = 80
const __PREFIX = ""
const __TICKTIME = 250
const __DEBUG = false
const __TRACE = false
*/

// ADVANCED SETTINGS
const __RESIZE_MS = 300
const __ZOOM_MS = 1000 / 60
// Fraction remaining per second
const __SHRINK_RATE = 0.5