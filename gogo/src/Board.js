import React, { useState, useEffect } from "react";
import {
  makeBoard,
  calcSide,
  connReady,
  calculateLocalMove,
  getStoneColor,
} from "./util";
import { states, moves } from "./constants";
import "./webapp.css";

const [initBoard, boardTemplate] = makeBoard();

const Board = ({ socket, playerName, gameData }) => {
  const [gameState, setGameState] = useState(initBoard);
  const [ourStone, setOurStone] = useState(states.BLACK);
  const [stoneLocation, setStoneLocation] = useState("");
  const [finishedTurn, setFinishedTurn] = useState(false);
  const [godMode] = useState(false); //temp, put pieces wherever

  useEffect(() => {
    if (gameData) {
      const { board, nextPlayer } = gameData;
      !godMode && setGameState(board);
      console.log("huh", nextPlayer, playerName);
      setFinishedTurn(nextPlayer !== playerName);
      setOurStone(getStoneColor(gameData, playerName));
    }
  }, [gameData]);

  useEffect(() => {
    if (connReady(socket) && finishedTurn) {
      console.log("sending game");
      socket.send(
        JSON.stringify({
          gameId: "theonlygame",
          player: playerName,
          color: ourStone,
          move: moves.PLAY,
          point: stoneLocation,
          finishedTurn,
          boardTemp: gameState,
        })
      );
    }
  }, [socket, finishedTurn]);

  const setStone = ({ target: { id: selectedLocation }, detail }) => {
    if (finishedTurn) {
      console.log("not my turn.");
      return;
    }
    const [newBoard, newStoneLocation, isFinished] = calculateLocalMove(
      gameState,
      stoneLocation,
      selectedLocation,
      ourStone,
      detail,
      godMode
    );
    setStoneLocation(newStoneLocation);
    setGameState(newBoard);
    if (isFinished) {
      setFinishedTurn(true);
    }
  };

  return (
    <div id="board-container">
      <div id="playing-area">
        {boardTemplate.map((rows, y) =>
          rows.map((_, x) => (
            <div
              key={y + ":" + x}
              id={y + ":" + x}
              className={`playing-square 
                ${calcSide(y, x)} 
                ${gameState[y][x]} ${
                stoneLocation === y + ":" + x && "selected"
              }`}
              onClick={setStone}
            />
          ))
        )}
      </div>
    </div>
  );
};

export default Board;
