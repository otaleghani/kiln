# Installation

**Kiln** is distributed as a single binary for macOS, Linux, and Windows. Select your operating system below for instructions.

## Using go
If you have go 1.25+ installed you can just simply install Kiln with go:
```bash
go install github.com/otaleghani/kiln/cmd/kiln@latest
```

## macOS (Apple Silicon / ARM64)

### Download the binary
Download the latest release for Apple Silicon:
```bash
curl -LO https://github.com/otaleghani/kiln/releases/latest/download/kiln_darwin_arm64
```

### Verify the Checksum (Recommended)
Ensure the file was downloaded correctly and has not been tampered with:
```bash
# Download checksums.txt
curl -LO https://github.com/otaleghani/kiln/releases/latest/download/checksums.txt

# This checks the downloaded binary against the checksum file
sha256sum -c checksums.txt --ignore-missing
```
_You should see `kiln_darwin_arm64: OK`._

### Install
Make the binary executable and move it to a directory in your `PATH` (e.g., `/usr/local/bin`).
```bash
chmod +x kiln_darwin_arm64
sudo mv kiln_darwin_arm64 /usr/local/bin/kiln
```

### Allow Execution (First Run Only)
_Since this binary is not notarized by Apple, you may need to allow it to run:_ Go to **System Settings** > **Privacy & Security**. Scroll down to the security section and click **Allow Anyway** next to the notification about `kiln`. _Alternatively, remove the quarantine attribute via terminal:_
```bash
xattr -d com.apple.quarantine /usr/local/bin/kiln
```

## Linux (AMD64)
### Download the binary
```bash
curl -LO https://github.com/otaleghani/kiln/releases/latest/download/kiln_linux_amd64
```

### Verify the Checksum (Recommended)
```bash
# Download checksums.txt
curl -LO https://github.com/otaleghani/kiln/releases/latest/download/checksums.txt

# This checks the downloaded binary against the checksum file
sha256sum -c checksums.txt --ignore-missing
```
_You should see `kiln_linux_amd64: OK`._

### Install 
Make the binary executable and move it to `/usr/local/bin`:
```bash
chmod +x kiln_linux_amd64
sudo mv kiln_linux_amd64 /usr/local/bin/kiln
```

## Windows (AMD64)
### Download the binary
Download `kiln_windows_amd64.exe` from the [Releases Page](https://github.com/otaleghani/kiln/releases/latest) or via PowerShell: 
```powershell
Invoke-WebRequest -Uri "https://github.com/otaleghani/kiln/releases/latest/download/kiln_windows_amd64.exe" -OutFile "kiln.exe"
```
### Verify the Checksum (Recommended) 
Run the following in PowerShell to verify the hash matches
```powershell
$expected = Select-String -Path .\checksums.txt -Pattern "kiln_windows_amd64.exe" | ForEach-Object { $_.Line.Split(' ')[0] };
(Get-FileHash .\kiln_windows_amd64.exe -Algorithm SHA256).Hash.ToLower() -eq $expected
```
_This should return `True`._

### Install
Move `kiln.exe` to a folder of your choice (e.g., `C:\Program Files\Kiln\`) and add that folder to your System `PATH` environment variable so you can run `kiln` from any terminal window.