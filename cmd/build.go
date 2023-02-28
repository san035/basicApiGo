/*
При компиляции использовать
go build -v -ldflags="-X 'main/build.Time=$(date)'"
тогда в переменной main.Time будет время компиляции программы
*/

package main

var TimeBuild string

//var Version = ""
