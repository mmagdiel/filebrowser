# Audio Transcription Setup Guide

Quick setup guide for the audio transcription feature.

## Quick Start

### 1. Place Requirements File

Ensure `requirements.txt` is in the same directory as the filebrowser binary:

```bash
# The requirements.txt should contain:
# torch>=2.0.0
# transformers>=4.30.0
# accelerate>=0.20.0
```

### 2. Start File Browser

```bash
./filebrowser
```

**That's it!** On first startup, filebrowser will:
- ✅ Automatically create a Python virtual environment (`.filebrowser_venv/`)
- ✅ Install all required dependencies
- ✅ Configure everything for you

The transcription feature is now ready to use!

---

## What Happens on First Startup

When you start filebrowser for the first time:

1. **Virtual Environment Creation** (~10 seconds)
   - Creates `.filebrowser_venv/` directory
   - Sets up isolated Python environment

2. **Dependency Installation** (~5-10 minutes)
   - Downloads and installs PyTorch
   - Installs Transformers library
   - Installs Accelerate library

3. **Ready to Use**
   - No global Python packages installed
   - Everything isolated in `.filebrowser_venv/`
   - Subsequent startups are instant

### Console Output

You'll see messages like:
```
Checking Python virtual environment...
Creating Python virtual environment...
Installing Python dependencies (this may take a few minutes)...
Upgrading pip...
Installing requirements from requirements.txt...
Dependencies installed successfully
Python virtual environment ready
```

---

## Detailed Setup

### Prerequisites

- Python 3.7 or higher
- pip (Python package manager)
- 4-8GB RAM (or GPU with 4GB+ VRAM)
- Internet connection (for first-time setup)

### Installation Steps

**No manual installation needed!** Just ensure:

1. Python 3 is installed:
   ```bash
   python3 --version
   ```

2. `requirements.txt` exists in the filebrowser directory

3. Start filebrowser - it handles the rest automatically

### Virtual Environment Location

The virtual environment is created at:
```
.filebrowser_venv/
├── bin/           # Python interpreter and tools (Linux/Mac)
├── Scripts/       # Python interpreter and tools (Windows)
├── lib/           # Installed packages
└── pyvenv.cfg     # Configuration
```

**Note:** This directory is automatically managed. Don't modify it manually.

---

## Configuration

### Script Location

The `transcribe.py` script must be in the same directory as filebrowser.

### Model Selection

Edit `transcribe.py` to change the Whisper model:

```python
# Faster but less accurate
model_id = "openai/whisper-tiny"    # ~39M parameters
model_id = "openai/whisper-base"    # ~74M parameters
model_id = "openai/whisper-small"   # ~244M parameters

# Balanced
model_id = "openai/whisper-medium"  # ~769M parameters

# Most accurate (default)
model_id = "openai/whisper-large-v3"  # ~1550M parameters
```

### Language Configuration

Edit `transcribe.py` to change the transcription language:

```python
# For Spanish
result = pipe(audio_file_path, generate_kwargs={"language": "spanish"})

# For French
result = pipe(audio_file_path, generate_kwargs={"language": "french"})

# For automatic detection (slower)
result = pipe(audio_file_path)
```

---

## Troubleshooting

### Common Issues

#### "Python not found in PATH"

```bash
# Install Python 3
sudo apt install python3 python3-venv  # Ubuntu/Debian
brew install python3                    # macOS
```

#### "Virtual environment creation failed"

```bash
# Ensure python3-venv is installed
sudo apt install python3-venv  # Ubuntu/Debian

# Or use python instead of python3
# (filebrowser will try both)
```

#### "Pip install failed"

Solutions:
1. Check internet connection
2. Ensure sufficient disk space (~5GB)
3. Try deleting `.filebrowser_venv/` and restart filebrowser

```bash
rm -rf .filebrowser_venv
./filebrowser
```

#### "requirements.txt not found"

```bash
# Create requirements.txt in the same directory as filebrowser
cat > requirements.txt << EOF
torch>=2.0.0
transformers>=4.30.0
accelerate>=0.20.0
EOF
```

#### "Warning: Failed to initialize Python virtual environment"

If you see this warning:
- Audio transcription will not be available
- Other features work normally
- Check the error message for details
- Verify Python 3 is installed

#### Virtual Environment Corrupted

```bash
# Delete and recreate
rm -rf .filebrowser_venv
./filebrowser
```

---

## Manual Testing

Test the virtual environment manually:

```bash
# Activate the venv
source .filebrowser_venv/bin/activate  # Linux/Mac
.filebrowser_venv\Scripts\activate     # Windows

# Test imports
python -c "import torch; import transformers; print('OK')"

# Test transcription script
python transcribe.py sample.mp3 output.txt

# Deactivate
deactivate
```

