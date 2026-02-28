# Audio Transcription Feature

The audio transcription feature allows you to convert audio files to SRT subtitle format with timestamps using OpenAI's Whisper model via a Python script.

## Output Format

The transcription generates **SRT (SubRip Subtitle)** format files with timestamps:

```srt
1
00:00:01,000 --> 00:00:04,500
Hello, welcome to today's program.

2
00:00:04,500 --> 00:00:08,200
In this episode, we'll discuss the latest developments.

3
00:00:08,200 --> 00:00:12,000
Let's get started with our first topic.
```

This format is compatible with:
- Video players (VLC, MPC-HC, etc.)
- Video editors (Premiere, Final Cut, DaVinci Resolve)
- Subtitle editors (Aegisub, Subtitle Edit)
- Streaming platforms (YouTube, Vimeo)
- Web video players

## Prerequisites

### Automatic Setup (Recommended)

Filebrowser automatically manages Python dependencies using a virtual environment:

1. **Python 3.7+** must be installed on the server
2. **requirements.txt** must be present in the filebrowser directory
3. **Internet connection** for first-time setup

On first startup, filebrowser will:
- Create `.filebrowser_venv/` directory
- Install all required packages automatically
- Configure everything for you

**No manual installation needed!**

### Manual Setup (Optional)

If you prefer to manage dependencies yourself, install packages globally:

```bash
pip install -r requirements.txt
```

Or install individually:

```bash
pip install torch transformers accelerate
```

### GPU Support (Optional but Recommended)

For faster transcription, install CUDA-enabled PyTorch:

Edit `requirements.txt` to include:
```
torch>=2.0.0 --index-url https://download.pytorch.org/whl/cu118
transformers>=4.30.0
accelerate>=0.20.0
```

Then restart filebrowser to recreate the virtual environment with GPU support.

### Script Location

The `transcribe.py` script must be accessible from the filebrowser working directory. You can:

1. Place it in the same directory as the filebrowser binary
2. Add it to your system PATH
3. Modify the backend code to use an absolute path

## Backend Implementation

The feature is implemented as a PATCH action on the `/api/resources` endpoint, following the same pattern as extract_audio and copy operations.

### API Usage

**Endpoint:** `PATCH /api/resources/{source_path}?action=transcribe_audio&destination={output_path}`

**Example:**
```
PATCH /api/resources/audio/recording.mp3?action=transcribe_audio&destination=audio/recording.srt
```

**Requirements:**
- User must have `Create` permission
- Python 3.7+ must be installed on the server
- Required Python packages must be installed (torch, transformers)
- Source file must be a valid audio file (mp3, wav, m4a, etc.)

**Output:**
- SRT subtitle file with timestamps
- Compatible with video players and editors
- Each segment includes start/end timestamps

### Supported Audio Formats

The Whisper model supports various audio formats including:
- MP3
- WAV
- M4A
- FLAC
- OGG
- And more

## Frontend Implementation

### User Interface

The transcribe audio button appears in the file listing header when:
- A single file is selected
- The selected file is an audio file (type === "audio")
- The user has create permissions

The button is available in:
- Desktop header bar (with icon `subtitles`)
- Mobile action bar
- Context menu (right-click)

### Behavior

When the user clicks the "Transcribe audio" button:
1. The system generates an output filename by replacing the audio extension with `.srt`
2. The API call is made to transcribe the audio
3. The Python script processes the audio file with timestamps (this may take several minutes)
4. A success toast notification is shown: "Audio transcribed successfully!"
5. The file list is automatically reloaded to show the new SRT subtitle file
6. If an error occurs, an error toast is displayed

### Frontend API

The frontend API provides a convenient function to transcribe audio:

```typescript
import { transcribeAudio } from "@/api/files";

// Transcribe audio file to SRT subtitle format
await transcribeAudio("/audio/recording.mp3", "/audio/recording.srt");
```

### Example Integration

```typescript
// In a Vue component
async function handleTranscribeAudio(audioFile: ResourceItem) {
  const outputPath = audioFile.url.replace(/\.[^.]+$/, '.srt');
  
  try {
    await transcribeAudio(audioFile.url, outputPath);
    // Refresh the file list or show success message
  } catch (error) {
    // Handle error
    console.error('Failed to transcribe audio:', error);
  }
}
```

## Translations

The feature includes translations for:
- English: "Transcribe audio" / "Audio transcribed successfully!"
- Spanish: "Transcribir audio" / "¡Audio transcrito exitosamente!"

