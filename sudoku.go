package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const printOnSetSolved bool = false

type suduku struct {
	board [10][10]int
	data  [10][10][10]int
}

func newSuduku() suduku {
	var s suduku
	for r := 1; r < 10; r++ {
		for c := 1; c < 10; c++ {
			s.data[r][c] = [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		}
	}

	return s
}
func setSolved(s *suduku, r int, c int, v int, doprint bool) {
	if v == 0 {
		return
	}
	if doprint {
		fmt.Printf("Setting (%d,%d) to %d\n", r, c, v)
		printBoard(s)
		showLevel(s, v)
	}

	s.board[r][c] = v
	for i := 1; i < 10; i++ {
		s.data[r][c][i] = 0
		s.data[r][i][v] = 0
		s.data[i][c][v] = 0
	}
	rr := r - (r-1)%3
	cc := c - (c-1)%3
	for rrr := 0; rrr < 3; rrr++ {
		for ccc := 0; ccc < 3; ccc++ {
			s.data[rr+rrr][cc+ccc][v] = 0
		}
	}
	if doprint {
		fmt.Printf("Done Setting (%d,%d) to %d\n", r, c, v)
		printBoard(s)
		showLevel(s, v)
	}
}
func solveBoard(s *suduku) (suduku, bool) {
	found := true
	for found {
		found = false
		for r := 1; r < 10; r++ {
			for c := 1; c < 10; c++ {
				if s.board[r][c] == 0 {
					vv := -1
					for v := 1; v < 10; v++ {
						d := s.data[r][c][v]
						if d != 0 {
							switch vv {
							case -1:
								vv = v
							default:
								{
									vv = 0
									break
								}
							}
						}
					}
					switch vv {
					case -1:
						return *s, false
					case 0:
					default:
						{
							setSolved(s, r, c, vv, printOnSetSolved)
							found = true
						}
					}
				}
			}
			for v := 1; v < 10; v++ {
				for r := 1; r < 10; r++ {
					var rfound, cfound [][]int
					for c := 1; c < 10; c++ {
						if s.data[r][c][v] != 0 {
							rfound = append(rfound, []int{r, c})
						}
						if s.data[c][r][v] != 0 {
							cfound = append(cfound, []int{c, r})
						}
					}

					if len(rfound) == 1 {
						setSolved(s, rfound[0][0], rfound[0][1], v, printOnSetSolved)
						found = true
					}
					if len(cfound) == 1 {
						setSolved(s, cfound[0][0], cfound[0][1], v, printOnSetSolved)
						found = true
					}
				}
			}
			for br := 1; br < 10; br += 3 {
				for bc := 1; bc < 10; bc += 3 {
					for v := 1; v < 10; v++ {
						var bfound [][]int
						for r := 0; r < 3; r++ {
							for c := 0; c < 3; c++ {
								if s.data[br+r][bc+c][v] != 0 {
									bfound = append(bfound, []int{br + r, bc + c})
								}
							}
						}
						switch len(bfound) {
						case 1:
							{
								setSolved(s, bfound[0][0], bfound[0][1], v, printOnSetSolved)
								found = true
							}
						case 2:
							{
								if bfound[0][0] == bfound[1][0] {
									r := bfound[0][0]
									c0 := bfound[0][1]
									c1 := bfound[1][1]
									for c := 1; c < 10; c++ {
										if c != c0 && c != c1 && s.data[r][c][v] != 0 {
											s.data[r][c][v] = 0
											found = true
										}
									}
								}
								if bfound[0][1] == bfound[1][1] {
									c := bfound[0][1]
									r0 := bfound[0][0]
									r1 := bfound[1][0]
									for r := 1; r < 10; r++ {
										if r != r0 && r != r1 && s.data[r][c][v] != 0 {
											s.data[r][c][v] = 0
											found = true
										}
									}
								}
							}
						case 3:
							{
								if bfound[0][0] == bfound[1][0] && bfound[1][0] == bfound[2][0] {
									r := bfound[0][0]
									c0 := bfound[0][1]
									c1 := bfound[1][1]
									c2 := bfound[2][1]
									for c := 1; c < 10; c++ {
										if c != c0 && c != c1 && c != c2 && s.data[r][c][v] != 0 {
											s.data[r][c][v] = 0
											found = true
										}
									}
								}
								if bfound[0][1] == bfound[1][1] && bfound[1][1] == bfound[2][1] {
									c := bfound[0][1]
									r0 := bfound[0][0]
									r1 := bfound[1][0]
									r2 := bfound[2][0]
									for r := 1; r < 10; r++ {
										if r != r0 && r != r1 && r != r2 && s.data[r][c][v] != 0 {
											s.data[r][c][v] = 0
											found = true
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	if checkSolved(s) {
		return *s, true
	}
	guesses := findGuess(s)
	for i := range guesses {
		r := guesses[i][0]
		c := guesses[i][1]
		v := guesses[i][2]
		fmt.Printf("guessing (%d, %d) %d\n", r, c, v)
		news := *s
		setSolved(&news, r, c, v, printOnSetSolved)
		news, solved := solveBoard(&news)
		if solved {
			return news, true
		}
		fmt.Printf("guess (%d, %d) %d failed\n", r, c, v)
	}
	return *s, false
}
func findGuess(s *suduku) [][]int {
	for r := 1; r < 10; r++ {
		for c := 1; c < 10; c++ {
			if s.board[r][c] == 0 {
				var vs [][]int
				for v := 1; v < 10; v++ {
					if s.data[r][c][v] != 0 {
						vs = append(vs, []int{r, c, v})
					}
				}
				return vs
			}
		}
	}
	return nil
}

func checkSolved(s *suduku) bool {
	for r := 1; r < 10; r++ {
		for c := 1; c < 10; c++ {
			if s.board[r][c] == 0 {
				return false
			}
		}
	}
	return true
}
func readSuduku(filename string) (suduku, error) {
	s := newSuduku()
	fmt.Printf("reading file %s\n", filename)
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return s, err
	}
	datstr := strings.Split(string(dat), "")
	i := 0
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			for {
				if i >= len(datstr) {
					return s, fmt.Errorf("Could not read %s", filename)
				}
				v, err := strconv.Atoi(string(datstr[i]))
				i++
				if err == nil {
					if v != 0 {
						setSolved(&s, r+1, c+1, v, false)
					}
					break
				}
			}
		}
	}
	printBoard(&s)
	return s, nil
}
func printBoard(s *suduku) {
	for q := range s.board {
		if q == 0 {
			continue
		}
		fmt.Println(s.board[q][1:])
	}
	fmt.Println()
}
func showLevel(s *suduku, v int) {
	for r := 1; r < 10; r++ {
		for c := 1; c < 10; c++ {
			fmt.Printf("%d ", s.data[r][c][v])
		}
		fmt.Println()
	}
}
func usage() {
	fmt.Printf("usage: %s <filename>\n", os.Args[0])
	fmt.Println("where <filename> is the name of a file defining sudoku board")
}
func main() {
	if len(os.Args) != 2 {
		usage()
		return
	}
	filename := os.Args[1]
	s, err := readSuduku(filename)
	if err != nil {
		fmt.Println(err)
		usage()
		return
	}
	s, solved := solveBoard(&s)
	if !solved {
		fmt.Println()
		fmt.Printf("Could not Solve %s not a legitimate sudoku board\n", filename)
	} else {
		fmt.Println("Solution:")
	}
	fmt.Println()
	printBoard(&s)

}
