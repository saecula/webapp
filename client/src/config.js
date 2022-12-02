const SERVER_HOST = process.env.REACT_APP_HOST || JSON.parse('throw error');
const SERVER_PORT = process.env.REACT_APP_PORT || JSON.parse('throw error');

export const SERVER_URL = `${SERVER_HOST}:${SERVER_PORT}`;