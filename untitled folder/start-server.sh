#!/bin/bash
# Simple script to start a local web server for the overlay
# This ensures localStorage works properly across tabs

cd "$(dirname "$0")"

# Check if portable server exists (preferred)
if [ -f "./server" ] && [ -x "./server" ]; then
    clear
    ./server
    exit 0
fi

# Check for Python 3
if command -v python3 &> /dev/null; then
    clear
    echo "=========================================="
    echo "   Pokexclusive Overlay Server"
    echo "=========================================="
    echo ""
    echo "Starting server..."
    echo ""
    echo "✓ Server is running!"
    echo ""
    echo "IMPORTANT: Keep this window open!"
    echo ""
    echo "Next: Open this URL in your browser:"
    echo "  http://localhost:8000/control.html"
    echo ""
    echo "To stop the server:"
    echo "  - Close this window, OR"
    echo "  - Press Ctrl+C"
    echo ""
    echo "=========================================="
    echo ""
    python3 -m http.server 8000
    exit 0
fi

# No server found
clear
echo "=========================================="
echo "    ERROR: No Server Found"
echo "=========================================="
echo ""
echo "Neither ./server nor python3 was found."
echo ""
echo "Please either:"
echo "  1. Build ./server (see BUILD-SERVER.md)"
echo "  2. Install Python 3"
echo ""
echo "=========================================="
echo ""
exit 1