---

## Performance Optimization

### GPU Acceleration

For best performance, use a CUDA-enabled GPU:

```bash
# Edit requirements.txt to use CUDA-enabled PyTorch
# Replace the torch line with:
torch>=2.0.0 --index-url https://download.pytorch.org/whl/cu118
```

Then delete the venv and restart:
```bash
rm -rf .filebrowser_venv
./filebrowser
```

### Model Size vs Speed

| Model | Size | Speed (GPU) | Speed (CPU) | Accuracy |
|-------|------|-------------|-------------|----------|
| tiny | 39M | Very Fast | Fast | Good |
| base | 74M | Fast | Moderate | Better |
| small | 244M | Moderate | Slow | Good |
| medium | 769M | Slow | Very Slow | Better |
| large-v3 | 1550M | Very Slow | Extremely Slow | Best |

### Memory Requirements

| Model | GPU VRAM | System RAM |
|-------|----------|------------|
| tiny | 1GB | 2GB |
| base | 1GB | 2GB |
| small | 2GB | 4GB |
| medium | 5GB | 8GB |
| large-v3 | 10GB | 16GB |

---

## Docker Setup

If running filebrowser in Docker:

### Dockerfile

```dockerfile
FROM filebrowser/filebrowser:latest

# Install Python and venv
RUN apk add --no-cache python3 py3-pip python3-dev

# Copy transcription files
COPY transcribe.py /app/transcribe.py
COPY requirements.txt /app/requirements.txt

# Set working directory
WORKDIR /app

# Virtual environment will be created automatically on first run
CMD ["filebrowser"]
```

### Docker Compose

```yaml
version: '3.8'
services:
  filebrowser:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/srv
      - ./database:/database
      - ./transcribe.py:/app/transcribe.py
      - ./requirements.txt:/app/requirements.txt
      - ./.filebrowser_venv:/app/.filebrowser_venv  # Persist venv
    environment:
      - FB_DATABASE=/database/filebrowser.db
```

**Note:** Persisting `.filebrowser_venv` as a volume avoids reinstalling dependencies on container restart.

---

## Production Deployment

### Systemd Service

Create `/etc/systemd/system/filebrowser.service`:

```ini
[Unit]
Description=File Browser
After=network.target

[Service]
Type=simple
User=filebrowser
WorkingDirectory=/opt/filebrowser
ExecStart=/opt/filebrowser/filebrowser
Restart=on-failure

# Ensure Python is available
Environment="PATH=/usr/local/bin:/usr/bin:/bin"

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable filebrowser
sudo systemctl start filebrowser
```

### Pre-building Virtual Environment

For faster deployment, pre-build the venv:

```bash
# On your build machine
./filebrowser  # Let it create venv
# Wait for "Python virtual environment ready"
# Stop filebrowser

# Package everything
tar czf filebrowser-with-venv.tar.gz \
  filebrowser \
  transcribe.py \
  requirements.txt \
  .filebrowser_venv/

# Deploy to production
scp filebrowser-with-venv.tar.gz server:/opt/
ssh server "cd /opt && tar xzf filebrowser-with-venv.tar.gz"
```

---

## Advantages of Virtual Environment

### ✅ Isolation
- No conflicts with system Python packages
- No global package pollution
- Each filebrowser instance has its own environment

### ✅ Portability
- Easy to move between servers
- Consistent dependencies
- No "works on my machine" issues

### ✅ Easy Cleanup
- Delete `.filebrowser_venv/` to remove everything
- No leftover global packages

### ✅ Version Control
- Pin exact package versions
- Reproducible deployments
- Easy rollback

---

## Monitoring

### Check Virtual Environment Status

```bash
# Check if venv exists
ls -la .filebrowser_venv/

# Check installed packages
.filebrowser_venv/bin/pip list

# Check Python version
.filebrowser_venv/bin/python --version
```

### Logs

Check filebrowser logs for venv initialization:

```bash
# If using systemd
journalctl -u filebrowser -f | grep -i python

# If running manually
./filebrowser 2>&1 | grep -i python
```

---

## Uninstallation

To remove the transcription feature:

```bash
# Remove virtual environment
rm -rf .filebrowser_venv

# Remove scripts
rm transcribe.py requirements.txt
```

The virtual environment is completely isolated, so removing it doesn't affect your system Python.

---

## Comparison: Manual vs Automatic Setup

| Aspect | Manual Setup | Automatic Setup (New) |
|--------|-------------|----------------------|
| Installation | `pip install -r requirements.txt` | Automatic on startup |
| Scope | Global packages | Isolated venv |
| Conflicts | Possible | None |
| Cleanup | Manual uninstall | Delete one folder |
| Updates | Manual | Automatic on restart |
| Portability | System-dependent | Self-contained |

