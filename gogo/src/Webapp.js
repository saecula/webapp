import React, { useState, useEffect } from "react";
import axios from "axios";
import "./webapp.css";

const BOARD_STATES = {
  empty: "empty",
  black: "black",
  white: "white",
};
const STONE_COLORS = {
  black: "black",
  white: "white",
};

const board = Array(19)
  .fill(null)
  .map((_, rowNum) =>
    Array(19)
      .fill(rowNum)
      .map((r, c) => ({
        row: r,
        col: c,
        state: BOARD_STATES.empty,
      }))
  );

const Webapp = () => {
  useEffect(() => {
    loadGameState();
  }, []);

  const loadGameState = async () => {
    try {
      // load game based on cookies ..... ?
      const retrievedGameState = await axios.get("http://localhost:4000");
      console.log("yey talked to back end", retrievedGameState);
    } catch (err) {
      console.error("Error loading game state", err);
    }
  };

  const [stoneColor] = useState(STONE_COLORS.black);
  const [gameState, setGameState] = useState(board);

  const setStone = (key) => {
    const [cRow, cCol] = key.split(":").map((s) => parseInt(s));
    console.log("current row and column:", cRow, cCol);
    setGameState((prevGameState) => {
      const prevState = prevGameState && prevGameState[cRow][cCol].state;
      console.log("prev state:", prevState);
      const newthing = prevGameState.map((prevRow) =>
        prevRow.map((sq) => {
          console.log(
            "trying to compare",
            sq.row,
            "to",
            cRow,
            "and",
            sq.col,
            "to",
            cCol
          );
          if (sq.row === cRow && sq.col === cCol) {
            console.log("FOUND");
            return {
              ...sq,
              state:
                prevState === BOARD_STATES.empty
                  ? stoneColor
                  : BOARD_STATES.empty,
            };
          } else {
            return { ...sq };
          }
        })
      );
      console.log("hm", newthing);
      return newthing;
    });
  };

  const makePlayingSquare = (rowNum, colNum) => {
    const key = rowNum + ":" + colNum;
    return (
      <div
        key={key}
        id={key}
        className={`playing-square ${
          gameState && gameState[rowNum] && gameState[rowNum][colNum]?.state
        }`}
        onClick={() => {
          console.log("hello", key);
          setStone(key);
        }}
      >
        {/* {rowNum != 0 && colNum != 0 && <div className="visible-square" />} */}
      </div>
    );
  };

  return (
    <div className="webapp">
      <header className="webapp-header">
        <div id="board-container">
          <div id="playing-area">
            {board.map((rows, y) =>
              rows.map((_, x) => makePlayingSquare(y, x))
            )}
          </div>
        </div>
      </header>
    </div>
  );
};

export default Webapp;
