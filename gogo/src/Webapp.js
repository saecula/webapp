import React, { useState, useEffect, useCallback } from "react";
import axios from "axios";
import { makeBoard } from "./util";
import { STATES } from "./constants";
import "./webapp.css";

const [board, boardTemplate] = makeBoard();

const Webapp = () => {
  const [socket, setSocket] = useState(null);

  useEffect(() => {
    initSocket();
    return disconnectSocket;
  }, []);

  const initSocket = useCallback(() => {
    const s = new WebSocket(`ws://localhost:4000/ws`);
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
  const [ourStone] = useState(STATES.black);
  const [isMyTurn, setIsMyTurn] = useState(true);
  const [finishedTurn, setFinishedTurn] = useState(false);
  const [gameState, setGameState] = useState(board);

  const setOurStone = (pointKey) => {
    if (!isMyTurn) {
      alert("not your turn :3");
      return;
    }
    const [cRow, cCol] = pointKey.split(":");
    const prevPlaceState = gameState[cRow][cCol];
    const newPlaceState =
      prevPlaceState === STATES.empty
        ? ourStone
        : prevPlaceState == ourStone
        ? STATES.empty
        : prevPlaceState;

    if (newPlaceState === prevPlaceState) {
      return;
    }
    setGameState((prevGameState) => ({
      ...prevGameState,
      [cRow]: { ...prevGameState[cRow], [cCol]: newPlaceState },
    }));
  };

  useEffect(() => {
    loadGameState();
  }, []);

  useEffect(() => {
    socket?.send(
      JSON.stringify({ game: gameState, played: ourStone, finishedTurn })
    );
  }, [gameState]);

  const loadGameState = async () => {
    try {
      const { data: retrievedGame } = await axios.get("http://localhost:4000");
      console.log("retrieved", retrievedGame);
      if (retrievedGame.game !== "hello:3") {
        setGameState(retrievedGame.game);
        setIsMyTurn(retrievedGame.next === ourStone);
      }
    } catch (err) {
      console.error("Error loading game state", err);
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
              setOurStone={setOurStone}
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

const PlayingSquare = ({ rowNum, colNum, setOurStone, gameState }) => {
  const key = rowNum + ":" + colNum;
  let side =
    rowNum === "0"
      ? "top"
      : rowNum === "18"
      ? "bottom"
      : colNum === "0"
      ? "left"
      : colNum === "18"
      ? "right"
      : "mid";

  side =
    rowNum == "0" && colNum == "0"
      ? "topleft"
      : rowNum == "18" && colNum == "18"
      ? "bottomright"
      : rowNum == "0" && colNum == "18"
      ? "topright"
      : rowNum == "18" && colNum == "0"
      ? "bottomleft"
      : side;

  return (
    <>
      <div
        key={key}
        id={key}
        className={`playing-square ${side} ${gameState[rowNum][colNum]}`}
        onClick={() => setOurStone(key)}
      >
        {/* {rowNum !== "0" && colNum !== "0" && <div className="visible-square" />} */}
      </div>
    </>
  );
};

export default Webapp;
