module.exports = {
  verbose: true,
  "transform": {
    "^.+\\.(ts|tsx)$": "ts-jest"
  },
  "testRegex": "(\\.|/)([jt]est)\\.[jt]s$",
  "moduleFileExtensions": [
    "ts",
    "tsx",
    "js",
    "jsx",
    "json"
  ],
  "coverageDirectory": "./tmp/coverage/",
};
