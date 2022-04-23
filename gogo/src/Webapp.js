import React, { useState, useEffect, useCallback } from "react";
import Board from "./Board";
import { PLAYER_NAME_LOCALSTORAGE } from "./constants";
import WhoAreUModal from "./WhoAreUModal";
import "./webapp.css";

const Webapp = () => {
  const [socket, setSocket] = useState(null);
  const [nameMissing, setNameMissing] = useState(true);

  useEffect(() => {
    const userName = localStorage.getItem(PLAYER_NAME_LOCALSTORAGE);
    if (!!userName) {
      setNameMissing(false);
    }
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
        {nameMissing && <WhoAreUModal />}
        <Board socket={socket} />
      </header>
    </div>
  );
};

export default Webapp;
