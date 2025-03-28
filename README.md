# Camaretto

A fun little card game implemented in Golang. This project is part of the [YNC](https://yn-corp.xyz/home) project.

This project is mainly to learn go in more details as well as server/client communication. While golang has an optimized package for go-to-go program: gob; we decided to use protobuf to learn more about it.

The game has been conceived in a MVC model:

- Model: you can find all the app logic and struct making it alive; as well as the server and client structs created to handle online multiplayer.
    - component: all re-usable components are to be found in this module.
    - game: every moving part of the camaretto game (the playable part)
    - dialog: this regroups every dialog of each character in the game. Since there's a lot of them, and an intricate logic about rarety and character progression, this made sens to seperate them in their own module.
- Event (Controller): A simple EventQueue struct to handle ebiten controller logic. This also helps us reduce code in the `main.go` file.
- View: Every asset loader and Sprite implementation.

In near future you'll be able to play this game directly in your web browser, [here](https://yn-corp.xyz/camaretto), thanks to a webassembly build.

__To build this project:__

    go build

__To run:__

    go run ./main.go
