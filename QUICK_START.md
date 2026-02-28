# Quick Start Guide - Audio Features

Get started with audio extraction and transcription in 5 minutes!

## 1. Extract Audio (Already Working!)

### Setup
✅ No setup needed if FFmpeg is installed

### Usage
1. Select a video file
2. Click "Extract audio" button (🎵)
3. Done! MP3 file created

---

## 2. Transcribe Audio (New Feature!)

### Quick Setup

```bash
# Just start filebrowser - it handles everything automatically!
./filebrowser
```

**That's it!** On first startup:
- ✅ Creates Python virtual environment automatically
- ✅ Installs all dependencies (takes ~5-10 minutes first time)
- ✅ Ready to use - no manual installation needed!

### First Use

1. Select an audio file (mp3, wav, etc.)
2. Click "Transcribe audio" button (📄)
3. Wait (first time downloads AI model ~3GB)
4. Done! Text file created

**Note:** First transcription takes longer due to model download. Subsequent transcriptions are much faster.

---

## Complete Example

### Video → Audio → Text

```
1. Upload: interview.mp4
   ↓
2. Click "Extract audio" → interview.mp3
   ↓
3. Click "Transcribe audio" → interview.srt (SRT subtitle format)
   ↓
4. View/edit: interview.srt (compatible with video players)
```

**Time**: 
- Extract audio: ~10 seconds
- Transcribe (first time): ~10 minutes (model download)
- Transcribe (after): ~2-5 minutes (depending on audio length)

---

## Troubleshooting

### Extract Audio Not Working?

```bash
# Install FFmpeg
sudo apt install ffmpeg  # Ubuntu/Debian
brew install ffmpeg      # macOS
```

### Transcribe Audio Not Working?

```bash
# Ensure Python 3 is installed
python3 --version

# Ensure requirements.txt exists
ls requirements.txt

# Delete venv and let filebrowser recreate it
rm -rf .filebrowser_venv
./filebrowser
```

### Still Having Issues?

See detailed guides:
- [SETUP_TRANSCRIPTION.md](SETUP_TRANSCRIPTION.md) - Full setup guide
- [AUDIO_FEATURES_USAGE.md](AUDIO_FEATURES_USAGE.md) - Complete documentation
- [TRANSCRIBE_AUDIO_USAGE.md](TRANSCRIBE_AUDIO_USAGE.md) - Transcription details

---

## Performance Tips

### Faster Transcription

1. **Use GPU** (10x faster)
   ```bash
   pip install torch --index-url https://download.pytorch.org/whl/cu118
   ```

2. **Use smaller model** (edit transcribe.py)
   ```python
   model_id = "openai/whisper-base"  # Instead of large-v3
   ```

3. **Process shorter files**
   - Split long audio into chunks
   - Or use extract audio to get specific segments

---

## Quick Reference

### Buttons

| Button | Icon | Shows For | Creates |
|--------|------|-----------|---------|
| Extract audio | 🎵 | Video files | .mp3 |
| Transcribe audio | 📄 | Audio files | .srt |

### Locations

- Header bar (desktop)
- Action bar (mobile)
- Context menu (right-click)

### Permissions

Both features require:
- ✅ Create permission
- ✅ Access to source file

---

## API Usage

### Extract Audio
```bash
curl -X PATCH \
  "http://localhost:8080/api/resources/files/video.mp4?action=extract_audio&destination=files%2Fvideo.mp3" \
  -H "X-Auth: TOKEN"
```

### Transcribe Audio
```bash
curl -X PATCH \
  "http://localhost:8080/api/resources/files/audio.mp3?action=transcribe_audio&destination=files%2Faudio.srt" \
  -H "X-Auth: TOKEN"
```

---

## What's Next?

1. ✅ Try extracting audio from a video
2. ✅ Try transcribing an audio file
3. 📖 Read full documentation for advanced features
4. 🚀 Integrate into your workflow

---

## Support

- 📖 [Full Documentation](AUDIO_FEATURES_USAGE.md)
- 🔧 [Setup Guide](SETUP_TRANSCRIPTION.md)
- 🐛 [Troubleshooting](SETUP_TRANSCRIPTION.md#troubleshooting)
- 💡 [Implementation Details](IMPLEMENTATION_SUMMARY.md)

---

**Ready to go!** Start by selecting a video or audio file and look for the new buttons in the header bar.
