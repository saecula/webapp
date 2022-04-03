import React, { useState, useEffect } from "react";
import axios from "axios";
import "./webapp.css";

const STATES = {
  empty: "e",
  black: "b",
  white: "w",
};

const board = Array(19)
  .fill(null)
  .map(() => Array(19).fill(null));

const boardMap = {};
board.forEach((row, i) => {
  const rowKey = i.toString();
  boardMap[rowKey] = {};
  row.forEach((_, y) => {
    const squareKey = y.toString();
    boardMap[rowKey][squareKey] = STATES.empty;
  });
});

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

  const [stoneColor] = useState(STATES.black);
  const [gameState, setGameState] = useState(boardMap);

  const setStone = (key) => {
    const [cRow, cCol] = key.split(":");
    const prevSquareState = gameState[cRow][cCol];
    const newGameState =
      prevSquareState === STATES.empty
        ? stoneColor
        : prevSquareState == stoneColor
        ? STATES.empty
        : prevSquareState;

    setGameState((prevGameState) => {
      const newthing = {
        ...prevGameState,
        [cRow]: { ...prevGameState[cRow], [cCol]: newGameState },
      };
      console.log("newthing", newthing);
      return newthing;
    });
  };

  const makePlayingSquare = (rowNum, colNum) => {
    const key = rowNum + ":" + colNum;
    return (
      <div
        key={key}
        id={key}
        className={`playing-square ${gameState[rowNum][colNum]}`}
        onClick={() => {
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
              rows.map((_, x) => makePlayingSquare(y.toString(), x.toString()))
            )}
          </div>
        </div>
      </header>
    </div>
  );
};

export default Webapp;
