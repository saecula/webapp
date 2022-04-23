import React, { useState, useCallback, useEffect } from "react";
import { validateNameInput } from "./util";
import "./webapp.css";

const WhoAreUModal = ({ loaded, playerNames, setPlayerName }) => {
  const [showNameInput, setShowNameInput] = useState(false);

  useEffect(() => {
    loaded && setShowNameInput(playerNames.length < 2);
  }, [playerNames]);

  console.log("shownamein", showNameInput);
  const handleSubmit = useCallback((e) => {
    console.log("here", e);
    e.preventDefault();
    const input = validateNameInput(e, playerNames);
    if (input) {
      console.log("name set:", input);
      setPlayerName(input);
    }
  }, []);
  console.log("pnames", playerNames.length);
  return (
    <div className="modal-container">
      <div className="modal">
        {loaded && (
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              padding: "10px",
              textAlign: "left",
              margin: "auto",
            }}
          >
            <div style={{ padding: "10px" }}> who are u?</div>
            {playerNames.map((pn) => (
              <button
                id={pn}
                value={pn}
                style={{
                  padding: "10px",
                  margin: "10px",
                  fontSize: ".8em",
                  backgroundColor: "white",
                }}
                onClick={handleSubmit}
              >
                {pn}
              </button>
            ))}
            {showNameInput && (
              <form onSubmit={handleSubmit}>
                <input type="text" className="modal-input" name={"newPlayer"} />
                <input type="submit" style={{ display: "none" }} />
              </form>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default WhoAreUModal;
