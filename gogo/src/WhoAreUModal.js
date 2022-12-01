import React, { useState, useCallback, useEffect } from "react";
import { validateNameInput, getColorInput } from "./util";
import "./webapp.css";

const WhoAreUModal = ({ loaded, playerNames, setPlayerName }) => {
  const [showNameInput, setShowNameInput] = useState(false);

  useEffect(() => {
    loaded && setShowNameInput(playerNames.length < 2);
  }, [playerNames]);

  const handleSubmit = useCallback(
    (e) => {
      e.preventDefault();
      const name = validateNameInput(e, playerNames);
      const color = getColorInput(e)
      if (name) {
        console.log("name set:", name);
        setPlayerName(name, color);
      }
    },
    [setPlayerName]
  );
  
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
              margin: "auto",
            }}
          >
            <div style={{ padding: "10px", margin: "auto" }}> who are u?</div>
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
                <div>
                  <input
                    type="text"
                    className="modal-input"
                    name={"newPlayer"}
  
                  />
                </div>
                {playerNames.length < 1 && (
                  <div>
                    <input
                      type="radio"
                      id="black-radio"
                      name="white"
                      value="b"
                    />
                    <label for="black-radio">black</label>
                    <input
                      type="radio"
                      id="white-radio"
                      name="black"
                      value="w"
                    />
                    <label for="white-radio">white</label>
                  </div>
                )}
                <input type="submit" style={{width: '200px', height: '50px', margin: '50px'}}/>
              </form>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default WhoAreUModal;
