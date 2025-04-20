# Camaretto

A fun little card game implemented in Golang. This project is part of the [YNC](https://yn-corp.xyz/home) project.

This project is mainly to learn go in more details as well as server/client communication. While golang has an optimized package for go-to-go program: gob; we decided to use protobuf to learn more about it.


- model: you can find all the app logic and struct making it alive; as well as the server and client structs created to handle online multiplayer.
    - component: Individuals elements following their own logic, implementing their own update and draw methods. Those components are to be reused by the Menu, Lobby & Game struct to interact with each other.
    - netplay: online multiplayer is implemented such as a routine is in the background communicating with the app. Sending & receiving data depending on the app state.
- View: Every routines or functions that are used to load or handle files. Every element found in those folders are to be used as much as possible by the `model/component/sprite.go` element.

In near future you'll be able to play this game directly in your web browser, [here](https://yn-corp.xyz/camaretto), thanks to a webassembly build.

__To build this project:__

    go build

__To run:__

    go run ./main.go
