import React, { useState, useEffect } from "react";
import {
  makeBoard,
  calcSide,
  connReady,
  calculateLocalMove,
  getStoneColor,
  removeLastPlayed,
} from "./util";
import { moves } from "./constants";
import { PassButton, ResignButton, SwitchButton } from "./Buttons";
import "./webapp.css";

const [initBoard, boardTemplate] = makeBoard();

const Board = ({ socket, playerName, gameData, ourStone, setOurStone }) => {
  console.log('board player name', playerName)
  const [gameState, setGameState] = useState(initBoard);

  const [stoneLocation, setStoneLocation] = useState("");
  const [finishedTurn, setFinishedTurn] = useState(false);
  const [sentGame, setSentGame] = useState(false);
  const [godMode] = useState(false); //temp, put pieces wherever

  useEffect(() => {
    if (gameData) {
      const { board, nextPlayer } = gameData;
      !godMode && setGameState(board);
      setFinishedTurn(nextPlayer !== playerName);
      setOurStone(getStoneColor(gameData, playerName));
    }
  }, [gameData, playerName]);

  useEffect(() => {
    if (connReady(socket) && finishedTurn && !sentGame) {
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

    const onResign = () => {
      const wants = confirm('are you sure?')
      if (!wants) return

      setFinishedTurn(true);
      const board = removeLastPlayed(gameState, stoneLocation);
      if (connReady(socket)) {
        socket.send(
          JSON.stringify({
            id: "theonlygame",
            player: playerName,
            color: ourStone,
            move: moves.RESIGN,
            point: "",
            finishedTurn: true,
            boardTemp:  null,
          })
        );
      }
    };

        const onSwitch = () => {
          if (connReady(socket)) {
            socket.send(
              JSON.stringify({
                id: "theonlygame",
                player: playerName,
                color: ourStone,
                move: moves.SWITCH,
                point: "",
                finishedTurn: true,
                boardTemp: null,
              })
            );
          }
        };

  return (
    <div className="content-container">
      <div className="buttons-and-cup" />
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
      <div className="buttons-and-cup">
        <div className="button-container">
          <PassButton onPass={onPass} disabled={finishedTurn} />
          <ResignButton onResign={onResign} />
        </div>
        <div className="empty" />
        <div
          className={`stonecup ${!gameData.started ? "switchable" : ""}`}
          onClick={onSwitch}
        >
          <div className={`${ourStone}-stonesincup`} />
        </div>
      </div>
    </div>
  );
};

export default Board;
