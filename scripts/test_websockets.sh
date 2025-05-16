PORT=8000
if [ "$(uname -s)" = "Linux" ]; then
    PYTHON_CMD=python3
    xdg-open "http://localhost:$PORT/client.html" &
else
    PYTHON_CMD=python
    start "http://localhost:$PORT/client.html"
fi

cd ../Tests
$PYTHON_CMD -m http.server $PORT &
cd -
echo "Python server running on port $PORT. Press Ctrl+C to stop."
wait