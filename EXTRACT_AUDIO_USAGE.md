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

## Frontend Implementation

### User Interface

The extract audio button appears in the file listing header when:
- A single file is selected
- The selected file is a video (type === "video")
- The user has create permissions

The button is available in:
- Desktop header bar (with icon `music_note`)
- Mobile action bar
- Context menu (right-click)

### Behavior

When the user clicks the "Extract audio" button:
1. The system generates an output filename by replacing the video extension with `.mp3`
2. The API call is made to extract the audio
3. A success toast notification is shown: "Audio extracted successfully!"
4. The file list is automatically reloaded to show the new audio file
5. If an error occurs, an error toast is displayed

### Frontend API

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

## Translations

The feature includes translations for:
- English: "Extract audio" / "Audio extracted successfully!"
- Spanish: "Extraer audio" / "¡Audio extraído exitosamente!"

To add more languages, update the corresponding JSON files in `frontend/src/i18n/`:
```json
{
  "buttons": {
    "extractAudio": "Your translation"
  },
  "success": {
    "audioExtracted": "Your success message"
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
- The output file is always created in the same directory as the source video
- The output filename matches the video filename with `.mp3` extension
