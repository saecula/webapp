import React, { useState, useEffect, useCallback } from "react";
import axios from "axios";
import { makeBoard, calcSide } from "./util";
import { states, moves } from "./constants";
import "./webapp.css";

const Webapp = () => {
  const [socket, setSocket] = useState(null);

  useEffect(() => {
    initSocket();
    return disconnectSocket;
  }, []);

  const initSocket = useCallback(() => {
    const s = new WebSocket(`ws://localhost:4000/ws`);
    s.addEventListener("message", function ({ data }) {
      const message = JSON.parse(data);
      console.log("message!", message);
    });
    setSocket(s);
  }, [setSocket]);

  const disconnectSocket = () => socket?.disconnect && socket.disconnect();

  return (
    <div className="webapp">
      <header className="webapp-header">
        <Board socket={socket} />
      </header>
    </div>
  );
};

const Board = ({ socket }) => {
  const [gameState, setGameState] = useState({});
  const [boardTemplate, setBoardTemplate] = useState([]);

  const [isMyTurn, setIsMyTurn] = useState(true);

  const [ourStone] = useState(states.BLACK);
  const [stoneLocation, setStoneLocation] = useState("");

  const setStone = (selectedLocation) => {
    console.log("hi", stoneLocation, "and", selectedLocation);
    if (!isMyTurn || (stoneLocation && selectedLocation !== stoneLocation)) {
      return;
    }
    console.log("hello", selectedLocation);
    const [curRow, curCol] = selectedLocation.split(":");
    const [oldRow, oldCol] = stoneLocation.split(":");
    const prevPointState = gameState[curRow][curCol];

    if (selectedLocation === stoneLocation) {
      console.log("yes here");
      setStoneLocation("");
      setGameState((prevGameState) => ({
        ...prevGameState,
        [curRow]: { ...prevGameState[curRow], [curCol]: states.EMPTY },
      }));
    } else if (prevPointState === states.EMPTY) {
      setStoneLocation(selectedLocation);
      setGameState((prevGameState) => ({
        ...prevGameState,
        [oldRow]: { ...prevGameState[oldRow], [oldCol]: states.EMPTY },
        [curRow]: { ...prevGameState[curRow], [curCol]: ourStone },
      }));
    }
  };

  useEffect(() => {
    const [initBoard, initBoardTemplate] = makeBoard();
    setGameState(initBoard);
    setBoardTemplate(initBoardTemplate);
    loadGameState();
  }, []);

  useEffect(() => {
    console.log("sendsinggs :3", {
      gameId: localStorage.getItem("gogameid"),
      playerId: ourStone,
      move: moves.PLAY,
      point: stoneLocation,
      finishedTurn: !isMyTurn,
      boardtemp: gameState,
    });
    socket?.send(
      JSON.stringify({
        gameId: localStorage.getItem("gogameid"),
        playerId: ourStone,
        move: moves.PLAY,
        point: stoneLocation,
        finishedTurn: !isMyTurn,
      })
    );
  }, [stoneLocation]);

  useEffect(() => {
    console.log("hm", gameState);
  }, [gameState]);

  const loadGameState = async () => {
    const gameId = localStorage.getItem("gogameid");
    console.log("fetchisng Gos game id", gameId);
    const url = `http://localhost:4000/`;
    try {
      const { data: retrievedGame } = await axios.get(url);
      localStorage.setItem("gogameid", retrievedGame.id);
      console.log("retrievaed", retrievedGame);
      if (retrievedGame.gameId !== "new") {
        // setGameState(retrievedGame.board);
        setIsMyTurn(retrievedGame.nextPlayer === ourStone);
      }
    } catch (err) {
      console.error("Errors loading game state", err);
    }
  };

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

export default Webapp;
