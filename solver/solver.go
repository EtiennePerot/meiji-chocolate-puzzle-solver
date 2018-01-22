package solver

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type PieceOnBoard struct {
	Piece        *Piece
	BoardBitmask BoardBitmask
}

type Solution struct {
	Solver        *Solver
	Board         Board
	FirstFreeSpot int
	PiecesSet     []PieceOnBoard
	PiecesLeft    []*Piece
}

var (
	boxLeft       = 1 << 0
	boxRight      = 1 << 1
	boxUp         = 1 << 2
	boxDown       = 1 << 3
	boxHorizontal = "─"
	boxVertical   = "│"
	boxDrawing    = map[int]string{
		boxLeft + boxRight + boxUp + boxDown: "┼",
		boxLeft + boxRight + boxUp + 0:       "┴",
		boxLeft + boxRight + 0 + boxDown:     "┬",
		boxLeft + boxRight + 0 + 0:           boxHorizontal,
		boxLeft + 0 + boxUp + boxDown:        "┤",
		boxLeft + 0 + boxUp + 0:              "┘",
		boxLeft + 0 + 0 + boxDown:            "┐",
		boxLeft + 0 + 0 + 0:                  boxHorizontal,
		0 + boxRight + boxUp + boxDown:       "├",
		0 + boxRight + boxUp + 0:             "└",
		0 + boxRight + 0 + boxDown:           "┌",
		0 + boxRight + 0 + 0:                 boxHorizontal,
		0 + 0 + boxUp + boxDown:              boxVertical,
		0 + 0 + boxUp + 0:                    boxVertical,
		0 + 0 + 0 + boxDown:                  boxVertical,
		0 + 0 + 0 + 0:                        " ",
	}
)

func (s Solution) Print() {
	width := s.Solver.Width
	height := s.Solver.Height
	lines := make([][]string, height*2+1)
	for y := 0; y < height*2+1; y++ {
		line := make([]string, width*2+1)
		for x := 0; x < width*2+1; x++ {
			line[x] = " "
		}
		lines[y] = line
	}
	setBoardCell := func(x, y int, piece *Piece) {
		lines[y*2+1][x*2+1] = piece.Name
	}
	setBoardCellEdge := func(x, y, xEdgeOffset, yEdgeOffset int, edge string) {
		current := lines[y*2+1][x*2+1]
		writeEdge := true
		x2 := x + xEdgeOffset
		y2 := y + yEdgeOffset
		if x2 >= 0 && x2 < width && y2 >= 0 && y2 < height {
			writeEdge = current != lines[y2*2+1][x2*2+1]
		}
		if writeEdge {
			lines[y*2+1+yEdgeOffset][x*2+1+xEdgeOffset] = edge
		} else {
			lines[y*2+1+yEdgeOffset][x*2+1+xEdgeOffset] = " "
		}
	}

	// Done with main grid; draw individual pieces.
	for _, pieceSet := range s.PiecesSet {
		positions := pieceSet.BoardBitmask.GridPositions(width, height)
		for _, position := range positions {
			x, y := position[0], position[1]
			setBoardCell(x, y, pieceSet.Piece)
		}
		for _, position := range positions {
			x, y := position[0], position[1]
			setBoardCellEdge(x, y, -1, 0, boxVertical)
			setBoardCellEdge(x, y, 1, 0, boxVertical)
			setBoardCellEdge(x, y, 0, -1, boxHorizontal)
			setBoardCellEdge(x, y, 0, 1, boxHorizontal)
		}
	}

	// Write characters in the corners
	for y := 0; y <= height; y++ {
		for x := 0; x <= width; x++ {
			box := 0
			if x > 0 && lines[y*2][x*2-1] != " " {
				box += boxLeft
			}
			if x < width && lines[y*2][x*2+1] != " " {
				box += boxRight
			}
			if y > 0 && lines[y*2-1][x*2] != " " {
				box += boxUp
			}
			if y < height && lines[y*2+1][x*2] != " " {
				box += boxDown
			}
			lines[y*2][x*2] = boxDrawing[box]
		}
	}

	// All done; recompute a single string so that it can be printed atomically.
	combined := make([]string, height*2+1)
	for i, line := range lines {
		combined[i] = strings.Join(line, "")
	}
	piecesLeft := ""
	if len(s.PiecesLeft) != 0 {
		piecesLeft = fmt.Sprintf("\nPieces left: %v", s.PiecesLeft)
	}
	freeSpot := ""
	if s.FirstFreeSpot != -1 {
		freeSpot = fmt.Sprintf("\nFree spot: x=%d / y=%d", s.FirstFreeSpot%width, s.FirstFreeSpot/width)
	}
	fmt.Printf("%s%s%s\n", strings.Join(combined, "\n"), piecesLeft, freeSpot)
}

