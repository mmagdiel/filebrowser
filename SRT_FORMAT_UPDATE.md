# SRT Format Update

## Overview

The transcription feature now generates **SRT (SubRip Subtitle)** format files with timestamps instead of plain text files.

## What Changed

### Output Format

**Before:**
- Plain text file (`.txt`)
- No timestamps
- Simple transcription

```
Hello, welcome to today's program. In this episode, we'll discuss the latest developments. Let's get started with our first topic.
```

**After:**
- SRT subtitle file (`.srt`)
- Precise timestamps
- Segmented by speech chunks

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

### File Extension

- **Before:** `audio.txt`
- **After:** `audio.srt`

### Python Script Changes

**File:** `transcribe.py`

**Key Updates:**
1. Added `format_timestamp()` function to convert seconds to SRT format
2. Added `generate_srt()` function to create SRT content
3. Enabled `return_timestamps=True` in Whisper pipeline
4. Added `chunk_length_s=30` for better timestamp accuracy
5. Changed output format from plain text to SRT

### Frontend Changes

**File:** `frontend/src/views/files/FileListing.vue`

**Key Updates:**
- Changed output extension from `.txt` to `.srt`
- Updated comments to reflect subtitle format

## SRT Format Specification

### Structure

Each subtitle entry consists of:

1. **Sequence number** (starting from 1)
2. **Timestamp range** (start --> end)
3. **Text content** (one or more lines)
4. **Blank line** (separator)

### Timestamp Format

`HH:MM:SS,mmm --> HH:MM:SS,mmm`

- `HH`: Hours (00-99)
- `MM`: Minutes (00-59)
- `SS`: Seconds (00-59)
- `mmm`: Milliseconds (000-999)
- Comma (`,`) separates seconds from milliseconds
- Arrow (`-->`) separates start from end time

### Example

```srt
1
00:00:00,000 --> 00:00:03,500
Welcome to the presentation.

2
00:00:03,500 --> 00:00:07,200
Today we'll cover three main topics.

3
00:00:07,200 --> 00:00:12,000
First, let's discuss the background.
```

## Benefits of SRT Format

### ✅ Compatibility

Works with:
- **Video Players**: VLC, MPC-HC, Windows Media Player, QuickTime
- **Video Editors**: Adobe Premiere, Final Cut Pro, DaVinci Resolve
- **Subtitle Editors**: Aegisub, Subtitle Edit, Subtitle Workshop
- **Streaming Platforms**: YouTube, Vimeo, Dailymotion
- **Web Players**: Video.js, Plyr, JW Player

### ✅ Timestamps

- Precise timing for each segment
- Easy synchronization with video
- Better for long-form content
- Useful for video editing

### ✅ Standard Format

- Industry-standard subtitle format
- Widely supported
- Simple text-based format
- Easy to edit manually

### ✅ Use Cases

- **Video Subtitles**: Add subtitles to videos
- **Accessibility**: Provide captions for hearing-impaired
- **Translation**: Base for subtitle translation
- **Searchability**: Find specific moments in audio/video
- **Editing**: Identify sections to cut or edit

## Technical Implementation

### Whisper Configuration

```python
pipe = pipeline(
    "automatic-speech-recognition",
    model=model,
    tokenizer=processor.tokenizer,
    feature_extractor=processor.feature_extractor,
    max_new_tokens=128,
    torch_dtype=torch_dtype,
    device=device,
    batch_size=1,
    return_timestamps=True,  # Enable timestamps
)
```

### Transcription with Chunking

```python
result = pipe(
    audio_file_path,
    generate_kwargs={"language": "english"},
    return_timestamps=True,
    chunk_length_s=30,  # Process in 30-second chunks
)
```

### SRT Generation

```python
def generate_srt(segments):
    """Generate SRT format from transcription segments"""
    srt_content = []
    
    for i, segment in enumerate(segments, start=1):
        start_time = format_timestamp(segment['timestamp'][0])
        end_time = format_timestamp(segment['timestamp'][1])
        text = segment['text'].strip()
        
        srt_content.append(f"{i}")
        srt_content.append(f"{start_time} --> {end_time}")
        srt_content.append(text)
        srt_content.append("")  # Empty line
    
    return "\n".join(srt_content)
```

## Usage Examples

### Basic Transcription

```bash
# From command line
python3 transcribe.py audio.mp3 subtitles.srt

# Output: subtitles.srt with timestamps
```

### From File Browser UI

1. Select audio file: `podcast.mp3`
2. Click "Transcribe audio" button
3. Wait for processing
4. Result: `podcast.srt` created

### With Video

1. Upload video: `lecture.mp4`
2. Extract audio: `lecture.mp3`
3. Transcribe: `lecture.srt`
4. Use SRT file with video in any player

