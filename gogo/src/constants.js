export const states = {
  EMPTY: "e",
  BLACK: "b",
  WHITE: "w",
};

export const moves = {
  // switch colors, only valid before first turn
  SWITCH: "switch",
  // play a stone
  PLAY: "play",
  // pass your turn
  PASS: "pass",
  // resign the game
  RESIGN: "resign",
};
