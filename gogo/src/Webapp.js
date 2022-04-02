import React, { useState, useEffect } from "react";
import axios from "axios";
import "./webapp.css";

const BOARDSTATES = {
  empty: "empty",
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
        state: BOARDSTATES.empty,
      }))
  );
console.log(board);

const visibleBoard = Array(18)
  .fill(null)
  .map(() => Array(18).fill(null));

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

  const makePlayingSquare = (rowNum, colNum) => {
    console.log("rownum", rowNum);
    return <div id={`${rowNum}:${colNum}`} class={"playing-square"}></div>;
  };

  const [gameState, setGameState] = useState(board);
  return (
    <div className="webapp">
      <header className="webapp-header">
        <div id="board-container">
          <div id="playing-area">
            {board.map((rows, rowNum) =>
              rows.map((_, colNum) => makePlayingSquare(rowNum, colNum))
            )}
          </div>
          <div id="visible-board">
            {visibleBoard.map((r) =>
              r.map((c) => <div className="board-square" />)
            )}
          </div>
        </div>
      </header>
    </div>
  );
};

export default Webapp;
