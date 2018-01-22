package solver

import (
	"fmt"
	"sort"
	"strings"
)

// Bit alignment divisor.
const boardAlignment = 64

// allSet has is an int64 with all bits set. Used to skip over full board cells.
const allSet = ^uint64(0)

// Board is an array of uint64's, meant to be interpreted as a single long bitfield.
// Unfortunately Go doesn't support arbitrary-size bitfields, so this is the next best thing.
// The board cell at 0-index row `i` and 0-index column `j` is located at bit `j*width + i`.
// A 1 set means that a piece occupies the cell. A 0 means no piece occupies the cell.
type Board []uint64

func (b Board) String() string {
	s := make([]string, len(b))
	for i, v := range b {
		s[i] = fmt.Sprintf("%d: %x", i, v)
	}
	return strings.Join(s, "\n")
}

// FirstFreeSpot finds the first free board position.
// If there are no free spots, returns -1.
func (b Board) FirstFreeSpot() int {
	for i, v := range b {
		if v^allSet == 0 {
			continue
		}
		for x := 0; x < boardAlignment; x++ {
			if v&1 == 0 {
				return i*boardAlignment + x
			}
			v = v >> 1
		}
		panic(fmt.Sprintf("%x (index %d) is not all set, yet could not find unset bit", b[i], i))
	}
	return -1
}

func NewBoard(width, height int) Board {
	boardSize := width * height
	roundedUp := boardSize / boardAlignment
	if boardSize%boardAlignment != 0 {
		roundedUp++
	}
	board := make([]uint64, roundedUp)
	// Flip the bits that are out of bounds.
	var lastValue uint64
	for b := uint(boardSize % boardAlignment); b < 64; b++ {
		lastValue |= 1 << b
	}
	board[roundedUp-1] = lastValue
	return board
}

type CellBitmask struct {
	// The index within the Board []uint64 for which the Mask applies.
	BoardIndex int
	// The bit mask to check.
	Mask uint64
}

type BoardBitmask []CellBitmask

// GridPositions returns a list of (x, y)-tuples that correspond to the bits
// set in the bitmask, assuming the given board width and height.
func (bitmask BoardBitmask) GridPositions(width, height int) [][]int {
	w := uint(width)
	var positions [][]int
	for _, cellBitmask := range bitmask {
		mask := cellBitmask.Mask
		for b := uint(0); b < boardAlignment; b++ {
			bit := uint64(1 << b)
			if mask&bit == bit {
				boardPosition := uint(cellBitmask.BoardIndex)*boardAlignment + b
				positions = append(positions, []int{
					int(boardPosition % w),
					int(boardPosition / w),
				})
			}
		}
	}
	return positions
}

// NewBoardBitmask returns a new BoardBitmask that contains a mask where each given
// position in `boardPositions` corresponds to a `1` bit.
func NewBoardBitmask(boardPositions []int) BoardBitmask {
	bitMaskMap := make(map[int]uint64, len(boardPositions))
	for _, boardPosition := range boardPositions {
		boardIndex := boardPosition / boardAlignment
		bitIndex := uint(boardPosition % boardAlignment)
		bitMaskMap[boardIndex] = bitMaskMap[boardIndex] | (1 << bitIndex)
	}
	sortedIndexes := make([]int, 0, len(bitMaskMap))
	for boardIndex := range bitMaskMap {
		sortedIndexes = append(sortedIndexes, boardIndex)
	}
	sort.Ints(sortedIndexes)
	boardBitMask := make([]CellBitmask, len(sortedIndexes))
	for i, boardIndex := range sortedIndexes {
		boardBitMask[i] = CellBitmask{
			BoardIndex: boardIndex,
			Mask:       bitMaskMap[boardIndex],
		}
	}
	return BoardBitmask(boardBitMask)
}

// Add adds the given BoardBitmask to the board.
// On success, returns a new Board with the bitmask applied.
// The original Board is not modified.
// On failure, returns nil.
func (b Board) Add(bitmask BoardBitmask) Board {
	for _, cellBitmask := range bitmask {
		if b[cellBitmask.BoardIndex]&cellBitmask.Mask != 0 {
			// Doesn't fit.
			return nil
		}
	}
	// Fits. Make new board.
	newBoard := make([]uint64, len(b))
	copy(newBoard, b)
	for _, cellBitmask := range bitmask {
		newBoard[cellBitmask.BoardIndex] |= cellBitmask.Mask
	}
	return newBoard
}