type Solver struct {
	Width, Height int
	Pieces        []*Piece
}

const pieceNames = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// New creates a new solver.
// `width` and `height` refer to the board size.
// `pieces` is a string that contains a visual representation of pieces using asterisks.
// Each piece should be separate by a whitespace-only line.
// Example valid `pieces` string:
//     `
//     ****
//     *
//
//     **
//     *
//
//     ***
//     * *
//
//     ***
//     `
func New(width, height int, piecesString string) (*Solver, error) {
	piecesSplit := strings.Split(piecesString, "\n")
	var pieces []*Piece
	var totalPieceSize int
	pieceNameIndex := 0
	for i := 0; i < len(piecesSplit); i++ {
		var currentPiece []string
		for j := i; j < len(piecesSplit); j++ {
			currentLine := piecesSplit[j]
			if strings.TrimSpace(currentLine) == "" {
				break
			}
			currentPiece = append(currentPiece, currentLine)
		}
		if len(currentPiece) != 0 {
			currentJoined := strings.Join(currentPiece, "\n")
			piece, err := NewPiece(pieceNames[pieceNameIndex:pieceNameIndex+1], currentJoined, width, height)
			if err != nil {
				return nil, fmt.Errorf("cannot parse piece: %v\n%v", err, currentJoined)
			}
			totalPieceSize += piece.NumCells
			pieces = append(pieces, piece)
			pieceNameIndex++
			if pieceNameIndex >= len(pieceNames) {
				return nil, errors.New("too many pieces; cannot give them unique names")
			}
			i += len(currentPiece) - 1
		}
	}
	if totalPieceSize != width*height {
		return nil, fmt.Errorf("pieces occupy a total of %d cells %v, but %d x %d board has %d cells", totalPieceSize, pieces, width, height, width*height)
	}
	return &Solver{width, height, pieces}, nil
}

func processSolution(sol Solution, work, found chan Solution) {
	for i, left := range sol.PiecesLeft {
		for _, bitmask := range left.BitmaskAtPosition[sol.FirstFreeSpot] {
			newBoard := sol.Board.Add(bitmask)
			if newBoard == nil {
				continue
			}
			newPiecesSet := make([]PieceOnBoard, len(sol.PiecesSet)+1)
			copy(newPiecesSet, sol.PiecesSet)
			newPiecesSet[len(sol.PiecesSet)] = PieceOnBoard{
				Piece:        left,
				BoardBitmask: bitmask,
			}
			newPiecesLeft := make([]*Piece, len(sol.PiecesLeft)-1)
			copy(newPiecesLeft[0:i], sol.PiecesLeft[0:i])
			copy(newPiecesLeft[i:len(newPiecesLeft)], sol.PiecesLeft[i+1:len(sol.PiecesLeft)])
			nextFreeSpot := newBoard.FirstFreeSpot()
			newSol := Solution{
				Solver:        sol.Solver,
				Board:         newBoard,
				FirstFreeSpot: nextFreeSpot,
				PiecesSet:     newPiecesSet,
				PiecesLeft:    newPiecesLeft,
			}
			if nextFreeSpot == -1 {
				select {
				case found <- newSol:
				default:
				}
				return
			}
			work <- newSol
		}
	}
}

func (s *Solver) Solve() (Solution, error) {
	work := make(chan Solution, 1<<24)
	found := make(chan Solution)
	for i := 0; i < runtime.NumCPU()*2; i++ {
		go func() {
			for sol := range work {
				processSolution(sol, work, found)
			}
		}()
	}
	work <- Solution{
		Solver:     s,
		Board:      NewBoard(s.Width, s.Height),
		PiecesLeft: s.Pieces,
	}
	return <-found, nil
}
