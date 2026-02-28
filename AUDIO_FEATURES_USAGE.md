# Audio & Video Processing Features

This document covers the audio extraction and transcription features available in File Browser.

## Features Overview

1. **Extract Audio** - Extract audio from video files to MP3 format
2. **Transcribe Audio** - Convert audio files to text using AI transcription

## Extract Audio Feature

### Description
Extract audio tracks from video files and save them as MP3 files using FFmpeg.

### Requirements
- FFmpeg must be installed on the server
- User must have Create permission

### Usage

**From UI:**
1. Select a video file (mp4, avi, mkv, etc.)
2. Click the "Extract audio" button (🎵 music note icon)
3. Audio file is created with the same name but .mp3 extension

**From API:**
```bash
curl -X PATCH "http://localhost:8080/api/resources/files/videos/myvideo.mp4?action=extract_audio&destination=files%2Fvideos%2Fmyvideo.mp3" \
  -H "X-Auth: YOUR_TOKEN"
```

### Technical Details
- Uses FFmpeg with libmp3lame codec
- Audio quality: Level 2 (high quality)
- Output format: MP3
- Video stream is removed (-vn flag)

See [EXTRACT_AUDIO_USAGE.md](EXTRACT_AUDIO_USAGE.md) for detailed documentation.

---

## Transcribe Audio Feature

### Description
Convert audio files to text using OpenAI's Whisper AI model.

### Requirements
- Python 3.7+ must be installed
- Required Python packages: `torch`, `transformers`, `accelerate`
- `transcribe.py` script must be accessible
- User must have Create permission

### Installation

**Automatic Setup (Recommended):**

Just start filebrowser - it will automatically:
1. Create a Python virtual environment (`.filebrowser_venv/`)
2. Install all required dependencies
3. Configure everything for you

```bash
# Ensure requirements.txt is present
./filebrowser
```

**Manual Setup (Optional):**

If you prefer to manage dependencies yourself:

```bash
# Install Python dependencies globally
pip install -r requirements.txt

# Or install individually
pip install torch transformers accelerate

# For GPU support (recommended)
pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118
```

**Note:** The automatic setup creates an isolated virtual environment, avoiding conflicts with system Python packages.

### Usage

**From UI:**
1. Select an audio file (mp3, wav, m4a, etc.)
2. Click the "Transcribe audio" button (📄 subtitles icon)
3. Text file is created with the same name but .txt extension
4. Wait for processing (may take several minutes)

**From API:**
```bash
curl -X PATCH "http://localhost:8080/api/resources/files/audio/recording.mp3?action=transcribe_audio&destination=files%2Faudio%2Frecording.srt" \
  -H "X-Auth: YOUR_TOKEN"
```

### Technical Details
- Uses OpenAI Whisper Large V3 model
- Supports multiple audio formats (mp3, wav, m4a, flac, ogg)
- GPU acceleration when available
- Automatic language detection (configured for English by default)
- Model is cached after first download (~3GB)

See [TRANSCRIBE_AUDIO_USAGE.md](TRANSCRIBE_AUDIO_USAGE.md) for detailed documentation.

---

## Complete Workflow Example

### Video to Transcription Pipeline

1. **Upload video file**
   - Upload: `interview.mp4`

2. **Extract audio**
   - Select `interview.mp4`
   - Click "Extract audio" button
   - Result: `interview.mp3`

3. **Transcribe audio**
   - Select `interview.mp3`
   - Click "Transcribe audio" button
   - Wait for processing
   - Result: `interview.srt` (SRT subtitle format with timestamps)

4. **View/edit transcription**
   - Click on `interview.srt` to view
   - Use with video players or editors
   - Edit if needed using the built-in editor

### Batch Processing

For multiple files, you can use the API in a script:

```bash
#!/bin/bash
# Extract audio from all videos and transcribe

for video in videos/*.mp4; do
  filename=$(basename "$video" .mp4)
  
  # Extract audio
  curl -X PATCH "http://localhost:8080/api/resources/files/videos/$video?action=extract_audio&destination=files%2Faudio%2F$filename.mp3" \
    -H "X-Auth: $TOKEN"
  
  # Transcribe audio
  curl -X PATCH "http://localhost:8080/api/resources/files/audio/$filename.mp3?action=transcribe_audio&destination=files%2Ftranscripts%2F$filename.txt" \
    -H "X-Auth: $TOKEN"
done
```

