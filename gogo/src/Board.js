import React, { useState, useEffect } from "react";
import {
  makeBoard,
  calcSide,
  connReady,
  calculateLocalMove,
  getStoneColor,
  removeLastPlayed,
} from "./util";
import { states, moves } from "./constants";
import { PassButton, ResignButton, SwitchButton } from "./Buttons";
import "./webapp.css";

const [initBoard, boardTemplate] = makeBoard();

const Board = ({ socket, playerName, gameData }) => {
  const [gameState, setGameState] = useState(initBoard);
  const [ourStone, setOurStone] = useState(states.BLACK);
  const [stoneLocation, setStoneLocation] = useState("");
  const [finishedTurn, setFinishedTurn] = useState(false);
  const [sentGame, setSentGame] = useState(false);
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
    console.log(
      "finished turn:",
      finishedTurn,
      ", game unsent:",
      !sentGame,
      "stone location:",
      stoneLocation
    );
    if (connReady(socket) && finishedTurn && !sentGame) {
      console.log("sending game", {
        id: "theonlygame",
        player: playerName,
        color: ourStone,
        move: moves.PLAY,
        point: stoneLocation,
        finishedTurn,
        boardTemp: gameState,
      });
      socket.send(
        JSON.stringify({
          id: "theonlygame",
          player: playerName,
          color: ourStone,
          move: moves.PLAY,
          point: stoneLocation,
          finishedTurn,
          boardTemp: gameState,
        })
      );
      setSentGame(true);
      setStoneLocation("");
    }
  }, [socket, finishedTurn, stoneLocation, playerName, ourStone, gameState]);

  useEffect(() => {
    if (!finishedTurn) {
      setSentGame(false);
    }
  }, [finishedTurn]);

  const setStone = ({
    target: { id: selectedLocation },
    detail: numClicks,
  }) => {
    if (finishedTurn) return;
    const [newBoard, newStoneLocation, isFinished] = calculateLocalMove(
      gameState,
      stoneLocation,
      selectedLocation,
      ourStone,
      numClicks,
      godMode
    );
    setStoneLocation(newStoneLocation);
    setGameState(newBoard);
    if (isFinished) {
      setFinishedTurn(true);
    }
  };

  const onPass = () => {
    setFinishedTurn(true);
    const board = removeLastPlayed(gameState, stoneLocation);
    if (connReady(socket)) {
      socket.send(
        JSON.stringify({
          id: "theonlygame",
          player: playerName,
          color: ourStone,
          move: moves.PASS,
          point: "",
          finishedTurn: true,
          boardTemp: board,
        })
      );
    }
  };

  return (
    <div
      style={{ display: "flex", justifyContent: "space-evenly", width: "100%" }}
    >
      <div
        id="board"
        className={`board-${
          finishedTurn ? "notmyturn" : "myturn"
        } board-${ourStone}`}
      >
        <div id="playing-area">
          {boardTemplate.map((rows, y) =>
            rows.map((_, x) => (
              <div
                key={y + ":" + x}
                id={y + ":" + x}
                className={`playing-square 
                ${calcSide(y, x)} 
                ${
                  gameState && gameState[y] && gameState[y][x] !== "e"
                    ? "stone"
                    : ""
                } ${gameState && gameState[y] && gameState[y][x]}${
                  stoneLocation === y + ":" + x ? "-selected" : "-unselected"
                }`}
                onClick={setStone}
              />
            ))
          )}
        </div>
      </div>
      <div style={{ display: "flex", flexDirection: "column", width: "30%" }}>
        <div className="button-container" style={{height: '30%'}}>
          <PassButton onPass={onPass} disabled={finishedTurn} />
          <ResignButton onResign={() => console.log("resign clicked.")} />
        </div>
        <div className="empty" />
        <div className="stonecup"><div className={`${ourStone}-stonesincup`}/></div>
      </div>
    </div>
  );
};

export default Board;
