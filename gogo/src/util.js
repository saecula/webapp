import { states } from "./constants";

export const makeBoard = () => {
  const board = {};

  const boardTemplate = Array(19)
    .fill(null)
    .map(() => Array(19).fill(null));

  boardTemplate.forEach((row, i) => {
    const rowKey = i;
    board[rowKey] = {};
    row.forEach((_, y) => {
      const squareKey = y;
      board[rowKey][squareKey] = states.EMPTY;
    });
  });
  console.log("hmmm", board);
  return [board, boardTemplate];
};

export const connReady = (socket) => socket?.readyState === 1;

export const calcSide = (rowNum, colNum) => {
  let side =
    rowNum === "0"
      ? "top"
      : rowNum === "18"
      ? "bottom"
      : colNum === "0"
      ? "left"
      : colNum === "18"
      ? "right"
      : "mid";

  side =
    rowNum === "0" && colNum === "0"
      ? "topleft"
      : rowNum === "18" && colNum === "18"
      ? "bottomright"
      : rowNum === "0" && colNum === "18"
      ? "topright"
      : rowNum === "18" && colNum === "0"
      ? "bottomleft"
      : side;

  return side;
};

export const calculateLocalMove = (
  gameState,
  oldLocation,
  attemptedLocation,
  ourStone
) => {
  let newGameState = gameState;
  let newStoneLocation = oldLocation;

  let hadBeenPlaced, oldRow, oldCol;
  const [curRow, curCol] = attemptedLocation.split(":");
  if (oldLocation) {
    hadBeenPlaced = true;
    oldRow = oldLocation.split(":")[0];
    oldCol = oldLocation.split(":")[1];
  }
  const prevPointState = gameState[curRow][curCol];

  if (attemptedLocation === oldLocation) {
    newStoneLocation = "";
    newGameState = {
      ...gameState,
      [curRow]: { ...gameState[curRow], [curCol]: states.EMPTY },
    };
  } else if (prevPointState === states.EMPTY) {
    newStoneLocation = attemptedLocation;
    const boardWithoutOldLocation = hadBeenPlaced
      ? {
          ...gameState,
          [oldRow]: { ...gameState[oldRow], [oldCol]: states.EMPTY },
        }
      : gameState;
    const boardWithNewLocation = {
      ...boardWithoutOldLocation,
      [curRow]: { ...boardWithoutOldLocation[curRow], [curCol]: ourStone },
    };
    newGameState = boardWithNewLocation;
  }
  return [newGameState, newStoneLocation];
};

export const getStoneColor = (retrievedGame, playerName) => {
  const { b } = retrievedGame.players;
  return b === playerName ? states.BLACK : states.WHITE;
};
