import React, { useState, useEffect, useCallback } from "react";
import axios from "axios";
import Board from "./Board";
import { PLAYER_NAME_LOCALSTORAGE } from "./constants";
import { getPlayerNames } from "./util";
import WhoAreUModal from "./WhoAreUModal";
import ErrorModal from "./ErrorModal";
import "./webapp.css";

const Webapp = () => {
  const [socket, setSocket] = useState(undefined);
  const [gameData, setGameData] = useState(null);
  const [loaded, setIsLoaded] = useState(false);
  const [playerName, setPlayerName] = useState(
    localStorage.getItem(PLAYER_NAME_LOCALSTORAGE)
  );
  const [fetchError, setFetchError] = useState(null);

  useEffect(() => {
    initSocket();
    loadGameState();
    return disconnectSocket;
  }, []);

  const loadGameState = useCallback(async () => {
    const id = window.location?.pathname?.slice(1);
    try {
      const { data } = await axios.get('/', id && { params: { id } });
      // determine if t sdhat id is valid on backend
      // if not send back 404...see how axios sends it
      // if on brand new game, be able to auto populate stored name
      // if (!getPlayerNames(data).includes(playerName)) {
      //   setPlayerName("");
      // }
      setGameData(data);
      setIsLoaded(true);
    } catch (err) {
      console.error("Errors loading game state", err);
      setFetchError(err);
    }
  }, [setPlayerName, setGameData, setIsLoaded, setFetchError]);

  const initSocket = useCallback(() => {
    const s = new WebSocket('ws://localhost:4000/ws');
    s.addEventListener("message", function ({ data }) {
      const parsedData = JSON.parse(data);
      console.log("got game data on socket:", parsedData);
      setGameData(parsedData);
    });
    setSocket(s);
  }, [setSocket]);

  const disconnectSocket = useCallback(() => {
    socket?.disconnect();
    setSocket(null);
  }, [socket, setSocket]);

  useEffect(() => {
    localStorage.setItem(PLAYER_NAME_LOCALSTORAGE, playerName);
  }, [playerName]);

  return (
    <div className="webapp">
      <header className="webapp-header">
        {fetchError && <ErrorModal err={fetchError} />}
        {!playerName && (
          <WhoAreUModal
            loaded={loaded}
            playerNames={getPlayerNames(gameData)}
            setPlayerName={setPlayerName}
          />
        )}
        <Board socket={socket} playerName={playerName} gameData={gameData} />
      </header>
    </div>
  );
};

export default Webapp;
