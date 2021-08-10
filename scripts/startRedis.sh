# Start redis server via homebrew
# Port defaults to localhost:6379
brew services start redis

# To test connection use (should receive PONG):
redis-cli ping

# To stop use
brew services stop redis