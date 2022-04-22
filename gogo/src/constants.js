export const SERVER_URL = `http://localhost:4000/`;
export const PLAYER_NAME_LOCALSTORAGE = "gogoname";

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
