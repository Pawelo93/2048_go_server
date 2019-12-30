package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Answer struct {
	Direction string
}

type Data struct {
	Board [][]int
}

func main() {
	// http.HandleFunc("/load", onLoad)
	// http.ListenAndServe(":8080", nil)

	var newGrid = [][]int{
		{0, 0, 0, 0},
		{0, 4, 0, 0},
		{0, 2, 2, 0},
		{0, 0, 0, 0},
	  }
	var grid = move(newGrid, "RIGHT")
	fmt.Println(grid)
}

func onLoad(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	boardJson := r.FormValue("board")
	fmt.Println(boardJson)

	var data Data
	err := json.Unmarshal([]byte(boardJson), &data)

	if err != nil {
		fmt.Fprintln(w, err.Error())
	}

	bestDirection := getBestDirection(data.Board)

	answer := Answer{Direction: bestDirection}
	b, err := json.Marshal(answer)

	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, string(b))
		fmt.Println(string(b))
	}
}

func getBestDirection(board [][]int) string {
	return "RIGHT"
}

func move(passedBoard [][]int, direction string) [][]int {
    var actualBoard [][]int
    if (direction == "DOWN") {
      actualBoard = passedBoard;
    } else if (direction == "RIGHT") {
      actualBoard = rotateTableLeft(passedBoard);
      actualBoard = rotateTableLeft(actualBoard);
      actualBoard = rotateTableLeft(actualBoard);
    } else if (direction == "UP") {
      actualBoard = rotateTableLeft(passedBoard);
      actualBoard = rotateTableLeft(actualBoard);
    } else if (direction == "LEFT") {
      actualBoard = rotateTableLeft(passedBoard);
    }

    for column := 0; column < 4; column++ {
      _moveDown(actualBoard, column);
      _join(actualBoard, column);
      _moveDown(actualBoard, column);
    }

	var resultBoard [][]int
    if (direction == "DOWN") {
      resultBoard = actualBoard;
    } else if (direction == "RIGHT") {
      resultBoard = rotateTableLeft(actualBoard);
    } else if (direction == "UP") {
      resultBoard = rotateTableLeft(actualBoard);
      resultBoard = rotateTableLeft(resultBoard);
    } else if (direction == "LEFT") {
      resultBoard = rotateTableLeft(actualBoard);
      resultBoard = rotateTableLeft(resultBoard);
      resultBoard = rotateTableLeft(resultBoard);
    }

    return resultBoard;
  }

  func _moveDown(actualBoard [][]int, column int) {
    for times := 0; times < 3; times++ {
      for row := 3; row > 0; row-- {
        if (actualBoard[row][column] == 0) {
          actualBoard[row][column] = actualBoard[row - 1][column];
          actualBoard[row - 1][column] = 0;
        }
      }
    }
  }

  func _join(actualBoard [][]int, column int) {
    for row := 3; row > 0; row-- {
      var value = actualBoard[row][column];
      var valueUp = actualBoard[row - 1][column];
      if (value == valueUp) {
        actualBoard[row][column] = actualBoard[row - 1][column] * 2;
        actualBoard[row - 1][column] = 0;
      }
    }
  }

  func rotateTableLeft(board [][]int) [][]int{
    var newGrid = [][]int{
      {0, 0, 0, 0},
      {0, 0, 0, 0},
      {0, 0, 0, 0},
      {0, 0, 0, 0},
	}
    for i := 0; i < 4; i++ {
      for j := 0; j < 4; j++ {
        var value = board[i][j];

        newGrid[3 - j][i] = value;
      }
    }
    return newGrid;
  }