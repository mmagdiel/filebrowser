#!/usr/bin/env python3
"""
Audio Transcription Script using Whisper
Transcribes audio files to SRT subtitle format with timestamps using OpenAI's Whisper model.
"""

import torch
from transformers import AutoModelForSpeechSeq2Seq, AutoProcessor, pipeline
import sys
import os
import gc


def format_timestamp(seconds):
    """Convert seconds to SRT timestamp format (HH:MM:SS,mmm)"""
    hours = int(seconds // 3600)
    minutes = int((seconds % 3600) // 60)
    secs = int(seconds % 60)
    millis = int((seconds % 1) * 1000)
    return f"{hours:02d}:{minutes:02d}:{secs:02d},{millis:03d}"


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
        srt_content.append("")  # Empty line between subtitles
    
    return "\n".join(srt_content)


def setup_whisper_pipeline():
    """Initialize the Whisper model and processor."""
    device = "cpu" # "cuda:0" if torch.cuda.is_available() else "cpu"
    # torch_dtype = torch.float16 if torch.cuda.is_available() else torch.float32
    torch_dtype = torch.float32
    model_id = "openai/whisper-small"
    
    print("Loading model (this may take a few minutes)...", file=sys.stderr)
    
    # Load model with memory optimization
    model = AutoModelForSpeechSeq2Seq.from_pretrained(
        model_id,
        torch_dtype=torch_dtype,
        low_cpu_mem_usage=True,
        use_safetensors=True
    )
    model.to(device)
    
    # Load processor
    processor = AutoProcessor.from_pretrained(model_id)
    
    # Create pipeline with timestamps enabled
    pipe = pipeline(
        "automatic-speech-recognition",
        model=model,
        tokenizer=processor.tokenizer,
        feature_extractor=processor.feature_extractor,
        max_new_tokens=128,
        torch_dtype=torch_dtype,
        device=device,
        batch_size=1,  # Reduced for memory efficiency
        return_timestamps=True,  # Enable timestamps for SRT format
    )
    
    return pipe


def transcribe_audio(audio_file_path, output_file=None):
    """Transcribe audio file to SRT subtitle format."""
    if not os.path.exists(audio_file_path):
        raise FileNotFoundError(f"Audio file not found: {audio_file_path}")
    
    print("Setting up Whisper pipeline...", file=sys.stderr)
    pipe = setup_whisper_pipeline()
    
    print(f"Transcribing: {audio_file_path}", file=sys.stderr)
    print("This may take several minutes depending on audio length...", file=sys.stderr)
    
    # Transcribe with timestamps and chunking
    result = pipe(
        audio_file_path,
        generate_kwargs={"language": "english"},
        return_timestamps=True,
        chunk_length_s=30,  # Process in 30-second chunks
    )
    
    # Generate SRT format
    if 'chunks' in result:
        srt_content = generate_srt(result['chunks'])
    else:
        # Fallback: create single subtitle if no chunks
        srt_content = "1\n00:00:00,000 --> 00:00:10,000\n" + result['text'].strip()
    
    # Clean up memory
    del pipe
    gc.collect()
    if torch.cuda.is_available():
        torch.cuda.empty_cache()
    
    # Save to file if specified
    if output_file:
        try:
            with open(output_file, 'w', encoding='utf-8') as f:
                f.write(srt_content)
            print(f"Transcription saved to: {output_file}", file=sys.stderr)
        except Exception as e:
            print(f"Error saving to file: {e}", file=sys.stderr)
            raise
    
    return srt_content


def main():
    if len(sys.argv) < 2:
        print("Usage: python transcribe.py <audio_file> [output_file]", file=sys.stderr)
        print("Example: python transcribe.py audio.mp3 subtitles.srt", file=sys.stderr)
        sys.exit(1)
    
    audio_file = sys.argv[1]
    output_file = sys.argv[2] if len(sys.argv) > 2 else None
    
    try:
        transcription = transcribe_audio(audio_file, output_file)
        
        # Print transcription to stdout if no output file specified
        if not output_file:
            print("\nTranscription (SRT format):")
            print("-" * 50)
            print(transcription)
    except KeyboardInterrupt:
        print("\nTranscription interrupted by user.", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
