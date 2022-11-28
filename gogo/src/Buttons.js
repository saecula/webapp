import React from "react";
import "./webapp.css";

export const PassButton = ({ onPass, disabled }) => (
  <div className="button" onClick={onPass} disabled={disabled}>
    Pass
  </div>
);

export const ResignButton = ({ onResign }) => (
  <div className="button" onClick={onResign}>
    Resign
  </div>
);

export default PassButton;
