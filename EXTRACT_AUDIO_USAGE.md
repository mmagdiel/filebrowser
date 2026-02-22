# Extract Audio Feature

The extract audio feature allows you to extract audio from video files using FFmpeg.

## Backend Implementation

The feature is implemented as a PATCH action on the `/api/resources` endpoint, following the same pattern as copy and rename operations.

### API Usage

**Endpoint:** `PATCH /api/resources/{source_path}?action=extract_audio&destination={output_path}`

**Example:**
```
PATCH /api/resources/videos/myvideo.mp4?action=extract_audio&destination=audio/myvideo.mp3
```

**Requirements:**
- User must have `Create` permission
- FFmpeg must be installed on the server
- Source file must be a valid video file with audio

## Frontend Usage

The frontend API provides a convenient function to extract audio:

```typescript
import { extractAudio } from "@/api/files";

// Extract audio from a video file
await extractAudio("/videos/myvideo.mp4", "/audio/myvideo.mp3");
```

### Example Integration

```typescript
// In a Vue component
async function handleExtractAudio(videoFile: ResourceItem) {
  const outputPath = videoFile.path.replace(/\.[^.]+$/, '.mp3');
  
  try {
    await extractAudio(videoFile.path, outputPath);
    // Refresh the file list or show success message
  } catch (error) {
    // Handle error
    console.error('Failed to extract audio:', error);
  }
}
```

## Technical Details

- The implementation uses FFmpeg with the following parameters:
  - `-vn`: No video (audio only)
  - `-acodec libmp3lame`: Use MP3 codec
  - `-q:a 2`: Audio quality level 2 (high quality)
- The feature properly handles virtual file systems (afero.Fs)
- Real paths are resolved for BasePathFs implementations
