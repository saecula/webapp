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
