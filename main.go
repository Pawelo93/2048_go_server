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

type MoveResult struct {
	board [][]int
	score int
}

func main() {
	// http.HandleFunc("/load", onLoad)
	// http.ListenAndServe(":8080", nil)

	var newGrid = [][]int{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 2, 2, 0},
		{0, 0, 0, 0},
	}
	// var grid = move(newGrid, "RIGHT")
	var nextMove = minimax(newGrid, 0, 1, true)
	fmt.Println(nextMove)
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

func move(passedBoard [][]int, direction string) MoveResult {
	var actualBoard [][]int
	var score int
	if direction == "DOWN" {
		actualBoard = passedBoard
	} else if direction == "RIGHT" {
		actualBoard = rotateTableLeft(passedBoard)
		actualBoard = rotateTableLeft(actualBoard)
		actualBoard = rotateTableLeft(actualBoard)
	} else if direction == "UP" {
		actualBoard = rotateTableLeft(passedBoard)
		actualBoard = rotateTableLeft(actualBoard)
	} else if direction == "LEFT" {
		actualBoard = rotateTableLeft(passedBoard)
	}

	for column := 0; column < 4; column++ {
		_moveDown(actualBoard, column)
		score += _join(actualBoard, column)
		_moveDown(actualBoard, column)
	}

	var resultBoard [][]int
	if direction == "DOWN" {
		resultBoard = actualBoard
	} else if direction == "RIGHT" {
		resultBoard = rotateTableLeft(actualBoard)
	} else if direction == "UP" {
		resultBoard = rotateTableLeft(actualBoard)
		resultBoard = rotateTableLeft(resultBoard)
	} else if direction == "LEFT" {
		resultBoard = rotateTableLeft(actualBoard)
		resultBoard = rotateTableLeft(resultBoard)
		resultBoard = rotateTableLeft(resultBoard)
	}

	return MoveResult{board: resultBoard, score: score}
}

func _moveDown(actualBoard [][]int, column int) {
	for times := 0; times < 3; times++ {
		for row := 3; row > 0; row-- {
			if actualBoard[row][column] == 0 {
				actualBoard[row][column] = actualBoard[row-1][column]
				actualBoard[row-1][column] = 0
			}
		}
	}
}

func _join(actualBoard [][]int, column int) int {
	var score int = 0
	for row := 3; row > 0; row-- {
		var value = actualBoard[row][column]
		var valueUp = actualBoard[row-1][column]
		if value == valueUp {
			var joinValue = actualBoard[row-1][column] * 2
			actualBoard[row][column] = joinValue
			actualBoard[row-1][column] = 0
			score += joinValue
		}
	}
	return score
}

func rotateTableLeft(board [][]int) [][]int {
	var newGrid = [][]int{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			var value = board[i][j]

			newGrid[3-j][i] = value
		}
	}
	return newGrid
}

// MINIMAX
const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func minimax(theBoard [][]int, score int, depth int, player bool) map[string]interface{} {
	result := make(map[string]interface{})
	var bestDirection string
	var bestScore int
	fmt.Printf("%d\n", depth)
	if depth == 0 || boardIsFull(theBoard) {
		fmt.Println("Depth is 0 (" + string(depth) + ") or board is full")
		bestScore = heuristicScore(score, getNumberOfEmptyCells(theBoard), calculateClusteringScore(theBoard))
	} else {
		if player == true {
			bestScore = MaxInt

			for _, direction := range [4]string{"RIGHT", "LEFT", "UP", "DOWN"} {
				var newBoard = theBoard

				var moveResult = move(newBoard, direction)

				var points = moveResult.score

				fmt.Println("New Direction " + direction)

				if points == 0 && areEqual(moveResult.board, theBoard) {
					continue
				}

				var currentResult map[string]interface{} = minimax(moveResult.board, moveResult.score, depth-1, false)
				var currentScore int = currentResult["Score"].(int)

				fmt.Println("In function current score " + string(currentScore))
				if currentScore > bestScore {
					//maximize score
					fmt.Println("New best score for direction " + direction)

					bestScore = currentScore
					bestDirection = direction
				}
			}
		} else {
			bestScore = MinInt

			var moves []int = getEmptyCellIds(theBoard)
			if len(moves) == 0 {
				bestScore = 0
			}

			var i, j int
			for _, cellId := range moves {
				i = (cellId / 4)
				j = cellId % 4

				for _, value := range [2]int{2, 4} {
					var newBoard [][]int = theBoard
					newBoard[i][j] = value

					var currentResult map[string]interface{} = minimax(newBoard, score, depth-1, true)
					var currentScore int = currentResult["Score"].(int)
					if currentScore < bestScore {
						//minimize best score
						bestScore = currentScore
					}
				}
			}
		}
	}

	fmt.Println("Best Score %d\n", bestScore)
	fmt.Println("Best Direction " + bestDirection)
	result["Score"] = int(bestScore)
	result["Direction"] = bestDirection

	return result
}

func getEmptyCellIds(theBoard [][]int) []int {
	var array []int
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if theBoard[i][j] == 0 {
				array = append(array, 4*i+j)
			}
		}
	}
	return array
}

func boardIsFull(board [][]int) bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

func areEqual(first [][]int, second [][]int) bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if first[i][j] != second[i][j] {
				return false
			}
		}
	}
	return true
}

func heuristicScore(actualScore int, numberOfEmptyCells int, clusteringScore int) int {
	//    print('actualScore $actualScore numberOfEmptyCells $numberOfEmptyCells, clusteringScore $clusteringScore, log ${log(actualScore)}');
	// int c;
	// if(actualScore > 0)
	// 	c = (log(actualScore) * numberOfEmptyCells).round();
	// else
	// 	c = 1;
	//
	// int score = (actualScore + c - clusteringScore);
	// return max(score, min(actualScore, 1));
	return 1
}

func getNumberOfEmptyCells(board [][]int) int {
	var emptyCells int = 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] == 0 {
				emptyCells++
			}
		}
	}
	return emptyCells
}

func calculateClusteringScore(board [][]int) int {
	//     var clusteringScore int = 0;
	//
	//     List<int> neighbors = [-1, 0, 1];
	//
	//     for (int i = 0; i < board.grid.length; ++i) {
	//       for (int j = 0; j < board.grid.length; ++j) {
	//         if (board.grid[i][j] == 0) {
	//           continue; //ignore empty cells
	//         }
	//
	//         //clusteringScore-=boardArray[i][j];
	//
	//         //for every pixel find the distance from each neightbors
	//         int numOfNeighbors = 0;
	//         int sum = 0;
	//         for (int k in neighbors) {
	//           int x = i + k;
	//           if (x < 0 || x >= board.grid.length) {
	//             continue;
	//           }
	//           for (int l in neighbors) {
	//             int y = j + l;
	//             if (y < 0 || y >= board.grid.length) {
	//               continue;
	//             }
	//
	//             if (board.grid[x][y] > 0) {
	//               ++numOfNeighbors;
	//               sum += (board.grid[i][j] - board.grid[x][y]).abs();
	//             }
	//           }
	//         }
	//
	//         clusteringScore += (sum / numOfNeighbors).round();
	//       }
	//     }
	//
	//     return clusteringScore;
	return 0
}
