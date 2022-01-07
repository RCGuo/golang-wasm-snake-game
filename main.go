package main

func main() {
	StartGame()
	<-make(chan struct{})
}