---

## Frontend Implementation

### Button Visibility

| Feature | Shows When | Icon | Location |
|---------|-----------|------|----------|
| Extract Audio | Single video file selected | 🎵 music_note | Header, Context Menu |
| Transcribe Audio | Single audio file selected | 📄 subtitles | Header, Context Menu |

### API Functions

```typescript
import { extractAudio, transcribeAudio } from "@/api/files";

// Extract audio from video
await extractAudio("/files/videos/myvideo.mp4", "/files/videos/myvideo.mp3");

// Transcribe audio to SRT subtitle format
await transcribeAudio("/files/audio/recording.mp3", "/files/audio/recording.srt");
```

---

## Translations

Both features support internationalization:

### English
- "Extract audio" / "Audio extracted successfully!"
- "Transcribe audio" / "Audio transcribed successfully!"

### Spanish
- "Extraer audio" / "¡Audio extraído exitosamente!"
- "Transcribir audio" / "¡Audio transcrito exitosamente!"

To add more languages, edit `frontend/src/i18n/{language}.json`:

```json
{
  "buttons": {
    "extractAudio": "Your translation",
    "transcribeAudio": "Your translation"
  },
  "success": {
    "audioExtracted": "Your success message",
    "audioTranscribed": "Your success message"
  }
}
```

---

## Performance Considerations

### Extract Audio
- Fast operation (seconds for most videos)
- Depends on video length and server CPU
- No GPU required

### Transcribe Audio
- Slower operation (minutes for long audio)
- First run downloads model (~3GB)
- GPU significantly improves speed
- CPU-only mode is slower but functional

### Resource Usage

| Operation | CPU | Memory | GPU | Disk |
|-----------|-----|--------|-----|------|
| Extract Audio | Low | Low | No | Minimal |
| Transcribe Audio | High | 4-8GB | Optional | 3GB (model) |

---

## Troubleshooting

### Extract Audio Issues

**Error: "ffmpeg: command not found"**
- Install FFmpeg: `sudo apt install ffmpeg` (Ubuntu/Debian)
- Or: `brew install ffmpeg` (macOS)

**Error: "403 Forbidden"**
- Check user has Create permission
- Verify file paths are correct

### Transcribe Audio Issues

**Error: "python3: command not found"**
- Install Python 3.7+
- Ensure python3 is in PATH

**Error: "No module named 'torch'"**
- Install dependencies: `pip install -r requirements.txt`

**Error: "Out of memory"**
- Use smaller Whisper model (edit transcribe.py)
- Close other applications
- Add more RAM or use GPU

**Slow transcription**
- Install CUDA-enabled PyTorch for GPU
- Use smaller Whisper model
- Process shorter audio files

---

## Security Notes

- Both features execute external commands (ffmpeg, python)
- Ensure scripts are not writable by untrusted users
- Consider running in sandboxed environment
- Monitor resource usage to prevent abuse
- Validate file paths to prevent directory traversal

---

## API Reference

### Extract Audio

```
PATCH /api/resources/{source_path}?action=extract_audio&destination={output_path}
```

**Parameters:**
- `source_path`: Path to video file
- `destination`: Path for output MP3 file

**Response:**
- 200: Success
- 403: Forbidden (no permission)
- 500: Internal error (ffmpeg failed)

### Transcribe Audio

```
PATCH /api/resources/{source_path}?action=transcribe_audio&destination={output_path}
```

**Parameters:**
- `source_path`: Path to audio file
- `destination`: Path for output text file

**Response:**
- 200: Success
- 403: Forbidden (no permission)
- 500: Internal error (script failed)

---

## Contributing

To improve these features:

1. **Extract Audio**: Modify `http/resource.go` → `extractAudio()` function
2. **Transcribe Audio**: Modify `transcribe.py` script or `http/resource.go` → `transcribeAudio()` function
3. **Frontend**: Modify `frontend/src/views/files/FileListing.vue`

Submit pull requests with improvements!