---

## Support

For issues or questions:

1. Check the [Troubleshooting](#troubleshooting) section
2. Review startup logs for error messages
3. Verify Python 3 is installed: `python3 --version`
4. Check disk space: `df -h`
5. Try recreating venv: `rm -rf .filebrowser_venv && ./filebrowser`

---

## Technical Details

### Virtual Environment Creation

Filebrowser uses Python's built-in `venv` module:
```bash
python3 -m venv .filebrowser_venv
```

### Dependency Installation

Uses pip within the venv:
```bash
.filebrowser_venv/bin/pip install -r requirements.txt
```

### Script Execution

Transcription uses the venv Python:
```bash
.filebrowser_venv/bin/python transcribe.py audio.mp3 output.txt
```

### Verification

On each startup, filebrowser:
1. Checks if `.filebrowser_venv/` exists
2. Verifies Python and pip are present
3. Checks if dependencies are installed
4. Recreates venv if corrupted
5. Installs missing dependencies

---

## License

The transcription feature uses:
- OpenAI Whisper (MIT License)
- Hugging Face Transformers (Apache 2.0)
- PyTorch (BSD License)

Ensure compliance with these licenses in your deployment.


---

## Detailed Setup

### Prerequisites

- Python 3.7 or higher
- pip (Python package manager)
- 4-8GB RAM (or GPU with 4GB+ VRAM)
- Internet connection (for first-time model download)

### Installation Steps

#### Option 1: Using requirements.txt (Recommended)

```bash
# Install all dependencies
pip install -r requirements.txt

# Verify installation
python3 -c "import torch; import transformers; print('Success!')"
```

#### Option 2: Manual Installation

```bash
# Install PyTorch (CPU version)
pip install torch torchvision torchaudio

# Install Transformers
pip install transformers

# Install Accelerate
pip install accelerate
```

#### Option 3: GPU Support (Recommended for Performance)

```bash
# Install CUDA-enabled PyTorch (for NVIDIA GPUs)
pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118

# Install other dependencies
pip install transformers accelerate

# Verify GPU is available
python3 -c "import torch; print(f'CUDA available: {torch.cuda.is_available()}')"
```

### First Run

The first time you transcribe audio:

1. The Whisper model will be downloaded (~3GB)
2. This may take 5-10 minutes depending on your connection
3. The model is cached for future use
4. Subsequent transcriptions will be much faster

### Testing

Test the script independently:

```bash
# Create a test audio file or use an existing one
python3 transcribe.py /path/to/audio.mp3 /path/to/output.txt

# Check the output
cat /path/to/output.txt
```

---

## Configuration

### Script Location

The `transcribe.py` script must be accessible from the filebrowser working directory.

**Option 1: Same directory as filebrowser**
```bash
# Place script next to filebrowser binary
cp transcribe.py /path/to/filebrowser/
```

**Option 2: System PATH**
```bash
# Make script executable
chmod +x transcribe.py

# Move to PATH location
sudo mv transcribe.py /usr/local/bin/

# Update backend to use 'transcribe.py' instead of './transcribe.py'
```

**Option 3: Absolute path in code**

Edit `http/resource.go`:
```go
cmd := exec.Command("python3", "/absolute/path/to/transcribe.py", srcRealPath, dstRealPath)
```

### Model Selection

Edit `transcribe.py` to change the Whisper model:

```python
# Faster but less accurate
model_id = "openai/whisper-tiny"    # ~39M parameters
model_id = "openai/whisper-base"    # ~74M parameters
model_id = "openai/whisper-small"   # ~244M parameters

# Balanced
model_id = "openai/whisper-medium"  # ~769M parameters

# Most accurate (default)
model_id = "openai/whisper-large-v3"  # ~1550M parameters
```

### Language Configuration

Edit `transcribe.py` to change the transcription language:

```python
# For Spanish
result = pipe(audio_file_path, generate_kwargs={"language": "spanish"})

# For French
result = pipe(audio_file_path, generate_kwargs={"language": "french"})

# For automatic detection (slower)
result = pipe(audio_file_path)
```

---

## Troubleshooting

### Common Issues

#### "No module named 'torch'"

```bash
# Install PyTorch
pip install torch
```

#### "No module named 'transformers'"

```bash
# Install Transformers
pip install transformers
```

#### "CUDA out of memory"

Solutions:
1. Use a smaller model (tiny, base, or small)
2. Close other GPU applications
3. Use CPU mode (automatic fallback)

Edit `transcribe.py`:
```python
# Force CPU mode
device = "cpu"
```

#### "Model download failed"

Solutions:
1. Check internet connection
2. Try again (downloads can be interrupted)
3. Manually download from Hugging Face:
   - Visit: https://huggingface.co/openai/whisper-large-v3
   - Download model files
   - Place in: `~/.cache/huggingface/hub/`

#### "Permission denied"

```bash
# Make script executable
chmod +x transcribe.py

# Or run with python3 explicitly
python3 transcribe.py audio.mp3 output.txt
```

#### Script not found by filebrowser

```bash
# Check script location
ls -la transcribe.py

# Ensure it's in the same directory as filebrowser
pwd
ls -la filebrowser transcribe.py

# Or use absolute path in resource.go
```

---

## Performance Optimization

### GPU Acceleration

For best performance, use a CUDA-enabled GPU:

```bash
# Install CUDA toolkit (Ubuntu/Debian)
sudo apt install nvidia-cuda-toolkit

# Install CUDA-enabled PyTorch
pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118

# Verify GPU is detected
python3 -c "import torch; print(torch.cuda.is_available())"
```

### Model Size vs Speed

| Model | Size | Speed (GPU) | Speed (CPU) | Accuracy |
|-------|------|-------------|-------------|----------|
| tiny | 39M | Very Fast | Fast | Good |
| base | 74M | Fast | Moderate | Better |
| small | 244M | Moderate | Slow | Good |
| medium | 769M | Slow | Very Slow | Better |
| large-v3 | 1550M | Very Slow | Extremely Slow | Best |

### Memory Requirements

| Model | GPU VRAM | System RAM |
|-------|----------|------------|
| tiny | 1GB | 2GB |
| base | 1GB | 2GB |
| small | 2GB | 4GB |
| medium | 5GB | 8GB |
| large-v3 | 10GB | 16GB |

---

## Docker Setup

If running filebrowser in Docker:

### Dockerfile

```dockerfile
FROM filebrowser/filebrowser:latest

# Install Python and dependencies
RUN apk add --no-cache python3 py3-pip

# Copy transcription script
COPY transcribe.py /app/transcribe.py
COPY requirements.txt /app/requirements.txt

# Install Python packages
RUN pip3 install -r /app/requirements.txt

# Set working directory
WORKDIR /app

# Run filebrowser
CMD ["filebrowser"]
```

### Docker Compose

```yaml
version: '3.8'
services:
  filebrowser:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/srv
      - ./database:/database
      - ./transcribe.py:/app/transcribe.py
    environment:
      - FB_DATABASE=/database/filebrowser.db
```

---

## Production Deployment

### Systemd Service

Create `/etc/systemd/system/filebrowser.service`:

```ini
[Unit]
Description=File Browser
After=network.target

[Service]
Type=simple
User=filebrowser
WorkingDirectory=/opt/filebrowser
ExecStart=/opt/filebrowser/filebrowser
Restart=on-failure

# Ensure Python is available
Environment="PATH=/usr/local/bin:/usr/bin:/bin"

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable filebrowser
sudo systemctl start filebrowser
```

### Resource Limits

Limit resource usage to prevent abuse:

```bash
# Limit CPU usage
nice -n 10 python3 transcribe.py audio.mp3 output.txt

# Limit memory
ulimit -v 8388608  # 8GB in KB
python3 transcribe.py audio.mp3 output.txt
```

Or use systemd limits in the service file:
```ini
[Service]
MemoryLimit=8G
CPUQuota=200%
```

---

## Monitoring

### Check Transcription Progress

```bash
# Monitor Python process
ps aux | grep transcribe.py

# Monitor GPU usage (if using CUDA)
nvidia-smi

# Monitor system resources
htop
```

### Logs

Check filebrowser logs for transcription errors:

```bash
# If using systemd
journalctl -u filebrowser -f

# If running manually
./filebrowser 2>&1 | tee filebrowser.log
```

---

## Uninstallation

To remove the transcription feature:

```bash
# Remove Python packages
pip uninstall torch transformers accelerate

# Remove cached models
rm -rf ~/.cache/huggingface/

# Remove script
rm transcribe.py requirements.txt
```

---

## Support

For issues or questions:

1. Check the [TRANSCRIBE_AUDIO_USAGE.md](TRANSCRIBE_AUDIO_USAGE.md) documentation
2. Review the [Troubleshooting](#troubleshooting) section
3. Check Whisper documentation: https://github.com/openai/whisper
4. Check Transformers documentation: https://huggingface.co/docs/transformers

---

## License

The transcription feature uses:
- OpenAI Whisper (MIT License)
- Hugging Face Transformers (Apache 2.0)
- PyTorch (BSD License)

Ensure compliance with these licenses in your deployment.
