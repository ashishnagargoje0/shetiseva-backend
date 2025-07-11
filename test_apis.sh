#!/bin/bash

TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFzaGlzaEBleGFtcGxlLmNvbSIsInJvbGUiOiJ1c2VyIiwiZXhwIjoxNzUyMzA2MDcyLCJpYXQiOjE3NTIyMTk2NzJ9.VZ3CbM4ZW9eLwhobBAgwaFIBiO7Rvtl2XQc-Zjj785E"

pretty_print() {
  if command -v jq >/dev/null; then
    jq .
  else
    cat
  fi
}

echo "🔍 Testing /weather/today"
curl -s -X GET http://localhost:8080/weather/today -H "Authorization: Bearer $TOKEN" | pretty_print
echo -e "\n-------------------------------------\n"

echo "🔍 Testing /weather/forecast"
curl -s -X GET "http://localhost:8080/weather/forecast?location=pune" -H "Authorization: Bearer $TOKEN" | pretty_print
echo -e "\n-------------------------------------\n"

echo "🔍 Testing /advisory/crop/123"
curl -s -X GET http://localhost:8080/advisory/crop/123 -H "Authorization: Bearer $TOKEN" | pretty_print
echo -e "\n-------------------------------------\n"

echo "🧠 Testing /ai/diagnose"
curl -s -X POST http://localhost:8080/ai/diagnose \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"image": "base64string"}' | pretty_print
echo -e "\n-------------------------------------\n"

echo "🤖 Testing /chatbot/message"
curl -s -X POST http://localhost:8080/chatbot/message \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello"}' | pretty_print
echo -e "\n-------------------------------------\n"

echo "📜 Testing /chatbot/history"
curl -s -X GET http://localhost:8080/chatbot/history -H "Authorization: Bearer $TOKEN" | pretty_print
echo -e "\n-------------------------------------\n"

echo "⚠️  Testing /ai/alerts"
curl -s -X GET http://localhost:8080/ai/alerts -H "Authorization: Bearer $TOKEN" | pretty_print
echo -e "\n✅ All API tests completed."

