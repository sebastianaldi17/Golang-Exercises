# Golang-Exercises
Bite sized projects made using Golang.

## SimpleAPI
The hello world for CRUD API making. Account creation and update via POST, read whole list or certain account via GET, and delete account via DELETE. Uses `gorilla/mux` for routing. 

## LinkShortener
Link shortener like bit.ly and tinyurl. Uses `gorilla/mux` for routing and serving html files.

## QuizGame
Quiz game that reads from a two column csv file (containing the question on the first column and the answer on the second column). Uses a timer and goroutine to set a time limit for the quiz.

## RediGoCasino
A command line "casino" game using `redis` as a database. Uses `redigo` so that `go` can interact with `redis`.