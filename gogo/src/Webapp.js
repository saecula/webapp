import React, { useState, useEffect, useCallback } from "react";
import axios from "axios";
import Board from "./Board";
import { PLAYER_NAME_LOCALSTORAGE, moves, states } from "./constants";
import { getPlayerNames, connReady } from "./util";
import WhoAreUModal from "./WhoAreUModal";
import ErrorModal from "./ErrorModal";
import "./webapp.css";

const Webapp = () => {
  const [socket, setSocket] = useState(undefined);
  const [gameData, setGameData] = useState({});
  const [ourStone, setOurStone] = useState(states.BLACK);
  const [loaded, setIsLoaded] = useState(false);
  const [playerName, setPlayerName] = useState(
    localStorage.getItem(PLAYER_NAME_LOCALSTORAGE)
  );
  const [fetchError, setFetchError] = useState(null);
  const [waitingForPlayer2, setWaitingForPlayer2] = useState(false)

  useEffect(() => {
    initSocket();
    loadGameState();
    return disconnectSocket;
  }, []);

  const loadGameState = useCallback(async () => {
    const id = window.location?.pathname?.slice(1);
    try {
      const { data } = await axios.get('http://localhost:4000/', id && { params: { id } });
      // determine if t sdhat id is valid on backend
      // if not send back 404...see how axios sends it
      // if on brand new game, be able to auto populate stored name
      // if (!getPlayerNames(data).includes(playerName)) {
      //   setPlayerName("");
      // }
      console.log('data???', data)
      setGameData(data);
    } catch (err) {
      console.error("Errors loading game state", err);
      setFetchError(err);
    }
  }, []);

  useEffect(() => {
    if (!loaded && Object.keys(gameData).length > 1) { 
    console.log('gamedata', gameData); 
  }
  if (!waitingForPlayer2 && getPlayerNames(gameData).includes(playerName) && getPlayerNames(gameData).length < 2) {
    setWaitingForPlayer2(true);
  }
  setIsLoaded(true); 
}, [gameData])

  const initSocket = useCallback(() => {
    const s = new WebSocket("ws://localhost:4000/ws");//("ws://143.198.127.101:4000/ws");
    s.addEventListener("message", function ({ data }) {
      const parsedData = JSON.parse(data);
      console.log("got game data on socket:", parsedData);
      setGameData(parsedData);
    });
    setSocket(s);
  }, []);

  const disconnectSocket = useCallback(() => {
    if (socket) {
      console.log('wtf socket?', socket)
      socket.disconnect();
      setSocket(null);
    }
  }, [socket]);

  const onSubmitName = useCallback(
    (name, color) => {
      setPlayerName(name);
      if (color) setOurStone(color)
      console.log("sending?", socket, connReady(socket));
      if (connReady(socket)) {
        console.log("sending");
        socket.send(
          JSON.stringify({
            id: "theonlygame",
            player: name,
            color,
            move: moves.NAME,
            point: "",
            finishedTurn: true,
            boardTemp: {},
          })
        );
      }
      setWaitingForPlayer2(true)
    },
    [socket]
  );

  useEffect(() => {
    localStorage.setItem(PLAYER_NAME_LOCALSTORAGE, playerName);
  }, [playerName]);

  return (
    <div className="webapp">
      <header className="webapp-header">
        {gameData?.ended && (
          <div
            style={{
              position: "absolute",
              height: "200vh",
              width: "200vw",
              top: 0,
              left: 0,
              zIndex: 3,
              backgroundColor: "#5d5d5d9e",
            }}
          >
            <div className="modal-container">
              <div className="modal">
                <div style={{ margin: "auto" }}>halo</div>
                <button
                  style={{ width: "200px", margin: "auto" }}
                  onClick={() => axios.post("http://localhost:4000/newgame")}
                >
                  new game
                </button>
              </div>
            </div>
          </div>
        )}
        {fetchError && <ErrorModal err={fetchError} />}
        {loaded && !getPlayerNames(gameData).includes(playerName) && (
          <WhoAreUModal
            loaded={loaded}
            playerNames={getPlayerNames(gameData)}
            setPlayerName={onSubmitName}
          />
        )}
        {getPlayerNames(gameData).includes(playerName) &&
          getPlayerNames(gameData).length == 1 && <Waiting />}
        <Board
          socket={socket}
          playerName={playerName}
          gameData={gameData}
          ourStone={ourStone}
          setOurStone={setOurStone}
        />
      </header>
    </div>
  );
};

const Waiting = () =>       ( 
  <div className="modal-container">
              <div className="modal">
                <div style={{ margin: "auto" }}>Waiting for partner...</div>
              </div>
            </div>
  )


export default Webapp;
