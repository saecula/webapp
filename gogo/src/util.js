import { STATES } from "./constants";

export const makeBoard = () => {
  const board = {};

  const boardTemplate = Array(19)
    .fill(null)
    .map(() => Array(19).fill(null));

  boardTemplate.forEach((row, i) => {
    const rowKey = i.toString();
    board[rowKey] = {};
    row.forEach((_, y) => {
      const squareKey = y.toString();
      board[rowKey][squareKey] = STATES.empty;
    });
  });

  return [board, boardTemplate];
};
