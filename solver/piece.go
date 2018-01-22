package solver

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Piece struct {
	Name       string
	CleanLines string
	NumCells   int
	// BitmaskAtPosition has boardWidth * boardHeight values,
	// one for each possible board position.
	// The value at position `p` is the list of possible PieceBitMasks
	// that could make the piece occupy this position.
	BitmaskAtPosition [][]BoardBitmask
}

func (p *Piece) String() string {
	return fmt.Sprintf("{%s, %d cells}", p.Name, p.NumCells)
}

func NewPiece(name, piecesString string, boardWidth, boardHeight int) (*Piece, error) {
	piecesString = strings.Replace(piecesString, "\r", "", -1)
	piecesString = strings.Replace(piecesString, "\t", " ", -1)
	lines := strings.Split(piecesString, "\n")
	earliestAsterisk := -1
	numCells := 0
	for i, line := range lines {
		firstAsterisk := -1
		for i, c := range line {
			if c != ' ' && c != '*' {
				return nil, fmt.Errorf("invalid character: %c", c)
			}
			if c == '*' {
				if firstAsterisk == -1 {
					firstAsterisk = i
				}
				numCells++
			}
		}
		if firstAsterisk == -1 {
			return nil, fmt.Errorf("line %d has no asterisks: %q", i, line)
		}
		if earliestAsterisk == -1 || firstAsterisk < earliestAsterisk {
			earliestAsterisk = firstAsterisk
		}
	}
	if earliestAsterisk == -1 {
		return nil, errors.New("found no asterisks")
	}
	pieceWidth := 0
	for i, line := range lines {
		sanitizedLine := strings.TrimRight(line[earliestAsterisk:], " ")
		if len(sanitizedLine) > pieceWidth {
			pieceWidth = len(sanitizedLine)
		}
		lines[i] = sanitizedLine
	}
	pieceHeight := len(lines)
	// Got sanitized piece string at this point. Convert to bool array[j][i].
	pieceBool := make([][]bool, pieceHeight)
	for j, line := range lines {
		pieceLine := make([]bool, pieceWidth)
		for i, c := range line {
			pieceLine[i] = c == '*'
		}
		pieceBool[j] = pieceLine
	}
	// Create list of orientations, still in array[j][i] format.
	var orientationsBool [][][]bool
	orientationsBool = append(orientationsBool, pieceBool)
	for o := 0; o < 3; o++ {
		// Rotate pieceBool by 90 degrees.
		pieceBoolWidth := len(pieceBool[0])
		pieceBoolHeight := len(pieceBool)
		ninetyDegree := make([][]bool, pieceBoolWidth)
		for j := 0; j < pieceBoolWidth; j++ {
			ninetyDegreeRow := make([]bool, pieceBoolHeight)
			for i := 0; i < pieceBoolHeight; i++ {
				ninetyDegreeRow[i] = pieceBool[i][pieceBoolWidth-j-1]
			}
			ninetyDegree[j] = ninetyDegreeRow
		}
		// Compare to existing rotations to avoid creating duplicates.
		foundSame := false
		for _, prev := range orientationsBool {
			if reflect.DeepEqual(prev, ninetyDegree) {
				foundSame = true
				break
			}
		}
		if foundSame {
			break
		}
		orientationsBool = append(orientationsBool, ninetyDegree)
		pieceBool = ninetyDegree
	}
	// Filter impossible orientations.
	var possibleOrientationsBool [][][]bool
	for _, o := range orientationsBool {
		oHeight := len(o)
		oWidth := len(o[0])
		if oWidth > boardWidth || oHeight > boardHeight {
			continue
		}
		possibleOrientationsBool = append(possibleOrientationsBool, o)
	}
	if len(possibleOrientationsBool) == 0 {
		return nil, errors.New("piece is too large to fit on the board in any orientation")
	}
	// For each possible orientation, determine all the ways in which it could fit a
	// given board position (0, 0).
	// occupiedOffsets represents a list of board bit offsets.
	// It also contains the rows above and below position (0, 0) that it needs,
	// and the number of columns it occupies horizontally.
	type occupiedOffsets struct {
		needRowsAbove, needRowsBelow      int
		needColumnsLeft, needColumnsRight int
		occupiedOffsets                   []int
	}
	var possibleOccupiedBitOffsets []occupiedOffsets
	for _, o := range possibleOrientationsBool {
		oHeight := len(o)
		oWidth := len(o[0])
		for heightOffset := 0; heightOffset < oHeight; heightOffset++ {
			// Determine the required width offset for the given height offset.
			widthOffset := -1
			rowAtOffsetHeight := o[heightOffset]
			for i := 0; i < oWidth && widthOffset == -1; i++ {
				if rowAtOffsetHeight[i] {
					widthOffset = i
				}
			}
			if widthOffset == -1 {
				panic("got empty row")
			}
			oo := occupiedOffsets{
				needRowsAbove:    heightOffset,
				needRowsBelow:    oHeight - heightOffset - 1,
				needColumnsLeft:  widthOffset,
				needColumnsRight: oWidth - widthOffset - 1,
				occupiedOffsets:  make([]int, 0, numCells),
			}
			for j := 0; j < oHeight; j++ {
				for i := 0; i < oWidth; i++ {
					if o[j][i] {
						occupiedOffset := boardWidth*(j-heightOffset) + (i - widthOffset)
						oo.occupiedOffsets = append(oo.occupiedOffsets, occupiedOffset)
					}
				}
			}
			possibleOccupiedBitOffsets = append(possibleOccupiedBitOffsets, oo)
		}
	}
	// Got full list of possible occupied bit offsets.
	// Now create the list of bitmasks this represents at each board position where the piece fits.
	bitmaskAtPosition := make([][]BoardBitmask, boardWidth*boardHeight)
	for _, oo := range possibleOccupiedBitOffsets {
		for y := oo.needRowsAbove; y < boardHeight-oo.needRowsBelow; y++ {
			for x := oo.needColumnsLeft; x < boardWidth-oo.needColumnsRight; x++ {
				boardOffset := y*boardWidth + x
				boardPositions := make([]int, len(oo.occupiedOffsets))
				for i, occupiedOffset := range oo.occupiedOffsets {
					boardPositions[i] = occupiedOffset + boardOffset
				}
				bitmaskAtPosition[boardOffset] = append(bitmaskAtPosition[boardOffset], NewBoardBitmask(boardPositions))
			}
		}
	}
	return &Piece{
		Name:              name,
		CleanLines:        strings.Join(lines, "\n"),
		NumCells:          numCells,
		BitmaskAtPosition: bitmaskAtPosition,
	}, nil
}
