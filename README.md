# MixLoop - Audio Sequence Mixer

Aplikasi web minimalist untuk menggabungkan multiple audio files menjadi sequence seamless dengan crossfade transitions dan audio enhancement.

**Developed by [BITZY.ID](https://bitzy.id)**

## Features

- 🎵 **Sequential Audio Processing** - Menggabungkan audio berurutan (audio1→audio2→audio3), bukan overlay
- 🔄 **Seamless Looping** - Loop dengan crossfade di boundaries untuk hasil seamless
- ✨ **Audio Enhancement** - Filter untuk meningkatkan kualitas audio
- 📱 **Mobile-First UI** - Responsive design terinspirasi Threads by Instagram
- 🎛️ **Drag & Drop Upload** - Interface modern untuk upload files
- 📤 **Multiple Export Formats** - MP3 320k dan WAV 24-bit

## Tech Stack

### Backend (Golang)
- Gorilla Mux router
- FFmpeg untuk audio processing
- CORS support
- Modular architecture

### Frontend (React)
- React 18 + Vite
- TailwindCSS untuk styling
- Axios untuk API calls
- Modern glass morphism UI

## Quick Start

### Prerequisites
- Go 1.19+
- Node.js 16+
- FFmpeg

### Installation

1. **Install FFmpeg**:
```bash
sudo apt install ffmpeg
```

2. **Backend Setup**:
```bash
cd backend
go mod tidy
go run main.go
```

3. **Frontend Setup**:
```bash
cd frontend
npm install
npm start
```

4. **Access Application**:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8081

## Audio Processing Pipeline

1. **Validation** - Validasi format dan integrity audio files
2. **Sequence Creation** - Concatenate files dengan crossfade transitions
3. **Loop Application** - Apply looping dengan crossfade di boundaries
4. **Enhancement** - Apply audio filters untuk kualitas optimal
5. **Export** - Output dalam format MP3 320k atau WAV 24-bit

## Enhancement Filters

- `loudnorm=I=-14:TP=-2:LRA=11` - Loudness normalization
- `highpass=f=80` - Remove low-frequency hum
- `lowpass=f=16000` - Remove ultrasonic noise
- `acompressor=threshold=-20dB:ratio=3` - Dynamic compression

## API Usage

```bash
curl -X POST http://localhost:8081/mix \
  -F "audio=@file1.mp3" \
  -F "audio=@file2.mp3" \
  -F "loops=2" \
  -F "crossfade=1.5" \
  -F "enhance=true" \
  -F "format=mp3" \
  -o output.mp3
```

## Project Structure

```
mixloop/
├── backend/           # Golang API server
│   ├── handlers/      # HTTP handlers
│   ├── utils/         # Audio processing utilities
│   ├── main.go        # Server entry point
│   └── go.mod         # Go dependencies
├── frontend/          # React application
│   ├── src/           # Source code
│   ├── public/        # Static assets
│   └── package.json   # NPM dependencies
├── test_audio/        # Sample audio files
└── test_output/       # Test results
```

## License

MIT License
