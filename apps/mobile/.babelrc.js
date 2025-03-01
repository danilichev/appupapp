module.exports = function (api) {
  api.cache(true);
  return {
    plugins: [["react-native-unistyles/plugin"]],
    presets: ["babel-preset-expo"],
  };
};
