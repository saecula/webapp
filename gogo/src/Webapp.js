import React, { useState, useEffect, useCallback } from "react";
import axios from "axios";
import { makeBoard, calcSide, connReady } from "./util";
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
    setSocket(s);
  }, [socket]);

  const disconnectSocket = () => {
    socket?.disconnect && socket.disconnect();
    setSocket(null);
  };

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

  const loadGameState = async () => {
    const me = localStorage.getItem("whoami");
    const url = `http://localhost:4000/`;
    try {
      const { data: retrievedGame } = await axios.get(url, {
        params: { id: me },
      });
      console.log("retrieved game:", retrievedGame);

      setGameState(retrievedGame.board);
      setIsMyTurn(retrievedGame.nextPlayer === me);
      const { b, w } = retrievedGame.players;
      const myStone = b === me ? b : w;
      console.log("setting my stone:", myStone);
      setStone(myStone);
    } catch (err) {
      console.error("Errors loading game state", err);
    }
  };

  const setStone = (selectedLocation) => {
    console.log("hi", stoneLocation, "and", selectedLocation);
    if (!isMyTurn) {
      console.log("not my turn.");
      return;
    }

    let hadBeenPlaced, oldRow, oldCol;
    const [curRow, curCol] = selectedLocation.split(":");
    if (stoneLocation) {
      hadBeenPlaced = true;
      oldRow = stoneLocation.split(":")[0];
      oldCol = stoneLocation.split(":")[1];
    }
    const prevPointState = gameState[curRow][curCol];

    if (selectedLocation === stoneLocation) {
      setStoneLocation("");
      setGameState((prevGameState) => ({
        ...prevGameState,
        [curRow]: { ...prevGameState[curRow], [curCol]: states.EMPTY },
      }));
    } else if (prevPointState === states.EMPTY) {
      setStoneLocation(selectedLocation);
      const boardWithoutOldLocation = hadBeenPlaced
        ? {
            ...gameState,
            [oldRow]: { ...gameState[oldRow], [oldCol]: states.EMPTY },
          }
        : gameState;
      const boardWithNewLocation = {
        ...boardWithoutOldLocation,
        [curRow]: { ...boardWithoutOldLocation[curRow], [curCol]: ourStone },
      };
      setGameState(boardWithNewLocation);
    }
  };

  useEffect(() => {
    socket?.addEventListener("message", function ({ data }) {
      const message = JSON.parse(data);
      console.log("got message!", message);
      if (typeof message.board === "object") {
        setGameState(message.board);
        const { b, w } = message.players;
        const myStone = b === localStorage.getItem("whoami") ? b : w;
        console.log("setting my stone:", myStone);
        setStone(myStone);
        setIsMyTurn(message.nextPlayer === myStone);
      }
    });
  }, [socket]);

  useEffect(() => {
    const [initBoard, initBoardTemplate] = makeBoard();
    setGameState(initBoard);
    console.log("init board:", initBoard);
    setBoardTemplate(initBoardTemplate);
    loadGameState();
  }, []);

  useEffect(() => {
    console.log("socket state:", socket?.readyState);
    console.log("about to send this board", gameState);
    if (connReady(socket)) {
      console.log("sending~ :3", {
        gameId: "theonlygame",
        playerId: localStorage.getItem("whoami"),
        color: ourStone,
        move: moves.PLAY,
        point: stoneLocation,
        finishedTurn: !isMyTurn,
        boardtemp: gameState,
      });
      socket.send(
        JSON.stringify({
          gameId: "theonlygame",
          playerId: localStorage.getItem("whoami"),
          color: ourStone,
          move: moves.PLAY,
          point: stoneLocation,
          finishedTurn: !isMyTurn,
          boardtemp: gameState,
        })
      );
    }
  }, [stoneLocation, socket]);

  useEffect(() => {
    console.log("game state changed:", gameState);
  }, [gameState]);

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