### API Call

```bash
curl -X PATCH \
  "http://localhost:8080/api/resources/files/audio/speech.mp3?action=transcribe_audio&destination=files%2Faudio%2Fspeech.srt" \
  -H "X-Auth: YOUR_TOKEN"
```

## Editing SRT Files

### Manual Editing

SRT files are plain text and can be edited with any text editor:

```srt
1
00:00:00,000 --> 00:00:03,500
[Edit this text as needed]

2
00:00:03,500 --> 00:00:07,200
[Adjust timestamps if necessary]
```

### Subtitle Editors

For advanced editing, use dedicated subtitle editors:

- **Aegisub**: Full-featured subtitle editor
- **Subtitle Edit**: Windows subtitle editor
- **Subtitle Workshop**: Classic subtitle tool
- **Jubler**: Cross-platform subtitle editor

### Common Edits

- Fix transcription errors
- Adjust timing
- Split long subtitles
- Merge short subtitles
- Add formatting (italics, bold)
- Synchronize with video

## Performance Impact

### Processing Time

- **Slightly slower** than plain text mode
- Timestamp generation adds ~10-20% overhead
- Chunking improves accuracy but takes more time

### Memory Usage

- Similar to plain text mode
- Timestamps stored in memory during processing
- Cleaned up after generation

### File Size

- SRT files are slightly larger than plain text
- Timestamps add ~30-40% to file size
- Still very small (few KB for most audio)

## Backward Compatibility

### Migration

If you have existing `.txt` transcriptions:

1. They will continue to work
2. New transcriptions create `.srt` files
3. No automatic conversion of old files
4. Both formats can coexist

### Manual Conversion

To convert existing plain text to SRT:

```python
# Simple conversion (no timestamps)
text = open('old.txt').read()
srt = "1\n00:00:00,000 --> 00:00:10,000\n" + text
open('new.srt', 'w').write(srt)
```

For proper timestamps, re-transcribe the audio.

## Troubleshooting

### No Timestamps in Output

**Symptom:** SRT file has generic timestamps

**Cause:** Whisper couldn't generate timestamps

**Solution:**
- Check audio quality
- Try shorter audio files
- Ensure `return_timestamps=True` in script

### Incorrect Timestamps

**Symptom:** Timestamps don't match audio

**Cause:** Audio processing issues

**Solutions:**
- Re-transcribe the audio
- Use shorter chunk_length_s
- Check audio file integrity

### SRT Format Errors

**Symptom:** Video player doesn't recognize SRT

**Cause:** Format issues

**Solutions:**
- Verify file encoding (UTF-8)
- Check for missing blank lines
- Validate with SRT validator tool

## Configuration Options

### Chunk Length

Adjust processing chunk size in `transcribe.py`:

```python
result = pipe(
    audio_file_path,
    chunk_length_s=30,  # Change this (10-60 seconds)
)
```

- **Smaller chunks** (10-20s): More accurate timestamps, slower
- **Larger chunks** (40-60s): Faster processing, less accurate

### Language

Change transcription language:

```python
result = pipe(
    audio_file_path,
    generate_kwargs={"language": "spanish"},  # or "french", "german", etc.
)
```

### Disable Timestamps

To revert to plain text (not recommended):

```python
pipe = pipeline(
    # ... other settings ...
    return_timestamps=False,  # Disable timestamps
)
```

## Future Enhancements

Potential improvements:

- [ ] VTT format support (WebVTT)
- [ ] ASS/SSA format support (Advanced SubStation Alpha)
- [ ] Automatic subtitle styling
- [ ] Speaker diarization (identify different speakers)
- [ ] Subtitle translation
- [ ] Subtitle synchronization tools
- [ ] Batch processing multiple files

## Comparison: TXT vs SRT

| Feature | TXT (Old) | SRT (New) |
|---------|-----------|-----------|
| Timestamps | ❌ No | ✅ Yes |
| Video Player Support | ❌ Limited | ✅ Full |
| Editing Tools | Basic text editors | Subtitle editors |
| File Size | Smaller | Slightly larger |
| Use Cases | Simple transcription | Video subtitles |
| Searchability | Text only | Text + timing |
| Accessibility | Basic | Full captions |
| Industry Standard | No | Yes |

## Conclusion

The SRT format provides:

- ✅ **Professional** subtitle format
- ✅ **Universal** compatibility
- ✅ **Precise** timestamps
- ✅ **Better** usability
- ✅ **Industry** standard

This makes the transcription feature much more useful for video production, accessibility, and professional workflows.

---

**Update Date**: February 2026
**Format**: SRT (SubRip Subtitle)
**Extension**: `.srt`
**Compatibility**: Universal
