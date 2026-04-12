# Ultimate Tic Tac Toe
A Go application for playing Ultimate Tic Tac Toe.

To run the application execute the following command:

```
go run ./tictactoe.go
```

The program will have the player play against the AI. It allows the player to select if they will play as X and O (determining who plays first) as well as the AI difficulty.

The "easy" difficulty uses the alpha-beta algorithm with a search depth of just 1 move. "Medium" uses a search depth of 6 moves.

The "hard" difficulty uses the MCTS (Monte Carlo Tree Search) algorithm with a time limit of 5 seconds. This MCTS algorithm is optimized to run simulations on multiple CPU cores using Goroutines and channels. This helps the algorithm run more simulations in order to make better estimations at the best next move. On my PC with two CPU cores and four logical CPUs, this increases the number of simulations run by about 80 %.
