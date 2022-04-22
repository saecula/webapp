import React, { useState, useEffect } from "react";
import axios from "axios";
import {
  makeBoard,
  calcSide,
  connReady,
  calculateLocalMove,
  getStoneColor,
} from "./util";
import {
  states,
  moves,
  SERVER_URL,
  PLAYER_NAME_LOCALSTORAGE,
} from "./constants";
import "./webapp.css";

const [initBoard, boardTemplate] = makeBoard();

const Board = ({ socket }) => {
  const [gameState, setGameState] = useState(initBoard);

  const [playerName, setPlayerName] = useState(
    localStorage.getItem(PLAYER_NAME_LOCALSTORAGE)
  );

  const [ourStone, setOurStone] = useState(states.BLACK);
  const [stoneLocation, setStoneLocation] = useState("");
  const [isMyTurn, setIsMyTurn] = useState(true);
  const loadGameState = async () => {
    try {
      const { data: retrievedGame } = await axios.get(SERVER_URL, {
        params: { id: playerName },
      });

      console.log("retrieved game:", retrievedGame);
      setGameState(retrievedGame.board);

      const myStoneColor = getStoneColor(retrievedGame, playerName);
      setOurStone(myStoneColor);
      setIsMyTurn(retrievedGame.nextPlayer === playerName);
    } catch (err) {
      console.error("Errors loading game state", err);
    }
  };

  const setStone = (selectedLocation) => {
    if (!isMyTurn) {
      console.log("not my turn.");
      return;
    }
    const [newBoard, newStoneLocation] = calculateLocalMove(
      gameState,
      stoneLocation,
      selectedLocation,
      ourStone
    );
    setStoneLocation(newStoneLocation);
    setGameState(newBoard);
  };

  useEffect(() => {
    socket?.addEventListener("message", function ({ data }) {
      const message = JSON.parse(data);
      console.log("got message!", message);
    });
  }, [socket]);

  useEffect(() => {
    loadGameState();
  }, []);

  useEffect(() => {
    if (connReady(socket)) {
      console.log("sending~ (not really)", {
        gameId: "theonlygame",
        player: playerName,
        color: ourStone,
        move: moves.PLAY,
        point: stoneLocation,
        finishedTurn: !isMyTurn,
        boardTemp: gameState,
      });
      //   socket.send(
      //     JSON.stringify({
      //         gameId: "theonlygame",
      //         player: playerName,
      //         color: ourStone,
      //         move: moves.PLAY,
      //         point: stoneLocation,
      //         finishedTurn: !isMyTurn,
      //         boardTemp: gameState,
      //     })
      //   );
    }
  }, [stoneLocation, socket]);

  return (
    <div id="board-container">
      <div id="playing-area">
        {boardTemplate.map((rows, y) =>
          rows.map((_, x) => (
            <PlayingSquare
              rowNum={y.toString()}
              colNum={x.toString()}
              setStone={setStone}
              gameState={gameState}
            >
              <hr style={{ color: "white", width: "1px" }} />
            </PlayingSquare>
          ))
        )}
      </div>
    </div>
  );
};

const PlayingSquare = ({ rowNum, colNum, setStone, gameState }) => {
  const key = rowNum + ":" + colNum;
  const side = calcSide(rowNum, colNum);
  return gameState ? (
    <>
      <div
        key={key}
        id={key}
        className={`playing-square ${side} ${gameState[rowNum][colNum]}`}
        onClick={() => setStone(key)}
      ></div>
    </>
  ) : null;
};

export default Board;
