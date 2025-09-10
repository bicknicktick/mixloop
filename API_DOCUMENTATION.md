# MixLoop Audio API Documentation

## Overview
MixLoop Audio API adalah backend service untuk memproses audio files menjadi sequence seamless dengan crossfade transitions dan audio enhancement.

## Base URL
```
http://localhost:8081
```

## Endpoints

### POST /mix
Memproses multiple audio files menjadi satu sequence dengan crossfade dan enhancement.

#### Request
- **Method**: POST
- **Content-Type**: multipart/form-data
- **Parameters**:
  - `audio` (files): Multiple audio files (MP3/WAV)
  - `loops` (int, optional): Jumlah loop (default: 1)
  - `crossfade` (float, optional): Durasi crossfade dalam detik (default: 2.0)
  - `enhance` (bool, optional): Enable audio enhancement (default: true)
  - `format` (string, optional): Output format "mp3" atau "wav" (default: "mp3")

#### Response
- **Content-Type**: audio/mpeg atau audio/wav
- **Body**: Binary audio file

#### Example cURL
```bash
curl -X POST http://localhost:8081/mix \
  -F "audio=@tone1.mp3" \
  -F "audio=@tone2.mp3" \
  -F "audio=@tone3.mp3" \
  -F "loops=2" \
  -F "crossfade=1.5" \
  -F "enhance=true" \
  -F "format=mp3" \
  -o output.mp3
```

### GET /health
Health check endpoint.

#### Response
```json
{
  "status": "ok",
  "message": "MixLoop Audio API is running"
}
```

## Audio Processing Features

### 1. Sequential Processing
- Audio files dirangkai berurutan: audio1 → audio2 → audio3
- BUKAN overlay mixing (audio1 + audio2 + audio3)

### 2. Crossfade Transitions
- Smooth transitions antar tracks
- Durasi crossfade dapat dikustomisasi
- Crossfade juga diterapkan pada loop boundaries

### 3. Audio Enhancement (Default: ON)
Filter chain yang diterapkan:
- `highpass=f=80` - Hilangkan low-frequency hum
- `lowpass=f=16000` - Hilangkan ultrasonic noise  
- `acompressor=threshold=-20dB:ratio=3` - Stabilkan dinamika
- `loudnorm=I=-14:TP=-2:LRA=11` - Normalisasi loudness

### 4. Export Quality
- **MP3**: 320kbps, 48kHz
- **WAV**: 24-bit PCM, 48kHz

## Error Responses

### 400 Bad Request
```json
{
  "error": "No audio files provided"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to process audio: [error details]"
}
```

## Dependencies
- FFmpeg (untuk audio processing)
- Go 1.19+
- Gorilla Mux (routing)
- rs/cors (CORS handling)

## Installation & Setup

1. Install FFmpeg:
```bash
sudo apt install ffmpeg
```

2. Run server:
```bash
cd backend
go mod tidy
go run main.go
```

Server akan berjalan di port 8081.