To add more languages, update the corresponding JSON files in `frontend/src/i18n/`:
```json
{
  "buttons": {
    "transcribeAudio": "Your translation"
  },
  "success": {
    "audioTranscribed": "Your success message"
  }
}
```

## Technical Details

### Whisper Model

- Uses OpenAI's Whisper Large V3 model
- Automatic language detection (configured for English by default)
- GPU acceleration when available (CUDA)
- Falls back to CPU if GPU is not available

### Performance Considerations

- **First run**: Model download may take several minutes (several GB)
- **GPU**: Transcription is significantly faster with CUDA-enabled GPU
- **CPU**: Transcription will work but may be slow for long audio files
- **Memory**: Requires sufficient RAM/VRAM (4-8GB recommended)

### Model Caching

The Whisper model is cached after the first download in:
- Linux/Mac: `~/.cache/huggingface/`
- Windows: `%USERPROFILE%\.cache\huggingface\`

### Script Behavior

The `transcribe.py` script:
1. Loads the Whisper model (cached after first run)
2. Processes the audio file in 30-second chunks
3. Generates transcription with timestamps
4. Formats output as SRT subtitle file
5. Saves to the specified output file
6. Cleans up memory and GPU cache

### SRT Format Details

Each subtitle entry includes:
- **Sequence number** (1, 2, 3, ...)
- **Timestamp range** (HH:MM:SS,mmm --> HH:MM:SS,mmm)
- **Text content** (transcribed speech)
- **Blank line** (separator)

Example:
```srt
1
00:00:00,000 --> 00:00:03,500
Welcome to the presentation.

2
00:00:03,500 --> 00:00:07,200
Today we'll cover three main topics.
```

### Error Handling

The backend captures script output and returns detailed error messages:
- File not found errors
- Model loading errors
- Transcription failures
- File write errors

## Configuration Options

### Using a Different Model

Edit `transcribe.py` to use a different Whisper model:

```python
# Options: tiny, base, small, medium, large, large-v2, large-v3
model_id = "openai/whisper-base"  # Faster but less accurate
# model_id = "openai/whisper-large-v3"  # Slower but more accurate
```

### Enabling Timestamps

To include timestamps in the transcription, modify the pipeline settings:

```python
pipe = pipeline(
    "automatic-speech-recognition",
    # ... other settings ...
    return_timestamps=True,  # Enable timestamps
)
```

### Language Configuration

To transcribe in a different language, modify the generate_kwargs:

```python
result = pipe(audio_file_path, generate_kwargs={"language": "spanish"})
```

Or remove the language parameter for automatic detection:

```python
result = pipe(audio_file_path)
```

## Troubleshooting

### Script Not Found

If you get "python3: command not found" or "transcribe.py not found":
1. Ensure Python 3 is installed: `python3 --version`
2. Place `transcribe.py` in the filebrowser working directory
3. Make the script executable: `chmod +x transcribe.py`

### Out of Memory Errors

If transcription fails with OOM errors:
1. Use a smaller Whisper model (base or small)
2. Reduce batch_size in the script
3. Close other applications to free memory
4. Use GPU if available

### Slow Transcription

To improve speed:
1. Install CUDA-enabled PyTorch for GPU acceleration
2. Use a smaller Whisper model
3. Ensure sufficient system resources

### Model Download Issues

If model download fails:
1. Check internet connection
2. Verify Hugging Face is accessible
3. Manually download the model and place in cache directory

## Example Workflow

1. **Upload video** → `myvideo.mp4`
2. **Extract audio** → `myvideo.mp3` (using extract audio feature)
3. **Transcribe audio** → `myvideo.txt` (using transcribe audio feature)
4. **View/edit transcription** → Open `myvideo.txt` in the editor

## Security Considerations

- The script executes on the server with filebrowser's permissions
- Ensure the script is not writable by untrusted users
- Consider running filebrowser in a sandboxed environment
- Monitor resource usage to prevent abuse
- Validate file paths to prevent directory traversal

## Performance Benchmarks

Approximate transcription times (1 minute of audio):

| Model | GPU (CUDA) | CPU |
|-------|-----------|-----|
| tiny | ~5 seconds | ~30 seconds |
| base | ~10 seconds | ~1 minute |
| small | ~20 seconds | ~3 minutes |
| medium | ~40 seconds | ~8 minutes |
| large-v3 | ~1 minute | ~20 minutes |

*Times vary based on hardware and audio complexity*
