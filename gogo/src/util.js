import { states } from "./constants";

export const makeBoard = () => {
  const board = {};

  const boardTemplate = Array(19)
    .fill(null)
    .map(() => Array(19).fill(null));

  boardTemplate.forEach((row, i) => {
    board[i] = {};
    row.forEach((_, y) => {
      board[i][y] = states.EMPTY;
    });
  });
  console.log("hmmm", board);
  return [board, boardTemplate];
};

export const connReady = (socket) => socket?.readyState === 1;

export const calcSide = (rowNum, colNum) => {
  let side =
    rowNum === 0
      ? "top"
      : rowNum === 18
      ? "bottom"
      : colNum === 0
      ? "left"
      : colNum === 18
      ? "right"
      : "mid";

  side =
    rowNum === 0 && colNum === 0
      ? "topleft"
      : rowNum === 18 && colNum === 18
      ? "bottomright"
      : rowNum === 0 && colNum === 18
      ? "topright"
      : rowNum === 18 && colNum === 0
      ? "bottomleft"
      : side;

  side =
    (rowNum === 3 && colNum === 3) ||
    (rowNum === 15 && colNum === 15) ||
    (rowNum === 3 && colNum === 15) ||
    (rowNum === 15 && colNum === 3) ||
    (rowNum === 9 && colNum === 3) ||
    (rowNum === 3 && colNum === 9) ||
    (rowNum === 9 && colNum === 9) ||
    (rowNum == 15 && colNum === 9) ||
    (rowNum == 9 && colNum === 15)
      ? "dot"
      : side;

  return side;
};

export const removeLastPlayed = (gameState, stoneLocation) => {
  if (stoneLocation) {
    const [row, col] = stoneLocation.split(":");
    gameState[row][col] = "e";
  }
  return gameState;
};

export const calculateLocalMove = (
  gameState,
  oldLocation,
  attemptedLocation,
  ourStone,
  numClicks,
  godMode = false
) => {
  let isFinished = false;
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

  if (godMode) {
    newStoneLocation = "";
    newGameState = {
      ...gameState,
      [curRow]: {
        ...gameState[curRow],
        [curCol]: prevPointState !== states.EMPTY ? states.EMPTY : ourStone,
      },
    };
  } else if (attemptedLocation === oldLocation) {
    console.log("numClicks....", numClicks);
    if (numClicks > 1) {
      return [newGameState, newStoneLocation, true];
    }
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
  return [newGameState, newStoneLocation, isFinished];
};

export const getStoneColor = (retrievedGame, playerName) => {
  const { b } = retrievedGame.players;
  return b === playerName ? states.BLACK : states.WHITE;
};

export const getPlayerNames = (gameData) =>
  Object.values(gameData?.players || []).filter((n) => !!n);

export const validateNameInput = (e) => {
  const input = e.target?.value || e.target[0]?.value;
  let alertMsg;
  if (!input) {
    alertMsg = "no value, hmm.";
  } else if (input.length > 32) {
    alertMsg = "pls pick shorter name";
  } else if (input.includes("<") || input.includes(";")) {
    alertMsg = "pls no weird symbol :3";
  }
  if (alertMsg) {
    alert(alertMsg);
  }
  return input;
};
