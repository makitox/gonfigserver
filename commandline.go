package main

import "flag"

var debug = flag.Bool("debug", false, "Add 'debug' option to start in debugging mode")

// var fswatch = flag.Bool("fswatch", false, "Watch changes in file system and reload properties map if any change")
// var git = flag.Bool("git", false, "Operate with property root as git repo")
// var loggerRoot = flag.String("log", "", "Specify log file full path")

var pFailOnDup = flag.Bool("failondup", true, "Fail to start if app finds key duplication")
var pReplaceChar = flag.String("replaceChar", ".", "Char to replace not allowed symbols in param key")
var iPort = flag.Int("port", 8080, "Specify port for Gonfig Server")
var pAddress = flag.String("address", "0.0.0.0", "Specify address for Gonfig Server")
var propertyRoot = flag.String("root", "", "Specify properties root")
